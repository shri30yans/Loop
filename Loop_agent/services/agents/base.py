"""
Base agent implementation providing common functionality for all agents.
"""

import os
import logging
import json
import requests
from typing import Dict, List, Any, Optional

# Configure logging
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class BaseAgent:
    """Base class for all agents in the system"""
    
    def __init__(self, name: str, role: str, model_name: str = None):
        """
        Initialize the base agent.
        
        Args:
            name: Agent's name
            role: Agent's role description
            model_name: Optional specific model to use (defaults to env var)
        """
        self.name = name
        self.role = role
        self.api_key = os.environ.get("API_KEY", "your_default_api_key")
        self.api_url = os.environ.get("API_URL", "https://api.groq.com/openai")
        self.model_name = model_name or os.environ.get("MODEL_NAME", "meta-llama/llama-4-maverick-17b-128e-instruct")
        
    async def call_llm(self, system_prompt: str, user_prompt: str, temperature: float = 0.7) -> Dict:
        """
        Call the LLM with the given prompts.
        
        Args:
            system_prompt: System-level instructions for the LLM
            user_prompt: User query or task description
            temperature: Sampling temperature (0.0 to 1.0)
            
        Returns:
            Dict containing the LLM's response
        """
        try:
            import aiohttp
            
            headers = {
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json"
            }
            
            async with aiohttp.ClientSession() as session:
                async with session.post(
                    f"{self.api_url}/v1/chat/completions",
                    headers=headers,
                    json={
                        "model": self.model_name,
                        "messages": [
                            {"role": "system", "content": system_prompt},
                            {"role": "user", "content": user_prompt}
                        ],
                        "temperature": temperature,
                        "max_tokens": 500
                    }
                ) as response:
                    response.raise_for_status()
                    return await response.json()
            
        except Exception as e:
            logger.error(f"Error calling LLM: {str(e)}")
            return {"error": str(e)}
    
    def parse_json_response(self, response: Dict) -> Dict:
        """
        Parse JSON response from LLM.
        
        Args:
            response: Raw response from LLM
            
        Returns:
            Dict containing parsed JSON data
        """
        try:
            # First check if we have a valid response
            if not response or "error" in response:
                return {"error": response.get("error", "Invalid response from LLM")}
            
            # Check if we have the expected structure
            if not response.get("choices") or not response["choices"]:
                return {"error": "No choices in LLM response"}
                
            content = response["choices"][0]["message"]["content"]
            if not content:
                return {"error": "Empty content in LLM response"}
            
            # Try to find JSON content within the response
            # Sometimes the LLM might include additional text before/after the JSON
            try:
                # First try direct JSON parsing
                return json.loads(content)
            except json.JSONDecodeError:
                # If that fails, try to find JSON-like content between curly braces
                import re
                json_match = re.search(r"\{[\s\S]*\}", content)
                if json_match:
                    return json.loads(json_match.group())
                else:
                    # If no JSON found, return a formatted error
                    return {
                        "error": "Could not parse JSON from response",
                        "raw_content": content
                    }
                    
        except Exception as e:
            logger.error(f"Error parsing JSON response: {str(e)}")
            return {
                "error": f"Failed to parse response: {str(e)}",
                "raw_content": response.get("choices", [{}])[0].get("message", {}).get("content", "")
            }
            
    def format_response(self, content: Dict) -> str:
        """
        Format a response for user display.
        
        Args:
            content: Dict containing response data
            
        Returns:
            Formatted string for display
        """
        try:
            return json.dumps(content, indent=2)
        except Exception as e:
            logger.error(f"Error formatting response: {str(e)}")
            return str(content)
            
    def log_error(self, error: str, context: Dict = None) -> None:
        """
        Log an error with optional context.
        
        Args:
            error: Error message
            context: Optional dict with additional context
        """
        if context:
            logger.error(f"{error} - Context: {json.dumps(context)}")
        else:
            logger.error(error)
