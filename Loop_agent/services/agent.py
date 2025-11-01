import logging
import os
import traceback
from typing import Dict, List, Optional
from .go_client import GoBackendClient

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class SearchAgent:
    """Advanced project search agent"""
    
    def __init__(self, go_client: GoBackendClient):
        self.go_client = go_client
        self.api_key = os.environ.get("API_KEY", "")
        self.api_url = os.environ.get("API_URL", "https://api.groq.com/openai")
        self.model_name = os.environ.get("MODEL_NAME", "meta-llama/llama-4-maverick-17b-128e-instruct")
        self.conversation_history = {}
        
        # Configure the agent with detailed instructions
        self.system_prompt = """
        You are an expert project discovery specialist with deep knowledge of software architecture.
        Your expertise encompasses:

        - Semantic similarity analysis
        - Technical stack evaluation
        - Project architecture understanding
        - Relationship pattern recognition
        - Technology trend analysis
        - Framework compatibility assessment
        - Development paradigm matching
        - Best practices identification
        - Cross-domain knowledge application
        - Technical requirement analysis

        Search Strategy Guidelines:
        
        1. Vector Search (Similarity-Based):
           - Use for finding semantically similar projects
           - Best for 'like X' or 'similar to X' queries
           - Handles natural language and conceptual matching
           - Example: 'Show me projects similar to Netflix'
        
        2. Graph Search (Conditional):
           - Use for specific attribute-based filtering
           - Best for queries with explicit conditions
           - Handles relationships and technical requirements
           - Example: 'Find React projects using GraphQL'
        
        3. Hybrid Search (Two-Stage):
           - First: Vector search for relevant candidates
           - Then: Graph search to filter and augment results
           - Best for complex queries combining similarity and conditions
           - Example: 'Projects like Airbnb but using Vue.js'
        """
        logger.info("Search agent initialized with system prompt")

    async def process_query(self, query: str, user_id: str, deep_search: bool = False) -> Dict:
        """
        Process a search query with optional deep search capabilities.
        
        Args:
            query: The search query string
            user_id: Unique identifier for the user
            deep_search: Whether to perform a more thorough analysis
        
        Returns:
            Dict containing search results and analysis
        """
        logger.info(f"Processing query - User: {user_id}, Query: {query}, Deep Search: {deep_search}")
        
        try:
            # Initialize or get conversation history
            if user_id not in self.conversation_history:
                self.conversation_history[user_id] = []
                logger.info(f"Initialized new conversation history for user: {user_id}")
                
            # Add query to history
            self.conversation_history[user_id].append({
                "role": "user",
                "content": query
            })
            
            # Check if this is a non-search query
            if self._is_non_search_query(query):
                logger.info(f"Detected non-search query: {query}")
                response = self._generate_non_search_response()
                
                # Add response to history
                self.conversation_history[user_id].append({
                    "role": "assistant",
                    "content": response
                })
                
                return {
                    "type": "non_search",
                    "message": response,
                    "results": []
                }

            # Use LLM to analyze the query
            logger.info("Analyzing query with LLM...")
            analysis = self._analyze_query_with_llm(query)
            
            # Check if clarification is needed
            if analysis.get("needs_clarification", False):
                logger.info(f"Clarification needed: {analysis.get('reason')}")
                questions = analysis.get("clarifying_questions", [])
                if questions:
                    formatted_questions = "\n".join(f"- {q}" for q in questions)
                    return {
                        "type": "clarification",
                        "message": f"To help you find the most relevant projects, I need some clarification:\n\n{formatted_questions}",
                        "partial_results": [],
                        "analysis": analysis.get("analysis", "")
                    }

            # Use LLM's analysis for search
            search_mode = analysis.get("search_mode", "vector")
            search_terms = analysis.get("extracted_terms", query.lower())
            logger.info(f"Using search mode: {search_mode}, terms: {search_terms}")
            
            # Execute the search with the determined mode
            logger.info(f"Executing search with mode: {search_mode}")
            try:
                search_results = self.go_client.search_tool(search_terms, mode=search_mode)
                logger.info(f"Search completed successfully. Found {len(search_results)} results.")
            except Exception as e:
                logger.error(f"Error during search execution: {str(e)}\n{traceback.format_exc()}")
                raise
            
            # Format the response
            logger.info("Formatting search results...")
            response_content = self._format_search_results(query, search_results, search_mode)
            
            # Add response to history
            self.conversation_history[user_id].append({
                "role": "assistant",
                "content": response_content
            })
            
            logger.info("Query processing completed successfully")
            return {
                "type": "results",
                "message": response_content,
                "results": search_results,
                "search_mode": search_mode
            }
            
        except Exception as e:
            logger.error(f"Error in process_query: {str(e)}\n{traceback.format_exc()}")
            raise

    def _extract_search_terms(self, query: str) -> str:
        """Extract relevant search terms from the query."""
        query_lower = query.lower()
        
        # Remove common prefixes
        prefixes = [
            "show me", "find", "search for", "looking for", "i want", "i need",
            "can you find", "please find", "get me", "projects about", "projects on",
            "projects using", "projects with"
        ]
        
        cleaned_query = query_lower
        for prefix in prefixes:
            if cleaned_query.startswith(prefix):
                cleaned_query = cleaned_query[len(prefix):].strip()
                break
        
        # Remove common suffixes
        suffixes = [
            "projects", "project", "please", "thanks", "thank you"
        ]
        
        for suffix in suffixes:
            if cleaned_query.endswith(suffix):
                cleaned_query = cleaned_query[:-len(suffix)].strip()
        
        logger.info(f"Extracted search terms from '{query}' -> '{cleaned_query}'")
        return cleaned_query

    def _is_basic_topic_search(self, query: str) -> bool:
        """Determine if this is a basic topic search."""
        query_lower = query.lower()
        
        # If query contains specific search indicators, it's not a basic search
        complex_indicators = [
            "using", "with", "built in", "framework", "language", "like",
            "similar to", "resembling", "same as", "comparable to"
        ]
        
        if any(indicator in query_lower for indicator in complex_indicators):
            return False
            
        # If query is just 1-3 words about a topic, it's a basic search
        words = query_lower.split()
        return 1 <= len(words) <= 3

    def _analyze_query(self, query: str) -> str:
        """Determine the best search mode for the query"""
        logger.info(f"Analyzing query: {query}")
        query_lower = query.lower()
        
        # Check for explicit conditions that suggest graph search
        graph_indicators = ["using", "with", "built in", "framework", "language", 
                           "technology", "database", "backend", "frontend"]
        
        # Check for similarity indicators that suggest vector search
        vector_indicators = ["like", "similar to", "resembling", "same as", 
                            "comparable to", "akin to", "equivalent to"]
        
        has_graph_indicators = any(indicator in query_lower for indicator in graph_indicators)
        has_vector_indicators = any(indicator in query_lower for indicator in vector_indicators)
        
        if has_vector_indicators and has_graph_indicators:
            logger.info("Query contains both vector and graph indicators -> Using hybrid search")
            return "hybrid"
        elif has_graph_indicators:
            logger.info("Query contains graph indicators -> Using graph search")
            return "graph"
        elif has_vector_indicators or self._is_basic_topic_search(query):
            logger.info("Query is vector-based or basic topic -> Using vector search")
            return "vector"
        else:
            logger.info("Complex query -> Using hybrid search")
            return "hybrid"

    def _analyze_query_with_llm(self, query: str) -> Dict:
        """
        Use LLM to analyze the query and determine if clarification is needed.
        Returns a dict with analysis results.
        """
        logger.info(f"Analyzing query with LLM: {query}")
        
        prompt = f"""As a project discovery specialist, analyze this search query: "{query}"

Your goal is to help users find relevant software projects. Consider:

1. Query Understanding:
   - Is the search intent clear and specific?
   - Are there implied technical requirements?
   - What domain or type of project is the user looking for?

2. Search Strategy:
   - Vector Search: Best for finding similar projects or topic-based searches
   - Graph Search: Best for technical requirements and specific features
   - Hybrid Search: For complex queries combining similarity and requirements

3. Clarification Needs:
   - Would additional context help find better results?
   - Are there ambiguous terms that need clarification?
   - Could the search be refined with more specific requirements?

4. Project Relevance:
   - What makes a project relevant to this query?
   - What aspects should be prioritized in the search?
   - What additional filters might be helpful?

Respond in JSON format:
{{
  "needs_clarification": true/false,
  "reason": "Detailed explanation of why clarification is needed or not",
  "clarifying_questions": [
    "Specific, contextual questions that will help refine the search"
  ],
  "search_mode": "vector/graph/hybrid",
  "extracted_terms": "Key search terms extracted from query",
  "analysis": "Comprehensive analysis of search intent and strategy",
  "suggested_filters": {{
    "technologies": ["relevant tech stack"],
    "domain": "suggested domain focus",
    "complexity": "suggested project size/complexity"
  }}
}}

Focus on generating highly relevant results that match the user's true intent."""

        try:
            headers = {
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json"
            }
            
            response = requests.post(
                f"{self.api_url}/v1/chat/completions",
                headers=headers,
                json={
                    "model": self.model_name,
                    "messages": [
                        {"role": "system", "content": self.system_prompt},
                        {"role": "user", "content": prompt}
                    ],
                    "temperature": 0.7,
                    "max_tokens": 500
                }
            )
            
            response.raise_for_status()
            result = response.json()
            analysis = json.loads(result['choices'][0]['message']['content'])
            logger.info(f"LLM analysis complete: {analysis}")
            return analysis
            
        except Exception as e:
            logger.error(f"Error in LLM analysis: {str(e)}\n{traceback.format_exc()}")
            # Return a default analysis if LLM fails
            return {
                "needs_clarification": False,
                "reason": "LLM analysis failed, proceeding with basic search",
                "search_mode": "vector",
                "extracted_terms": query.lower(),
                "analysis": "Using default analysis due to LLM error"
            }

    def _needs_clarification(self, query: str) -> bool:
        """Use LLM to determine if query needs clarification"""
        analysis = self._analyze_query_with_llm(query)
        return analysis.get("needs_clarification", False)

    def _generate_clarification_questions(self, query: str) -> str:
        """Generate clarification questions using LLM analysis"""
        analysis = self._analyze_query_with_llm(query)
        
        if not analysis.get("needs_clarification"):
            return ""
            
        questions = analysis.get("clarifying_questions", [])
        if not questions:
            # Fallback questions if LLM didn't provide any
            questions = [
                f"Could you provide more details about what aspects of '{query}' are most important to you?",
                "Are you looking for specific technologies or frameworks?",
                "Do you have any preferences regarding project size or complexity?"
            ]
            
        formatted_questions = "\n".join(f"- {q}" for q in questions)
        logger.info(f"Generated clarification questions using LLM")
        return formatted_questions

    def _format_search_results(self, query: str, results: List[Dict], search_mode: str) -> str:
        """Format search results in a structured way"""
        logger.info(f"Formatting search results - Query: {query}, Mode: {search_mode}, Results: {len(results)}")
        
        if not results:
            logger.info("No results found")
            return f"# No matching projects found\n\nI couldn't find any projects matching your query: '{query}'"
            
        # Start with the header and query analysis
        formatted_response = f"""# Search Results 🎯

## Query Analysis
I analyzed your query: "{query}"

This query is best suited for {search_mode} search because """
        
        if search_mode == "vector":
            formatted_response += "it focuses on finding conceptually similar projects."
        elif search_mode == "graph":
            formatted_response += "it contains specific technical requirements or conditions."
        else:
            formatted_response += "it combines both similarity and specific requirements."
            
        # Add primary matches section
        formatted_response += "\n\n## Primary Matches\n"
        
        for i, result in enumerate(results[:5], 1):
            project_name = result.get("name", f"Project {i}")
            description = result.get("description", "No description available")
            score = result.get("score", 0)
            
            formatted_response += f"\n### {i}. {project_name}\n"
            formatted_response += f"**Relevance Score**: {score:.2f}\n\n"
            formatted_response += f"{description}\n"
            
            # Add technical details if available
            if "technologies" in result:
                techs = result.get("technologies", [])
                if techs:
                    formatted_response += "\n**Technologies**: " + ", ".join(techs) + "\n"
                    
        # Add search methodology
        formatted_response += f"""
## Search Methodology
**Search Type Used**: {search_mode.capitalize()} Search
**Ranking Factors**: Semantic relevance, keyword matching, and technical alignment
**Results Count**: {len(results)}

---
Search conducted by Project Discovery Specialist
Using: Vector Search 🎯 | Graph Search 🔍 | Hybrid Search 🔄
"""
        
        logger.info("Results formatting completed")
        return formatted_response

    def _is_non_search_query(self, query: str) -> bool:
        """
        Determine if the query is not a search request.
        Examples: greetings, casual conversation, etc.
        """
        query_lower = query.lower().strip()
        
        # Common non-search patterns
        non_search_patterns = [
            "hi", "hello", "hey", "good morning", "good afternoon", 
            "good evening", "how are you", "what's up", "sup",
            "thanks", "thank you", "bye", "goodbye"
        ]
        
        # Check if query matches any non-search pattern
        if any(query_lower == pattern for pattern in non_search_patterns):
            logger.info(f"Query '{query}' identified as non-search query")
            return True
            
        # Check if query is too short or lacks search intent
        if len(query_lower.split()) < 2 and not any(char.isdigit() for char in query_lower):
            logger.info(f"Query '{query}' is too short for search")
            return True
            
        return False

    def _generate_non_search_response(self) -> str:
        """Generate a response for non-search queries."""
        response = """Welcome to Loop!\n I'm here to help you discover amazing projects.\n
 Search by Similarity:\n
   - "Show me projects similar to Netflix"\n
   - "Find projects like Airbnb"\n\n

Search by Technology\n
   - "Find React projects using GraphQL"\n
   - "Show me Python machine learning projects"\n\n

Combined Search:\n
   - "Projects like Twitter but using Vue.js"\n
   - "E-commerce projects using Node.js and MongoDB"\n\n

What kind of projects would you like to discover?"""

        logger.info("Generated non-search response with search suggestions")
        return response

    def clear_conversation(self, user_id: str) -> None:
        """Clear the conversation history for a user"""
        logger.info(f"Clearing conversation history for user: {user_id}")
        if user_id in self.conversation_history:
            self.conversation_history[user_id] = []
            logger.info(f"Conversation history cleared for user: {user_id}")
