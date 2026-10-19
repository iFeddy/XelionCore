from trace_login import symbols

print("=== ShinePlayer Relevant Methods ===")
for va, name in sorted(symbols.items()):
    if 'ShinePlayer' in name and any(k in name for k in ['Store', 'so_char', 'so_Save', 'so_ply']):
        print(f"0x{va:x}: {name.split('@')[0]}")
