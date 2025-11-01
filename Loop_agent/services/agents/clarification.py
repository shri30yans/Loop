"""
Clarification Agent implementation.
This agent is responsible for analyzing queries and generating clarifying questions when needed.
"""

import logging
from typing import Dict, List, Any, Optional
from services.agents.base import BaseAgent

logger = logging.getLogger(__name__)

class ClarificationAgent(BaseAgent):
    """Agent responsible for generating clarifying questions for ambiguous queries"""
    
    def __init__(self):
        """Initialize the clarification agent."""
        super().__init__(
            name="Clarification Agent",
            role="Generate clarifying questions for ambiguous project search queries"
        )
        self.system_prompt = """
        You are an expert at analyzing project search queries and identifying ambiguities.
        Your task is to determine if a query needs clarification and generate appropriate questions.
        
        IMPORTANT: Simple queries like "projects about agriculture" or "machine learning projects" 
        are acceptable and DO NOT need clarification. Only ask for clarification when the query is 
        truly ambiguous or lacks critical information needed to provide relevant results.
        
        Focus on:
        - Highly ambiguous technical requirements
        - Extremely vague project descriptions with multiple possible interpretations
        - Missing critical context that would significantly impact search results
        - Unclear user intent that could lead to completely irrelevant results
        
        When generating clarifying questions:
        - Be specific and contextual
        - Focus on technical aspects that would improve search results
        - Limit to a MAXIMUM of 4 questions
        - For simple queries, prefer to search directly rather than asking for clarification
        """
    
    async def analyze_query(self, query: str, context: List[Dict] = None) -> Dict:
        """
        Analyze if a query needs clarification and generate questions.
        
        Args:
            query: The user's search query
            context: Previous conversation history
            
        Returns:
            Dict containing analysis and clarifying questions
        """
        # Format conversation context if available
        context_str = ""
        if context and len(context) > 0:
            context_str = "Previous conversation:\n"
            for msg in context:
                role = msg.get("role", "unknown")
                content = msg.get("content", "")
                context_str += f"{role}: {content}\n"
        
        prompt = f"""Analyze this project search query: "{query}"

{context_str}

Determine if clarification is needed and generate appropriate questions.

Respond in JSON format:
{{
  "needs_clarification": true/false,
  "reason": "Detailed explanation of why clarification is needed or not",
  "clarifying_questions": [
    "Specific, contextual questions that will help refine the search"
  ]
}}
"""
        
        response = await self.call_llm(self.system_prompt, prompt)
        analysis = self.parse_json_response(response)
        
        logger.info(f"Clarification analysis for '{query}': needs_clarification={analysis.get('needs_clarification', False)}")
        return analysis
    
    async def generate_questions(self, query: str, context: List[Dict] = None) -> List[str]:
        """
        Generate clarifying questions for a query.
        
        Args:
            query: The user's search query
            context: Previous conversation history
            
        Returns:
            List of clarifying questions
        """
        analysis = await self.analyze_query(query, context)
        
        if not analysis.get("needs_clarification", False):
            logger.info(f"No clarification needed for query: '{query}'")
            return []
            
        questions = analysis.get("clarifying_questions", [])
        
        if not questions:
            # Fallback questions if LLM didn't provide any
            questions = [
                f"Could you provide more details about what aspects of '{query}' are most important to you?",
                "Are you looking for specific technologies or frameworks?",
                "Do you have any preferences regarding project size or complexity?"
            ]
            logger.info(f"Using fallback questions for query: '{query}'")
        
        logger.info(f"Generated {len(questions)} clarifying questions for query: '{query}'")
        return questions
    
    def format_questions(self, questions: List[str]) -> str:
        """
        Format clarifying questions for display.
        
        Args:
            questions: List of clarifying questions
            
        Returns:
            Formatted string for display
        """
        if not questions:
            return ""
            
        formatted = "To help you find the most relevant projects, I need some clarification:\n\n"
        formatted += "\n".join(f"- {q}" for q in questions)
        return formatted
