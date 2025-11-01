"""
Team coordination module for managing agent interactions.
"""

import logging
from typing import Dict, List, Any, Optional
from services.agents.project_search import ProjectSearchAgent
from services.agents.clarification import ClarificationAgent
from services.agents.deep_search import DeepSearchAgent

logger = logging.getLogger(__name__)

class ProjectDiscoveryTeam:
    """Coordinates the agents for project discovery"""
    
    def __init__(self, go_client):
        """
        Initialize the project discovery team.
        
        Args:
            go_client: Client for communicating with the Go backend
        """
        self.project_search = ProjectSearchAgent(go_client)
        self.clarification = ClarificationAgent()
        self.deep_search = DeepSearchAgent(go_client)
        self.conversation_history = {}
        self.sid_to_user = {}  # Map socket IDs to user IDs
        self.sid_to_token = {}  # Map socket IDs to auth tokens
        
    def register_session(self, sid: str, user_id: str, token: str = None) -> None:
        """
        Register a socket session for a user.
        
        Args:
            sid: Socket session ID
            user_id: User identifier
            token: Authentication token
        """
        self.sid_to_user[sid] = user_id
        if token:
            self.sid_to_token[sid] = token
            logger.info(f"Stored auth token for session {sid}")
        
        if user_id not in self.conversation_history:
            self.conversation_history[user_id] = []
        logger.info(f"Registered session {sid} for user {user_id}")
    
    def get_conversation_context(self, user_id: str, max_history: int = 5) -> List[Dict]:
        """
        Get recent conversation history for context.
        
        Args:
            user_id: User identifier
            max_history: Maximum number of recent messages to include
            
        Returns:
            List of recent messages
        """
        if user_id not in self.conversation_history:
            return []
        return self.conversation_history[user_id][-max_history:]
    
    async def process_query(self, query: str, user_id: str, deep_search: bool = False, token: str = None) -> Dict:
        """
        Process a user query through the agent team.
        
        Args:
            query: The user's search query
            user_id: Unique identifier for the user
            deep_search: Whether to perform deep analysis
            
        Returns:
            Dict containing search results and analysis
        """
        logger.info(f"Processing query - User: {user_id}, Query: {query}, Deep Search: {deep_search}")
        
        # Initialize conversation history if needed
        if user_id not in self.conversation_history:
            self.conversation_history[user_id] = []
            
        # Add query to history
        self.conversation_history[user_id].append({
            "role": "user",
            "content": query
        })
        
        try:
            # Get conversation context
            context = self.get_conversation_context(user_id)
            
            # Get token from parameters or lookup by user_id
            access_token = token
            if not access_token:
                # Try to find token from any session associated with this user
                for sid, uid in self.sid_to_user.items():
                    if uid == user_id and sid in self.sid_to_token:
                        access_token = self.sid_to_token[sid]
                        break
            
            # Commenting out clarification agent as per requirements
            # We'll directly search without asking clarification questions
            """
            # Check if clarification is needed
            clarification_analysis = await self.clarification.analyze_query(query, context)
            if clarification_analysis.get("needs_clarification", False):
                questions = await self.clarification.generate_questions(query, context)
                if questions:
                    formatted_questions = self.clarification.format_questions(questions)
                    return {
                        "type": "clarification",
                        "message": formatted_questions,
                        "partial_results": [],
                        "analysis": clarification_analysis.get("reason", "")
                    }
            """
            
            # Process search with context
            search_results = await self.project_search.search(query, context, access_token)
            
            # Handle search errors
            if "error" in search_results:
                return {
                    "type": "error",
                    "message": search_results["error"],
                    "query": query
                }
            
            # Perform deep analysis if requested
            if deep_search and search_results.get("results"):
                deep_analysis = await self.deep_search.analyze_results(
                    query, 
                    search_results.get("results", []),
                    context,
                    access_token
                )
                search_results["deep_analysis"] = deep_analysis
            
            # Format response
            response_content = self._format_response(search_results)
            
            # Add response to history
            self.conversation_history[user_id].append({
                "role": "assistant",
                "content": response_content
            })
            
            # Check if there's a specific message to use (for vague queries)
            message = search_results.get("message", response_content)
            
            return {
                "type": "results",
                "message": message,  # Use the message from search_results if available
                "results": search_results.get("results", []),
                "search_mode": search_results.get("search_mode", "vector"),
                "extracted_terms": search_results.get("extracted_terms", ""),
                "deep_analysis": search_results.get("deep_analysis") if deep_search else None
            }
            
        except Exception as e:
            logger.error(f"Error processing query: {str(e)}")
            return {
                "type": "error",
                "message": f"An error occurred while processing your request: {str(e)}",
                "query": query
            }
    
    def _format_response(self, search_results: Dict) -> str:
        """
        Format the complete response including search results and analysis.
        
        Args:
            search_results: Dict containing search results and metadata
            
        Returns:
            Formatted string for display
        """
        # Get basic search results formatting
        formatted = self.project_search.format_results(search_results)
        
        # Add deep analysis if available
        if "deep_analysis" in search_results:
            deep_analysis = self.deep_search.format_analysis(search_results["deep_analysis"])
            formatted += f"\n\n{deep_analysis}"
        
        return formatted
    
    def clear_conversation(self, user_id: str) -> None:
        """
        Clear the conversation history for a user.
        
        Args:
            user_id: Unique identifier for the user
        """
        if user_id in self.conversation_history:
            self.conversation_history[user_id] = []
            logger.info(f"Cleared conversation history for user: {user_id}")
