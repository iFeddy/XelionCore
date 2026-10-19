from parse_pdb_types import PDBTypeParser
p = PDBTypeParser('pdo/Zone.pdb')

structs_to_check = [
    'PROTO_NC_CHAR_BASE_CMD',
    'PROTO_NC_CHAR_LOOK_CMD',
    'PROTO_NC_QUEST_BUSY_LIST_CMD',
    'PROTO_NC_QUEST_DONE_LIST_CMD',
    'PROTO_NC_SKILL_ACTIVE_LIST_CMD',
    'PROTO_NC_SKILL_PASSIVE_LIST_CMD',
    'PROTO_NC_ITEM_LIST_CMD',
    'PROTO_NC_TITLE_LIST_CMD',
    'PROTO_NC_CHAR_INFO_END_CMD',
    'PROTO_NC_CHAR_MAPLOGIN_ACK'
]

for name in structs_to_check:
    info = p.structs.get(name)
    if info:
        print(f"=== {name} (Size: {info['size']}) ===")
        members = p.get_struct_members(name)
        if members:
            for m in members:
                print(f"  offset {m['offset']:3d}: {m['type']} {m['name']}")
        else:
            print("  (no members or empty)")
    else:
        # Search partial matches
        matches = [k for k in p.structs.keys() if name in k]
        print(f"=== {name}: NOT FOUND. Possible matches: {matches} ===")
