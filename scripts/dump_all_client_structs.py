from parse_pdb_types import PDBTypeParser
p = PDBTypeParser('pdo/Zone.pdb')

target_names = [
    'BASE', 'SHAPE', 'QUEST_DOING', 'QUEST_DONE', 'ABSTATE', 'SKILL', 'PASSIVE',
    'HOUSE', 'FRIEND', 'MASPUP', 'GUILD', 'EMBLEM', 'PET', 'ARENA', 'SEAWAR',
    'ITEM', 'GAME', 'CHARTITLE', 'CHARGEDBUFF'
]

for t in target_names:
    found = False
    for name, info in p.structs.items():
        if f'CLIENT_{t}_CMD' in name or (f'CHAR_{t}_CMD' in name and not name.startswith('PROTO_NC_CHARSAVE') and not 'ZONE' in name):
            members = p.get_struct_members(name)
            print(f"=== {name} (size {info['size']}) ===")
            if members:
                for m in members:
                    print(f"  {m['offset']}: {m['type']} {m['name']}")
            found = True
    if not found:
        print(f"*** No struct found for {t} ***")
