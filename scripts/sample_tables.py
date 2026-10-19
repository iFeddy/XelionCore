import glob
import os

for f in sorted(glob.glob('docs/Servicios/*.md')):
    print("=" * 60)
    print(f"File: {os.path.basename(f)}")
    with open(f, 'r', encoding='utf-8', errors='ignore') as fp:
        lines = fp.readlines()
    count = 0
    for l in lines:
        if l.strip().startswith('| `0x') or l.strip().startswith('| `0X'):
            print(l.strip())
            count += 1
            if count >= 3:
                break
