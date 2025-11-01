import numpy as np
from sentence_transformers import SentenceTransformer
import logging

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class EmbeddingService:
    """Service for generating text embeddings using sentence-transformers."""
    
    def __init__(self, model_name='all-MiniLM-L6-v2'):
        """
        Initialize the embedding service with the specified model.
        
        Args:
            model_name (str): Name of the sentence-transformers model to use.
                Default is 'all-MiniLM-L6-v2' which provides a good balance
                between performance and quality.
        """
        logger.info(f"Loading embedding model: {model_name}")
        self.model = SentenceTransformer(model_name)
        self.dimensions = self.model.get_sentence_embedding_dimension()
        logger.info(f"Model loaded successfully. Embedding dimensions: {self.dimensions}")
    
    def generate_embedding(self, text):
        """
        Generate an embedding vector for the provided text.
        
        Args:
            text (str): The text to embed.
            
        Returns:
            numpy.ndarray: The embedding vector.
        """
        if not text or not isinstance(text, str):
            logger.warning(f"Invalid text provided: {type(text)}")
            raise ValueError("Text must be a non-empty string")
        
        try:
            logger.debug(f"Generating embedding for text of length {len(text)}")
            embedding = self.model.encode(text)
            return embedding
        except Exception as e:
            logger.error(f"Error generating embedding: {str(e)}")
            raise
    
    def generate_batch_embeddings(self, texts):
        """
        Generate embeddings for a batch of texts.
        
        Args:
            texts (list): List of strings to embed.
            
        Returns:
            list: List of embedding vectors.
        """
        if not texts or not isinstance(texts, list):
            logger.warning(f"Invalid texts provided: {type(texts)}")
            raise ValueError("Texts must be a non-empty list of strings")
        
        try:
            logger.debug(f"Generating embeddings for {len(texts)} texts")
            embeddings = self.model.encode(texts)
            return embeddings
        except Exception as e:
            logger.error(f"Error generating batch embeddings: {str(e)}")
            raise
    
    def get_dimensions(self):
        """
        Get the dimensions of the embedding vectors.
        
        Returns:
            int: Number of dimensions in the embedding vectors.
        """
        return self.dimensions
