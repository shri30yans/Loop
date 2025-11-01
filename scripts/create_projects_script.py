import os
import json
import requests
import random
import time

# Define input directory
input_dir = "dataset"

# API Endpoint
api_url = "http://localhost:8080/api/project/create"
#auth_token = "your_jwt_token_here"  # Replace with your actual token

auth_tokens_list = [
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiN2QwYjE4YzktNGExNy00ZDlhLTllMGQtMmIxNGFiNjUzZjczIiwiZXhwIjoxNzU0MDgxMjA3LCJpYXQiOjE3NDU0NDEyMDd9.v2zKSdMAtr8U7ZbmsuT2s1LjcgjCl9B3shx_V6KMlMw",
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYzdjM2RkODAtNzM0NC00NWViLTkxZmQtZTU5OWMyNDI1OWFhIiwiZXhwIjoxNzU0MDgxMjY1LCJpYXQiOjE3NDU0NDEyNjV9.V7MfE0or7aksYmOqXnazee1TL8rTZb2iw1F6C4tPmU0",
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTgzYzdlMTMtYzQxYi00ZDM3LWIzYmItNDhiZDAxZWY2MjFhIiwiZXhwIjoxNzU0MDgxMjg4LCJpYXQiOjE3NDU0NDEyODh9.zfRbst4M3bGa-02cV1BqC0o3e3WQqNc9CwoYi6IhhwU",
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiY2EyMGU3NzgtNWUwZC00OTE2LWIwNTUtOWNkYTc5M2UzODkwIiwiZXhwIjoxNzU0MDgxMzA5LCJpYXQiOjE3NDU0NDEzMDl9.hjM__Kk7mnoowR0ckvXbPeWVA0IYMQVtm7s-gai6Ag8",
"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzcxNWI0YmEtMzg5ZC00NWM0LWEzYTUtZmVhODdkMzczMmI3IiwiZXhwIjoxNzU0MDgxMzI2LCJpYXQiOjE3NDU0NDEzMjZ9.u1KrGZ8HrtXqbL5pM_Jgu4VSykQ207iwTsF_OOxPKxQ",
]

auth_token = random.choice(auth_tokens_list)
    
# Headers
headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {auth_token}"
}

# Process all JSON files
for filename in os.listdir(input_dir):
    time.sleep(5)
    
    if filename.endswith(".json"):
        input_filepath = os.path.join(input_dir, filename)

        try:
            # Read JSON file
            with open(input_filepath, "r", encoding="utf-8") as infile:
                content = infile.read()
                content = content.replace("body","content")
                json_data = json.loads(content)

            # Send POST request
            response = requests.post(api_url, json=json_data, headers=headers)

            # Print response
            print(f"Response for {filename}: {response.status_code} - {response.text}")

        except json.JSONDecodeError as e:
            print(f"JSON syntax error in {filename}: {e}")
        except Exception as e:
            print(f"Error processing {filename}: {e}")

print("All requests completed.")
