import os
import json
import requests
import random

# Define input directory
input_dir = "dataset"

# API Endpoint
api_url = "http://localhost:8080/api/project/create"
#auth_token = "your_jwt_token_here"  # Replace with your actual token

auth_tokens_list = [
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZGJkNGFiOGItMDVhMS00NjQ2LTg5OTgtZWZkMjMxZGU4OTg5IiwiZXhwIjoxNzUwOTk1ODY1LCJpYXQiOjE3NDIzNTU4NjV9.auJEGegv99NS6LS_PFgKlzMtgij5ugHDwiJ1BP2SdvU",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYjQ3MDVlNDUtYTBhYS00YWU3LWE1MjItZDhhZmZlMjFkOGVkIiwiZXhwIjoxNzUwOTk1OTMyLCJpYXQiOjE3NDIzNTU5MzJ9.H7o_BJgVhTGhwZ4ALR_zWw-OROyqli9UMMBAc98moGk",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMDQ4YTBkMGMtYWU2YS00NWY0LTk2NTQtYzM0MmNkNDNmNWU1IiwiZXhwIjoxNzUwOTk1OTQ2LCJpYXQiOjE3NDIzNTU5NDZ9.-mLMB34KdOZ5qqk9JewAe0BYADQu5fqnAI-Kcb1zTkg",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMGM4MTljMGQtNTQ4ZS00MDg4LWEyM2UtYzljNWZjNTRhYjdlIiwiZXhwIjoxNzUwOTk1OTU2LCJpYXQiOjE3NDIzNTU5NTZ9.rDl0awHDa3lJO8_NdA_kMD_7R35lueIts3S_i-1rb3c",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMjk4OWYwNjEtMmQyYy00MGNmLTk1YTEtNDZjOWNkNjVjMWMyIiwiZXhwIjoxNzUwOTk1OTY5LCJpYXQiOjE3NDIzNTU5Njl9.Sufz6_FLc0l3xpLeeOKuybv-5WvE72a4R_wxq9QqogE",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTNkN2JhNzktMGI0Yy00OTRlLWJmM2QtYzg0ZWQ5NWE2OTgxIiwiZXhwIjoxNzUwOTk1OTk5LCJpYXQiOjE3NDIzNTU5OTl9.w5JQiQfGWMkcRPaSBOB8EvBiATrRoQDXBrco8G0HJy4",
    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTMxYjY5NDItNDUyOS00MWUxLWFiYTUtODUxNDdmYTFjNjdhIiwiZXhwIjoxNzUwOTk2MDE0LCJpYXQiOjE3NDIzNTYwMTR9.WsKpJpXq1k-3qY2atgnIpaehAVMMp4brq9Sild69cKs",
]

# pop_auth_token_list = []

# while auth_tokens_list:
#     try:
#         auth_token = auth_tokens_list.pop(0)
#         pop_auth_token_list.append(auth_token)
#     except:
#         auth_tokens_list = pop_auth_token_list
#         pop_auth_token_list = []

auth_token = random.choice(auth_tokens_list)
    
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