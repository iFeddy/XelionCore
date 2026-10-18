import os
import glob
import re
import json

base_dir = "./data/MobAttackSequence" if os.path.exists("./data/MobAttackSequence") else "../AstraServer/Server/9Data/Shine/MobAttackSequence"
output_path = "./data/mob_attack_sequence.json"

files = sorted(glob.glob(f"{base_dir}/*.txt"))

data = {}

total_seq_records = 0
total_change_records = 0

for f in files:
    mob_name = os.path.splitext(os.path.basename(f))[0]
    with open(f, "r", encoding="cp949", errors="replace") as fp:
        lines = fp.readlines()

    current_table = None
    seq = []
    changes = []

    for line in lines:
        line = line.strip()
        if not line or line.startswith(";"):
            continue

        # Check table
        m_tab = re.match(r"#Table\s+([a-zA-Z0-9_]+)", line, re.I)
        if m_tab:
            current_table = m_tab.group(1).lower()
            continue

        if line.lower().startswith("#record"):
            # Strip inline comments (e.g. ; comment)
            if ";" in line:
                line = line.split(";")[0].strip()
            
            parts = line.split()[1:] # drop #record
            if not parts:
                continue

            if current_table == "attseq":
                if len(parts) >= 2:
                    order_val = int(parts[0]) if parts[0].isdigit() else parts[0]
                    seq.append({
                        "order": order_val,
                        "attack": parts[1]
                    })
                    total_seq_records += 1
            elif current_table == "skillchange":
                if len(parts) >= 4:
                    val = int(parts[1]) if parts[1].isdigit() else parts[1]
                    from_val = int(parts[2]) if parts[2].isdigit() else parts[2]
                    changes.append({
                        "condition": parts[0],
                        "value": val,
                        "from": from_val,
                        "to": parts[3]
                    })
                    total_change_records += 1
                elif len(parts) == 3:
                    # Rare variation: condition, value, to
                    val = int(parts[1]) if parts[1].isdigit() else parts[1]
                    changes.append({
                        "condition": parts[0],
                        "value": val,
                        "from": "-",
                        "to": parts[2]
                    })
                    total_change_records += 1

    entry = {
        "sequence": seq
    }
    if changes:
        entry["skill_changes"] = changes

    data[mob_name] = entry

with open(output_path, "w", encoding="utf-8") as out:
    json.dump(data, out, indent=2, ensure_ascii=False)

file_size = os.path.getsize(output_path)
print(f"Generated {output_path} successfully!")
print(f"Total Mobs: {len(data)}/297")
print(f"Total Attack Sequence Records: {total_seq_records}")
print(f"Total Skill Change Rules: {total_change_records}")
print(f"File Size: {file_size} bytes ({file_size / 1024:.2f} KB)")
