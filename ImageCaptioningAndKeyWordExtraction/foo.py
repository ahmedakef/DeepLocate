import sys
import json

if __name__ == "__main__":
    if len(sys.argv) > 1:
        file_name = sys.argv[1]
        print(json.dumps({"Name": "ahmed", "FileName": file_name}))
