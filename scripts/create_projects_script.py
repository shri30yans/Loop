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
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzdkZDkzNjktYTM5MC00OThkLTllNWQtZjA3ZWUxZTRmMmNjIiwiZXhwIjoxNzUzOTg4NTY5LCJpYXQiOjE3NDUzNDg1Njl9.XvFvxIWBHAhBw1_eclQJ7RlqrAZgWM9H9YQ1GVpRfwE",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGRmMTIxN2ItZjNiYi00NWU2LWIwYjQtMGE2ZWY5YWRkYjFjIiwiZXhwIjoxNzUzOTg4NTk4LCJpYXQiOjE3NDUzNDg1OTh9.aKyLRuBTNcn98DAd-X9s8sg6hQ3crF7rgcg12q2JYZI",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTdkYWE5OGEtOGNjMy00YzZkLWEyMzctN2JkYzNhN2ZjZDRlIiwiZXhwIjoxNzUzOTg4NjE5LCJpYXQiOjE3NDUzNDg2MTl9.WDDS64AyKF_Ve_BZOczHO8_hxtAYXHc2IsyIRLZSuc0",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGZjNmQ1Y2QtMmU0Yy00ZDYyLTg5MmQtMmZmMTU4ZTNmN2ViIiwiZXhwIjoxNzUzOTg4NjQyLCJpYXQiOjE3NDUzNDg2NDJ9.hShdu3TmH6m8zu7mMbLWvoK11ZyY822Zue2EKLGFPyA",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNGJlNmUyNTEtYzU3Ni00Yzk1LTkxOTctYWIwYTk1ZjFkOWExIiwiZXhwIjoxNzUzOTg4NjYyLCJpYXQiOjE3NDUzNDg2NjJ9.jDuqXcvCA_lNA2vHCFNP2URuDePMaFhv-cjT1N2uDHQ"
]

auth_token = random.choice(auth_tokens_list)
    
# Headers
headers = {
    "Content-Type": "application/json",
    "Authorization": f"Bearer {auth_token}"
}

# Process all JSON files
for filename in os.listdir(input_dir):
    time.sleep(2)
    
    if filename.endswith(".json"):
        input_filepath = os.path.join(input_dir, filename)

        try:
            # Read JSON file
            with open(input_filepath, "r", encoding="utf-8") as infile:
                json_data = json.load(infile)

            json_data.replace("body","content")
            # Send POST request
            response = requests.post(api_url, json=json_data, headers=headers)

            # Print response
            print(f"Response for {filename}: {response.status_code} - {response.text}")

        except json.JSONDecodeError as e:
            print(f"JSON syntax error in {filename}: {e}")
        except Exception as e:
            print(f"Error processing {filename}: {e}")

print("All requests completed.")