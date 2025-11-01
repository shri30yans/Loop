"""Client for interacting with the Go backend service"""

import aiohttp
import json
import logging
import traceback
import os
from typing import Dict, List, Optional

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class GoBackendClient:
    """Client for interacting with the Go backend service"""
    
    def __init__(self, base_url: str, access_token: str = None):
        self.base_url = base_url.rstrip('/')
        self.access_token = access_token or os.getenv('GO_BACKEND_API_TOKEN')
        logger.info(f"Initialized Go backend client with base URL: {self.base_url}")
        self._session = None
        self._session_lock = asyncio.Lock()

    async def _get_session(self) -> aiohttp.ClientSession:
        """Get or create an aiohttp session with proper locking to avoid race conditions."""
        async with self._session_lock:
            if self._session is None or self._session.closed:
                # Configure timeout to prevent hanging requests
                timeout = aiohttp.ClientTimeout(total=30, connect=10, sock_connect=10, sock_read=10)
                connector = aiohttp.TCPConnector(limit=100, ttl_dns_cache=300)
                self._session = aiohttp.ClientSession(timeout=timeout, connector=connector)
            return self._session

    async def search_tool(self, query: str, mode: str = "hybrid", access_token: str = None) -> List[Dict]:
        """
        Execute a search query using the specified mode.
        
        Args:
            query: The search query string (already cleaned and processed)
            mode: Search mode ('vector', 'graph', or 'hybrid')
            access_token: Optional token to override the default one
            
        Returns:
            List of search results
        """
        logger.info(f"Executing search - Query: {query}, Mode: {mode}")
        
        # Use the specific endpoint for the requested mode
        if mode == "vector":
            endpoint = f"{self.base_url}/api/tools/search/vector"
        elif mode == "graph":
            endpoint = f"{self.base_url}/api/tools/search/graph"
        elif mode == "hybrid":
            endpoint = f"{self.base_url}/api/tools/search/hybrid"
        else:
            # Fallback to the unified endpoint for backward compatibility
            endpoint = f"{self.base_url}/api/search_tools"
        
        # Prepare the request with mode
        payload = {
            "query": query,
            "mode": mode  # Include mode in payload for backward compatibility
        }
        
        try:
            # Make the request
            return await self._execute_search_request(endpoint, payload, access_token)
        except Exception as e:
            error_msg = str(e)
            logger.error(f"Error in search_tool: {error_msg}\n{traceback.format_exc()}")
            
            # If the error is related to the endpoint not being found, try the legacy endpoint
            if isinstance(e, aiohttp.ClientResponseError) and e.status == 404:
                logger.info("Endpoint not found, falling back to legacy endpoint")
                try:
                    return await self._execute_search_request(f"{self.base_url}/api/search_tools", payload, access_token)
                except Exception as fallback_error:
                    logger.error(f"Fallback request also failed: {str(fallback_error)}")
                    raise
            
            # Propagate the original error
            raise
    
    async def _execute_search_request(self, endpoint: str, payload: Dict, access_token: str = None) -> List[Dict]:
        """
        Execute a search request to the specified endpoint.
        
        Args:
            endpoint: The API endpoint to call
            payload: The request payload
            access_token: Optional token to override the default one
            
        Returns:
            List of processed search results
        """
        # Set up headers with authentication
        token = access_token or self.access_token
        headers = {
            "Content-Type": "application/json",
            "Accept": "application/json"
        }
        
        if token:
            headers["Authorization"] = f"Bearer {token}"
            logger.info("Including authorization token in request")
        else:
            logger.warning("No access token available for authentication")
        
        logger.info(f"Sending request to {endpoint}")
        session = await self._get_session()
        
        response = None
        try:
            # Use a timeout context to prevent hanging requests
            async with session.post(endpoint, json=payload, headers=headers) as response:
                # Check for successful response
                response.raise_for_status()
                
                # Parse response - use a separate task to avoid blocking
                data = await asyncio.create_task(response.json())
                results = data.get("results", [])
                logger.info(f"Search completed successfully. Found {len(results)} results.")
                
                # Process results asynchronously to avoid blocking
                return await self._process_search_results(results, payload.get("mode", "hybrid"))
                
        except aiohttp.ClientResponseError as e:
            # Handle HTTP errors with detailed information
            error_details = {}
            if response:
                try:
                    error_text = await response.text()
                    error_details = json.loads(error_text)
                except Exception:
                    error_details = {"raw_error": error_text if 'error_text' in locals() else "Unknown error"}
            
            logger.error(f"HTTP error {e.status}: {e.message}")
            logger.error(f"Error details: {error_details}")
            
            # Re-raise with more context
            raise Exception(f"API request failed with status {e.status}: {error_details.get('message', e.message)}")
            
        except aiohttp.ClientError as e:
            logger.error(f"HTTP request failed: {str(e)}\n{traceback.format_exc()}")
            raise Exception(f"Connection error: {str(e)}")
            
        except asyncio.TimeoutError:
            logger.error("Request timed out")
            raise Exception("Request timed out. The server took too long to respond.")
            
        except Exception as e:
            logger.error(f"Unexpected error: {str(e)}\n{traceback.format_exc()}")
            raise
            
    async def _process_search_results(self, results: List[Dict], mode: str) -> List[Dict]:
        """Process search results asynchronously to avoid blocking."""
        processed_results = []
        
        # Process each result
        for result in results:
            # Extract core fields
            processed_result = {
                "name": result.get("name", "Unnamed Project"),
                "description": result.get("description", "No description available"),
                "score": float(result.get("score", 0.0)),  # Ensure score is float
                "technologies": result.get("technologies", []),
                "url": result.get("url", ""),
                "stars": result.get("stars", 0),
                "last_updated": result.get("last_updated", ""),
                "language": result.get("language", "")
            }
            
            # Add mode-specific fields
            if mode == "vector":
                processed_result["similarity_score"] = float(result.get("similarity", 0.0))
            elif mode == "graph":
                processed_result["match_score"] = float(result.get("match_score", 0.0))
                processed_result["relationships"] = result.get("relationships", [])
            else:  # hybrid
                processed_result["combined_score"] = float(result.get("combined_score", 0.0))
                processed_result["vector_score"] = float(result.get("vector_score", 0.0))
                processed_result["graph_score"] = float(result.get("graph_score", 0.0))
            
            processed_results.append(processed_result)
        
        return processed_results

    async def _make_request(self, method: str, endpoint: str, **kwargs) -> Dict:
        """
        Make an HTTP request to the Go backend.
        
        Args:
            method: HTTP method ('GET', 'POST', etc.)
            endpoint: API endpoint path
            **kwargs: Additional arguments for aiohttp
            
        Returns:
            Response data as dictionary
        """
        url = f"{self.base_url}/{endpoint.lstrip('/')}"
        logger.info(f"Making {method} request to {url}")
        
        # Ensure headers are set
        kwargs.setdefault('headers', {})
        if 'headers' in kwargs and isinstance(kwargs['headers'], dict):
            # Set default content type if not present
            kwargs['headers'].setdefault('Content-Type', 'application/json')
            kwargs['headers'].setdefault('Accept', 'application/json')
            
            # Add authorization if token is available and not already set
            if self.access_token and 'Authorization' not in kwargs['headers']:
                kwargs['headers']['Authorization'] = f"Bearer {self.access_token}"
        
        response = None
        try:
            session = await self._get_session()
            async with session.request(method, url, **kwargs) as response:
                response.raise_for_status()
                # Use create_task to avoid blocking
                return await asyncio.create_task(response.json())
                
        except aiohttp.ClientResponseError as e:
            # Try to parse error response
            error_details = {}
            if response:
                try:
                    error_text = await response.text()
                    error_details = json.loads(error_text)
                    logger.error(f"API error response: {error_details}")
                except Exception:
                    logger.error(f"Could not parse error response: {e}")
            
            # Raise a more informative exception
            error_message = error_details.get('message', str(e))
            raise Exception(f"API request failed ({e.status}): {error_message}")
            
        except aiohttp.ClientError as e:
            logger.error(f"HTTP request failed: {str(e)}\n{traceback.format_exc()}")
            raise Exception(f"Connection error: {str(e)}")
            
        except asyncio.TimeoutError:
            logger.error("Request timed out")
            raise Exception("Request timed out. The server took too long to respond.")
            
        except Exception as e:
            logger.error(f"Error in _make_request: {str(e)}\n{traceback.format_exc()}")
            raise

    async def close(self):
        """Close the aiohttp session."""
        async with self._session_lock:
            if self._session and not self._session.closed:
                await self._session.close()
                self._session = None
