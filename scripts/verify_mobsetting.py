import glob
import json
import os
import sys

def verify_mobsetting():
    input_dir = 'data/MobSetting/Action'
    json_path = 'data/MobSetting/mob_actions.json'
    
    if not os.path.exists(json_path):
        print(f"Error: {json_path} does not exist")
        sys.exit(1)
        
    with open(json_path, 'r', encoding='utf-8') as f:
        data = json.load(f)
        
    files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    errors = 0
    total_verified = 0
    
    for f in files:
        mob_name = os.path.splitext(os.path.basename(f))[0]
        if mob_name not in data:
            print(f"Error: {mob_name} missing from JSON!")
            errors += 1
            continue
            
        json_records = data[mob_name]
        txt_records = []
        with open(f, 'r', encoding='latin1') as fp:
            for line in fp:
                line = line.strip()
                if not line or line.startswith(';'):
                    continue
                if line.startswith('#Record'):
                    parts = line.split()
                    txt_records.append({
                        "condition_id": int(parts[1]),
                        "action_id": int(parts[2])
                    })
                    
        if len(json_records) != len(txt_records):
            print(f"Count mismatch for {mob_name}: TXT={len(txt_records)}, JSON={len(json_records)}")
            errors += 1
            continue
            
        for i, (tr, jr) in enumerate(zip(txt_records, json_records)):
            if tr != jr:
                print(f"Discrepancy in {mob_name}[{i}]: TXT={tr} vs JSON={jr}")
                errors += 1
            else:
                total_verified += 1
                
    if errors == 0:
        print(f"VERIFICATION SUCCESSFUL: {len(files)} files and {total_verified} records matched with 0 discrepancies!")
    else:
        print(f"VERIFICATION FAILED: {errors} errors found!")
        sys.exit(1)

if __name__ == '__main__':
    verify_mobsetting()
