"""
Agents package initialization.
This package contains specialized agents for project discovery and analysis.
"""

from services.agents.base import BaseAgent
from services.agents.project_search import ProjectSearchAgent
from services.agents.clarification import ClarificationAgent
from services.agents.deep_search import DeepSearchAgent

__all__ = [
    'BaseAgent',
    'ProjectSearchAgent',
    'ClarificationAgent',
    'DeepSearchAgent',
]

# Package-level logger configuration
import logging

logging.getLogger(__name__).addHandler(logging.NullHandler())
