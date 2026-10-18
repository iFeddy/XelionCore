import glob
import json
import os

def convert_mobsetting():
    input_dir = 'data/MobSetting/Action'
    output_file = 'data/MobSetting/mob_actions.json'
    
    files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    print(f"Processing {len(files)} action files from {input_dir}...")
    
    data = {}
    total_records = 0
    
    for file_path in files:
        mob_name = os.path.splitext(os.path.basename(file_path))[0]
        records = []
        with open(file_path, 'r', encoding='latin1') as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith(';'):
                    continue
                if line.startswith('#Record'):
                    parts = line.split()
                    if len(parts) >= 3:
                        cond_id = int(parts[1])
                        act_id = int(parts[2])
                        records.append({
                            "condition_id": cond_id,
                            "action_id": act_id
                        })
                    else:
                        print(f"Warning: malformed record in {file_path}: {line}")
        data[mob_name] = records
        total_records += len(records)
        
    os.makedirs(os.path.dirname(output_file), exist_ok=True)
    with open(output_file, 'w', encoding='utf-8') as out:
        json.dump(data, out, indent=2, ensure_ascii=False)
        out.write('\n')
        
    print(f"Successfully generated {output_file}")
    print(f"Entities: {len(data)}, Total records: {total_records}")

if __name__ == '__main__':
    convert_mobsetting()
