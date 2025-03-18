import os
import json
import requests

# Define input directory
input_dir = "dataset"

# API Endpoint
api_url = "http://localhost:8080/api/project/create"
auth_token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDU1ZjBjMjctZDczZC00Y2M2LTg3Y2QtYjUxZTczZTQ1YTM0IiwiZXhwIjoxNzUwNTczNTkwLCJpYXQiOjE3NDE5MzM1OTB9.8Ei2joAzBgTPj3f6-_Jx82KjqDc2EZMSkv_x1K7cORU"  # Replace with your actual token

# Headers
headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {auth_token}"
}

# Process all JSON files
for filename in os.listdir(input_dir):
    if filename.endswith(".json"):
        input_filepath = os.path.join(input_dir, filename)

        try:
            # Read JSON file
            with open(input_filepath, "r", encoding="utf-8") as infile:
                json_data = json.load(infile)

            # Send POST request
            response = requests.post(api_url, json=json_data, headers=headers)

            # Print response
            print(f"Response for {filename}: {response.status_code} - {response.text}")

        except json.JSONDecodeError as e:
            print(f"JSON syntax error in {filename}: {e}")
        except Exception as e:
            print(f"Error processing {filename}: {e}")

print("All requests completed.")