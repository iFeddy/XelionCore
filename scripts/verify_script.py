import glob
import json
import os
import sys

def verify_script():
    input_dir = 'data/Script'
    txt_files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    
    errors = 0
    total_verified = 0

    for txt_path in txt_files:
        base_name = os.path.splitext(os.path.basename(txt_path))[0]
        json_path = os.path.join(input_dir, f"{base_name}.json")

        if not os.path.exists(json_path):
            print(f"Missing JSON: {json_path}")
            errors += 1
            continue

        with open(json_path, 'r', encoding='utf-8') as f:
            json_data = json.load(f)

        txt_entries = {}
        with open(txt_path, 'rb') as fp:
            for line in fp:
                try:
                    line_str = line.decode('cp949').rstrip('\r\n')
                except UnicodeDecodeError:
                    line_str = line.decode('latin1').rstrip('\r\n')

                if not line_str.startswith('#Record'):
                    continue

                parts = line_str.split('\t')
                tokens = [p.strip() for p in parts if p.strip()]
                if len(tokens) < 2:
                    continue

                key = tokens[1]
                if key.startswith(';') or not key:
                    continue

                val = tokens[2] if len(tokens) > 2 else ""
                if val.startswith(';'):
                    val = ""

                if len(val) >= 2 and val.startswith('"') and val.endswith('"'):
                    val = val[1:-1]

                txt_entries[key] = val

        for key, val in txt_entries.items():
            if key not in json_data:
                print(f"Missing key in {base_name}.json: {key}")
                errors += 1
            elif json_data[key] != val:
                print(f"Mismatch in {base_name}.json [{key}]: TXT='{val}' vs JSON='{json_data[key]}'")
                errors += 1
            else:
                total_verified += 1

    if errors == 0:
        print(f"VERIFICATION SUCCESSFUL: {len(txt_files)} files and {total_verified} entries verified with 0 discrepancies!")
    else:
        print(f"VERIFICATION FAILED: {errors} errors found!")
        sys.exit(1)

if __name__ == '__main__':
    verify_script()
