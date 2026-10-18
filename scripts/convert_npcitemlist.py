import glob
import json
import os

def convert_npcitemlist():
    input_dir = 'data/NPCItemList'
    files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    print(f"Found {len(files)} merchant text files in {input_dir}...")
    
    total_tabs = 0
    total_rows = 0
    total_items = 0
    
    for f in files:
        merchant_name = os.path.splitext(os.path.basename(f))[0]
        tabs = []
        cur_tab = None
        
        with open(f, 'rb') as fp:
            for line in fp:
                try:
                    line_str = line.decode('cp949').strip()
                except UnicodeDecodeError:
                    line_str = line.decode('latin1').strip()
                    
                if not line_str:
                    continue
                    
                if line_str.startswith('#Table'):
                    # e.g.: #Table Tab00 ;Comment
                    parts = line_str.split(';', 1)
                    table_tokens = parts[0].split()
                    tab_name = table_tokens[1] if len(table_tokens) > 1 else 'Tab00'
                    title = parts[1].strip() if len(parts) > 1 else ''
                    
                    cur_tab = {
                        "tab": tab_name,
                        "title": title,
                        "rows": []
                    }
                    tabs.append(cur_tab)
                    total_tabs += 1
                    
                elif line_str.startswith('#Record') and cur_tab is not None:
                    parts = line_str.split()
                    # parts[0] is '#Record', parts[1] is row index (Rec), parts[2..7] are 6 items
                    items = parts[2:]
                    cur_tab["rows"].append(items)
                    total_rows += 1
                    for item in items:
                        if item != '-':
                            total_items += 1
                            
        out_data = {
            "merchant": merchant_name,
            "tabs": tabs
        }
        
        out_path = os.path.join(input_dir, f"{merchant_name}.json")
        with open(out_path, 'w', encoding='utf-8') as out_fp:
            json.dump(out_data, out_fp, indent=2, ensure_ascii=False)
            out_fp.write('\n')
            
    print(f"Converted {len(files)} merchants successfully.")
    print(f"Total tabs: {total_tabs}, Total rows: {total_rows}, Total valid items: {total_items}")

if __name__ == '__main__':
    convert_npcitemlist()
