"""
Services package initialization.
This package contains core services for the Loop agent system.
"""

from services.embedding import EmbeddingService
from services.go_client import GoBackendClient

__all__ = [
    'EmbeddingService',
    'GoBackendClient',
]

# Package-level logger configuration
import logging

logging.getLogger(__name__).addHandler(logging.NullHandler())
