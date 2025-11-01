FROM python:3.9-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

# Expose port
EXPOSE 5000

# Set environment variables
ENV PYTHONUNBUFFERED=1
ENV EMBEDDING_MODEL=all-MiniLM-L6-v2

# Run the application with Uvicorn (ASGI server)
# Using --workers based on CPU cores (2 is a reasonable default)
# Using --lifespan=on to properly handle startup/shutdown events
# Using --http httptools for faster HTTP parsing
CMD ["uvicorn", "app:app", "--host", "0.0.0.0", "--port", "5000", "--workers", "2", "--lifespan", "on", "--http", "httptools", "--ws", "websockets", "--loop", "uvloop"]
