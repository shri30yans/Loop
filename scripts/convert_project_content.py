import json
import os

def convert_project_content(input_file, output_file):
    with open(input_file, 'r') as infile:
        data = json.load(infile)
    
    # Extract the required fields
    text = f"{data['title']}: {data['introduction']}\n\n"
    for section in data['sections']:
        text += f"{section['title']}: {section['body']}\n\n"
    tags = data['tags']
    
    # Create the new structure
    converted_data = {
        "text": text,
        "tags": tags
    }
    
    # Write to the output file
    with open(output_file, 'w') as outfile:
        json.dump(converted_data, outfile, indent=2)

input_dir = "d:/sem_6/code/genai-main-1/Loop/converted_dataset_temp"
output_dir = "d:/sem_6/code/genai-main-1/Loop/converted_dataset_temp_v2"
os.makedirs(output_dir, exist_ok=True)
for filename in os.listdir(input_dir):
    if filename.endswith(".json"):
        input_filepath = os.path.join(input_dir, filename)
        output_filepath = os.path.join(output_dir, filename)
        # Read JSON file
        with open(input_filepath, "r", encoding="utf-8") as infile:
            json_data = json.load(infile)
    convert_project_content(input_filepath, output_filepath)
