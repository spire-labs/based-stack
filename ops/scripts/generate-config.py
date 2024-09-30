import json
import os

def write_json_files(json_data):
    for path, files_dict in json_data.items():
        for filename, data in files_dict.items():
            full_path = os.path.join(path, f"{filename}.json")
            # Create directories if they don't exist
            os.makedirs(os.path.dirname(full_path), exist_ok=True)
            # Write the data to the file
            with open(full_path, 'w') as f:
                json.dump(data, f, indent=2)

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.abspath(os.path.join(script_dir, '..', '..'))
    # Path to the JSON file in the project root
    json_file_path = os.path.join(project_root, 'based.json')
    with open(json_file_path, 'r') as f:
        json_data = json.load(f)

    write_json_files(json_data)




if __name__ == "__main__":
    try:
       main()
    except Exception as e:
       print(e)
       exit(1)