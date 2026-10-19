import glob
import re
import os

def parse_go_opcodes():
    go_opcodes = {}
    with open('pkg/protocol/opcodes.go', 'r', encoding='utf-8') as f:
        content = f.read()
    
    for m in re.finditer(r'(PROTO_NC_[A-Z0-9_]+)\s*(?:uint16)?\s*=\s*(?:\((\d+)\s*<<\s*10\)|(\d+)\s*<<\s*10)\s*\|\s*(\d+)', content):
        name = m.group(1)
        dept = int(m.group(2) if m.group(2) is not None else m.group(3))
        cmd = int(m.group(4))
        opcode = (dept << 10) | cmd
        go_opcodes[opcode] = (name, dept, cmd)
    return go_opcodes

def parse_docs_dispatchers():
    doc_opcodes = {}
    for doc in sorted(glob.glob('docs/Servicios/*.md')):
        service_name = os.path.basename(doc).replace('.md', '')
        with open(doc, 'r', encoding='utf-8', errors='ignore') as f:
            for line in f:
                line = line.strip()
                if not line.startswith('| `0x') and not line.startswith('| `0X'):
                    continue
                parts = [p.strip().strip('`') for p in line.split('|')[1:-1]]
                if len(parts) >= 6:
                    hex_str = parts[0]
                    try:
                        op = int(hex_str, 16)
                    except ValueError:
                        continue
                    dept = op >> 10
                    cmd = op & 0x3FF
                    cb_func = parts[4]
                    proto_struct = parts[5]
                    direction = parts[6] if len(parts) > 6 else 'UNK'
                    
                    if op not in doc_opcodes:
                        doc_opcodes[op] = []
                    doc_opcodes[op].append({
                        'service': service_name,
                        'dept': dept,
                        'cmd': cmd,
                        'cb': cb_func,
                        'proto': proto_struct,
                        'dir': direction
                    })
    return doc_opcodes

def main():
    go_ops = parse_go_opcodes()
    doc_ops = parse_docs_dispatchers()
    
    print(f"Total opcodes in pkg/protocol/opcodes.go: {len(go_ops)}")
    print(f"Total opcodes documented in official service dispatchers: {len(doc_ops)}")
    
    missing_ops = {}
    for op, list_info in sorted(doc_ops.items()):
        if op not in go_ops:
            missing_ops[op] = list_info
            
    print(f"Opcodes in official dispatchers MISSING in opcodes.go: {len(missing_ops)}")
    print("=" * 70)
    
    by_service = {}
    by_dept = {}
    
    for op, entries in missing_ops.items():
        for e in entries:
            srv = e['service']
            by_service.setdefault(srv, set()).add(op)
        dept = op >> 10
        by_dept.setdefault(dept, []).append((op, entries))
        
    print("\nMissing opcodes count by service:")
    for srv, ops in sorted(by_service.items()):
        print(f"  - {srv}: {len(ops)} opcodes faltantes")
        
    dept_names = {
        2: "NC_MISC (Utilidades, Tiempos, Claves)",
        3: "NC_USER (Autenticación y Cuentas)",
        4: "NC_CHAR (Selección y Datos Base del Avatar)",
        5: "NC_CHAR Management (Creación / Eliminación)",
        6: "NC_MAP (Traspasos y Navegación de Mapas)",
        7: "NC_BRIEFINFO (Spawns, Visibilidad de Entidades)",
        8: "NC_ACT (Movimiento, Acciones, Interacciones)",
        9: "NC_BAT (Combate, Habilidades, Estados Anormales)",
        10: "NC_PARTY (Grupos y Partys)",
        11: "NC_FELLOW (Compañeros / Mascotas)",
        12: "NC_ITEM (Inventario, Uso, Drops, Comercio)",
        13: "NC_STORE (Tiendas Personales / Mercado)",
        14: "NC_ORDER (Órdenes / Subastas)",
        15: "NC_OP (Operaciones / GM / Administración)",
        16: "NC_COMM (Comunicaciones / Mails)",
        17: "NC_GUILD (Gremios y Academias)",
        18: "NC_QUEST (Misiones y Progresión)",
        19: "NC_FRIEND (Lista de Amigos)",
        20: "NC_MINIHOUSE (Mini-Casas)",
        21: "NC_CHAT / Social",
        22: "NC_KQ (Kingdom Quests)",
        28: "NC_SETTING (Configuraciones y QuickBar)",
        31: "NC_SYSTEM (Sistema y Errores)"
    }
    
    print("\nResumen detallado por Departamento (Department / Head):")
    for dept, items in sorted(by_dept.items()):
        dname = dept_names.get(dept, f"Departamento {dept}")
        print(f"\n--- Dept {dept}: {dname} ({len(items)} faltantes) ---")
        for op, entries in items[:8]:
            protos = "/".join(set(e['proto'] for e in entries))
            cbs = "/".join(set(e['cb'] for e in entries))
            cmd = op & 0x3FF
            print(f"  0x{op:04X} (Head: {dept:2d}, Tail: {cmd:3d}): {protos} [{cbs}]")
        if len(items) > 8:
            print(f"  ... y {len(items)-8} más en este departamento.")

if __name__ == '__main__':
    main()
