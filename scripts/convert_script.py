import glob
import json
import os

def convert_script_files():
    input_dir = 'data/Script'
    txt_files = sorted(glob.glob(os.path.join(input_dir, '*.txt')))
    print(f"Converting {len(txt_files)} script text files in {input_dir}...")

    total_records = 0

    for file_path in txt_files:
        base_name = os.path.splitext(os.path.basename(file_path))[0]
        entries = {}

        with open(file_path, 'rb') as fp:
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

                # Strip outer quotes if both start and end with quote
                if len(val) >= 2 and val.startswith('"') and val.endswith('"'):
                    val = val[1:-1]

                entries[key] = val
                total_records += 1

        out_path = os.path.join(input_dir, f"{base_name}.json")
        with open(out_path, 'w', encoding='utf-8') as out_fp:
            json.dump(entries, out_fp, indent=2, ensure_ascii=False)
            out_fp.write('\n')

    print(f"Successfully converted {len(txt_files)} files into individual JSON files.")
    print(f"Total script string entries: {total_records}")

if __name__ == '__main__':
    convert_script_files()
