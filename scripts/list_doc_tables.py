import glob
import os

for f in sorted(glob.glob('docs/Servicios/*.md')):
    with open(f, 'r', encoding='utf-8', errors='ignore') as fp:
        lines = fp.readlines()
    table_lines = [l for l in lines if l.strip().startswith('| `0x') or l.strip().startswith('| `0X')]
    print(f"{os.path.basename(f)}: {len(table_lines)} opcode table rows")
