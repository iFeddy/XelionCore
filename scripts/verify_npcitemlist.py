import glob
import json
import os
import sys

def verify_npcitemlist():
    input_dir = 'data/NPCItemList'
    txt_files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    
    errors = 0
    total_verified_tabs = 0
    total_verified_rows = 0
    total_verified_items = 0
    
    for txt_path in txt_files:
        base_name = os.path.splitext(os.path.basename(txt_path))[0]
        json_path = os.path.join(input_dir, f"{base_name}.json")
        
        if not os.path.exists(json_path):
            print(f"Error: Missing JSON file for {base_name}")
            errors += 1
            continue
            
        with open(json_path, 'r', encoding='utf-8') as f:
            json_data = json.load(f)
            
        # Parse TXT
        txt_tabs = []
        cur_tab = None
        with open(txt_path, 'rb') as fp:
            for line in fp:
                try:
                    line_str = line.decode('cp949').strip()
                except UnicodeDecodeError:
                    line_str = line.decode('latin1').strip()
                    
                if not line_str:
                    continue
                    
                if line_str.startswith('#Table'):
                    parts = line_str.split(';', 1)
                    table_tokens = parts[0].split()
                    tab_name = table_tokens[1] if len(table_tokens) > 1 else 'Tab00'
                    title = parts[1].strip() if len(parts) > 1 else ''
                    cur_tab = {
                        "tab": tab_name,
                        "title": title,
                        "rows": []
                    }
                    txt_tabs.append(cur_tab)
                elif line_str.startswith('#Record') and cur_tab is not None:
                    parts = line_str.split()
                    items = parts[2:]
                    cur_tab["rows"].append(items)
                    
        # Compare
        if json_data.get("merchant") != base_name:
            print(f"Merchant name mismatch in {base_name}: {json_data.get('merchant')}")
            errors += 1
            
        json_tabs = json_data.get("tabs", [])
        if len(json_tabs) != len(txt_tabs):
            print(f"Tab count mismatch in {base_name}: TXT={len(txt_tabs)}, JSON={len(json_tabs)}")
            errors += 1
            continue
            
        for t_idx, (t_txt, t_json) in enumerate(zip(txt_tabs, json_tabs)):
            total_verified_tabs += 1
            if t_txt["tab"] != t_json["tab"]:
                print(f"Tab name mismatch in {base_name}[{t_idx}]: {t_txt['tab']} vs {t_json['tab']}")
                errors += 1
            if t_txt["title"] != t_json["title"]:
                print(f"Tab title mismatch in {base_name}[{t_idx}]: '{t_txt['title']}' vs '{t_json['title']}'")
                errors += 1
                
            rows_txt = t_txt["rows"]
            rows_json = t_json.get("rows", [])
            if len(rows_txt) != len(rows_json):
                print(f"Row count mismatch in {base_name} tab {t_txt['tab']}: TXT={len(rows_txt)}, JSON={len(rows_json)}")
                errors += 1
                continue
                
            for r_idx, (r_txt, r_json) in enumerate(zip(rows_txt, rows_json)):
                total_verified_rows += 1
                if len(r_txt) != len(r_json):
                    print(f"Col count mismatch in {base_name} tab {t_txt['tab']} row {r_idx}")
                    errors += 1
                    continue
                for c_idx, (i_txt, i_json) in enumerate(zip(r_txt, r_json)):
                    total_verified_items += 1
                    if i_txt != i_json:
                        print(f"Item mismatch in {base_name} tab {t_txt['tab']} row {r_idx} col {c_idx}: '{i_txt}' vs '{i_json}'")
                        errors += 1

    if errors == 0:
        print(f"VERIFICATION SUCCESSFUL: {len(txt_files)} merchants, {total_verified_tabs} tabs, {total_verified_rows} rows, and {total_verified_items} item slots matched with 0 discrepancies!")
    else:
        print(f"VERIFICATION FAILED: {errors} errors found!")
        sys.exit(1)

if __name__ == '__main__':
    verify_npcitemlist()
