"""
Deep Search Agent implementation.
This agent is responsible for performing deep analysis of projects and search results.
"""

import logging
import json
from typing import Dict, List, Any, Optional
from services.agents.base import BaseAgent

logger = logging.getLogger(__name__)

class DeepSearchAgent(BaseAgent):
    """Agent responsible for performing deep analysis of projects and search results"""
    
    def __init__(self, go_client):
        """
        Initialize the deep search agent.
        
        Args:
            go_client: Client for communicating with the Go backend
        """
        super().__init__(
            name="Deep Search Agent",
            role="Perform comprehensive analysis of projects and search results"
        )
        self.go_client = go_client
        self.system_prompt = """
        You are an expert at deep technical analysis of software projects.
        Your task is to analyze projects in depth, identify relationships between projects,
        and provide comprehensive insights beyond basic search results.
        
        Focus on:
        - Architecture patterns and design choices
        - Technical compatibility and integration possibilities
        - Implementation approaches and trade-offs
        - Related projects and technologies
        - Common patterns across similar projects
        - Unique or innovative aspects
        - Potential use cases and applications
        """
    
    async def analyze_results(self, query: str, results: List[Dict], context: List[Dict] = None, access_token: str = None, deep_analysis: bool = True) -> Dict:
        """
        Perform deep analysis on search results.
        
        Args:
            query: The original search query
            results: List of search results to analyze
            context: Previous conversation history
            access_token: JWT token for authentication
            deep_analysis: Whether to perform full deep analysis
            
        Returns:
            Dict containing analysis results
        """
        if not results:
            return {"analysis": "No results to analyze"}
            
        # Format results for LLM consumption
        results_text = json.dumps(results[:5], indent=2)  # Limit to top 5 for LLM context
        
        # Format conversation context if available
        context_str = ""
        if context and len(context) > 0:
            context_str = "Previous conversation:\n"
            for msg in context:
                role = msg.get("role", "unknown")
                content = msg.get("content", "")
                context_str += f"{role}: {content}\n"
        
        prompt = f"""Analyze these search results for the query: "{query}"

{context_str}

Results:
{results_text}

Provide a comprehensive analysis including:
1. Common patterns across these projects
2. Technical compatibility considerations
3. Implementation approaches
4. Suggested related areas to explore
5. Unique or innovative aspects
6. Potential use cases and applications

Respond in JSON format:
{{
  "patterns": ["List of common patterns identified"],
  "compatibility": {{
    "technical_stack": "Analysis of technical compatibility",
    "integration": "Integration possibilities",
    "challenges": ["Potential challenges or considerations"]
  }},
  "implementation": {{
    "approaches": ["Different implementation approaches identified"],
    "trade_offs": ["Key trade-offs to consider"],
    "best_practices": ["Recommended best practices"]
  }},
  "innovation": {{
    "unique_aspects": ["Unique or innovative features"],
    "novel_approaches": ["Novel implementation approaches"]
  }},
  "applications": {{
    "use_cases": ["Potential use cases"],
    "industries": ["Relevant industries or domains"]
  }},
  "related_areas": ["Suggested related areas to explore"],
  "summary": "Overall analysis summary"
}}
"""
        
        response = await self.call_llm(self.system_prompt, prompt, temperature=0.5)
        analysis = self.parse_json_response(response)
        
        logger.info(f"Deep analysis completed for query: '{query}'")
        return analysis
    
    def format_analysis(self, analysis: Dict) -> str:
        """
        Format deep analysis results for display.
        
        Args:
            analysis: Dict containing analysis results
            
        Returns:
            Formatted string for display
        """
        if "error" in analysis or not analysis:
            return "Unable to perform deep analysis at this time."
            
        formatted = "# Deep Analysis Results 🔍\n\n"
        
        # Add summary if available
        if "summary" in analysis:
            formatted += f"## Overview\n{analysis['summary']}\n\n"
        
        # Add patterns section
        if "patterns" in analysis and analysis["patterns"]:
            formatted += "## Common Patterns\n"
            for pattern in analysis["patterns"]:
                formatted += f"- {pattern}\n"
            formatted += "\n"
        
        # Add technical compatibility section
        if "compatibility" in analysis:
            compat = analysis["compatibility"]
            formatted += "## Technical Considerations\n"
            if "technical_stack" in compat:
                formatted += f"**Stack Compatibility**: {compat['technical_stack']}\n\n"
            if "integration" in compat:
                formatted += f"**Integration**: {compat['integration']}\n\n"
            if "challenges" in compat:
                formatted += "**Challenges to Consider**:\n"
                for challenge in compat["challenges"]:
                    formatted += f"- {challenge}\n"
            formatted += "\n"
        
        # Add implementation approaches
        if "implementation" in analysis:
            impl = analysis["implementation"]
            formatted += "## Implementation Insights\n"
            if "approaches" in impl:
                formatted += "**Approaches**:\n"
                for approach in impl["approaches"]:
                    formatted += f"- {approach}\n"
            if "best_practices" in impl:
                formatted += "\n**Best Practices**:\n"
                for practice in impl["best_practices"]:
                    formatted += f"- {practice}\n"
            formatted += "\n"
        
        # Add innovation section
        if "innovation" in analysis:
            innov = analysis["innovation"]
            formatted += "## Innovation & Unique Aspects\n"
            if "unique_aspects" in innov:
                for aspect in innov["unique_aspects"]:
                    formatted += f"- {aspect}\n"
            formatted += "\n"
        
        # Add applications section
        if "applications" in analysis:
            apps = analysis["applications"]
            formatted += "## Applications & Use Cases\n"
            if "use_cases" in apps:
                formatted += "**Potential Use Cases**:\n"
                for use_case in apps["use_cases"]:
                    formatted += f"- {use_case}\n"
            if "industries" in apps:
                formatted += "\n**Relevant Industries**:\n"
                for industry in apps["industries"]:
                    formatted += f"- {industry}\n"
            formatted += "\n"
        
        # Add related areas
        if "related_areas" in analysis and analysis["related_areas"]:
            formatted += "## Related Areas to Explore\n"
            for area in analysis["related_areas"]:
                formatted += f"- {area}\n"
        
        return formatted
