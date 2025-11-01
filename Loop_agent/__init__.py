"""
Loop Agent package initialization.
This package provides an intelligent project discovery system using specialized agents.
"""

from team import ProjectDiscoveryTeam
from services.go_client import GoBackendClient
from services.embedding import EmbeddingService

__version__ = "0.1.0"
__author__ = "Loop Team"
__description__ = "Intelligent project discovery system using specialized agents"

__all__ = [
    'ProjectDiscoveryTeam',
    'GoBackendClient',
    'EmbeddingService',
]

# Package-level logger configuration
import logging

logging.getLogger(__name__).addHandler(logging.NullHandler())
