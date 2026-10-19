import glob
import os
import sys
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from extract_pdb_enums import get_enums

all_protocol_enums = {}

for pdb_path in sorted(glob.glob('pdo/*.pdb')):
    pdb_name = os.path.basename(pdb_path)
    enums = get_enums(pdb_path)
    p_enums = {k: v for k, v in enums.items() if k.startswith('PROTOCOL_')}
    print(f"{pdb_name}: {len(p_enums)} PROTOCOL_ enums")
    for k, v in p_enums.items():
        if k not in all_protocol_enums:
            all_protocol_enums[k] = (pdb_name, v)
        else:
            # merge if more values
            if len(v) > len(all_protocol_enums[k][1]):
                all_protocol_enums[k] = (pdb_name, v)

print("=" * 60)
print(f"Total unique PROTOCOL_ enums found: {len(all_protocol_enums)}")
for name, (pdb_src, fields) in sorted(all_protocol_enums.items()):
    print(f"{name:<35} ({len(fields):3d} fields) from {pdb_src}")
