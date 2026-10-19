from parse_pdb_types import PDBTypeParser
p = PDBTypeParser('pdo/Zone.pdb')

print("=== All matching PROTO_NC_ structs ===")
for name in sorted(p.structs.keys()):
    if name.startswith('PROTO_NC_'):
        if any(k in name for k in ['QUEST', 'SKILL', 'ITEM', 'TITLE', 'LOOK', 'END', 'CHAR_']):
            print(f"{name} (size: {p.structs[name]['size']})")
