from parse_pdb_types import PDBTypeParser

p = PDBTypeParser(r'pdo/Zone.pdb')

for struct_name in ['SHINE_CHAR_STATVAR', 'CHAR_PARAMETER_DATA::<unnamed-type-PwrStone>', 'PROTO_NC_CHAR_LOOK_CMD']:
    info = p.structs.get(struct_name)
    if info:
        print(f"=== {struct_name} (Size: {info['size']}) ===")
        members = p.get_struct_members(struct_name)
        for m in members:
            print(f"  offset {m['offset']:3d} (0x{m['offset']:02x}): {m['type']} {m['name']}")
