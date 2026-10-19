import struct

def main():
    with open('pdo/WorldManager.pdb', 'rb') as f:
        pdb = f.read()
    with open('../NA2016/Server/WorldManager/WorldManager.exe', 'rb') as f:
        exe = f.read()

    symbols = {}
    pos = 0
    while True:
        idx = pdb.find(b'\x0e\x11', pos)
        if idx == -1: break
        if idx >= 2:
            reclen = struct.unpack_from('<H', pdb, idx-2)[0]
            if reclen >= 10:
                flags, off, seg = struct.unpack_from('<IIH', pdb, idx+2)
                name_end = pdb.find(b'\x00', idx+12)
                if name_end != -1 and name_end <= idx + reclen:
                    try:
                        name = pdb[idx+12:name_end].decode('latin1')
                        if seg == 1:
                            symbols[0x400000 + 0x1000 + off] = name
                    except: pass
        pos = idx + 2

    target_va = 0x427260
    target_rva = target_va - 0x400000

    records = []
    for i in range(len(exe) - 0x400 - 5):
        if exe[0x400 + i] == 0xe8:
            rel = struct.unpack_from('<i', exe, 0x400 + i + 1)[0]
            target = (0x1000 + i + 5 + rel) & 0xFFFFFFFF
            if target == target_rva:
                call_raw = 0x400 + i
                call_rva = 0x1000 + i
                
                funcs = sorted([r for r in symbols if r <= 0x400000 + call_rva])
                enc_func = symbols.get(funcs[-1], 'unknown') if funcs else 'unknown'
                clean_enc = enc_func.split('@')[0].replace('?', '').replace('0', '')
                
                chunk = exe[call_raw-25:call_raw]
                p = len(chunk) - 5
                while p >= 0:
                    if chunk[p] == 0x68:
                        fva = struct.unpack_from('<I', chunk, p+1)[0]
                        if 0x401000 <= fva <= 0x490000:
                            cur = p + 5
                            cmd = None
                            if cur < len(chunk) and chunk[cur] == 0x6a:
                                cmd = chunk[cur+1]; cur += 2
                            elif cur + 5 <= len(chunk) and chunk[cur] == 0x68:
                                cmd = struct.unpack_from('<I', chunk, cur+1)[0]; cur += 5
                            dept = None
                            if cur < len(chunk) and chunk[cur] == 0x6a:
                                dept = chunk[cur+1]; cur += 2
                            elif cur + 5 <= len(chunk) and chunk[cur] == 0x68:
                                dept = struct.unpack_from('<I', chunk, cur+1)[0]; cur += 5
                            if cmd is not None and dept is not None:
                                opcode = (dept << 10) | cmd
                                sym = symbols.get(fva, f'sub_{hex(fva)}')
                                clean_sym = sym.split('@')[0].replace('?', '')
                                records.append((clean_enc, opcode, dept, cmd, clean_sym, fva))
                                break
                    p -= 1

    # Group by parser
    by_parser = {}
    for enc, op, dept, cmd, sym, fva in records:
        by_parser.setdefault(enc, {})[op] = (dept, cmd, sym, fva)

    lines = []
    lines.append("## Dispatcher de Protocolos Oficial (WorldManager.exe)\n")
    lines.append(f"Mapeo de ingeniería inversa de las tablas de despacho `PROTOCOLFUNCTIONTEMPLETE` del binario oficial `WorldManager.exe` contra `WorldManager.pdb`.\n")
    lines.append(f"> **Total de Protocolos Oficiales Registrados:** `{len(records)}` distribuidos en `{len(by_parser)}` subsistemas/parsers.\n")

    parser_names = [
        ("CParserClient", "1. Protocolos Cliente -> WorldManager (C2S - Selección de Personajes, Gremios, Chat)", "C2S"),
        ("CParserCharDB", "2. Protocolos Character Server -> WorldManager (S2S - Persistencia de Personajes)", "S2S"),
        ("CParserZone", "3. Protocolos Zone Server -> WorldManager (S2S - Coordinación de Mapas)", "S2S"),
        ("CParserLogin", "4. Protocolos Login Server -> WorldManager (S2S - Transferencia de Sesiones)", "S2S"),
        ("CParserAccountDB", "5. Protocolos Account Server -> WorldManager (S2S)", "S2S"),
        ("CParserAccountLogDB", "6. Protocolos AccountLog Server -> WorldManager (S2S)", "S2S"),
        ("CParserGameLogDB", "7. Protocolos GameLog Server -> WorldManager (S2S)", "S2S"),
        ("CParserOPTool", "8. Protocolos Herramientas Administrativas / GM (OPTool)", "S2S"),
    ]

    for pkey, ptitle, pdir in parser_names:
        if pkey in by_parser:
            dict_ops = by_parser[pkey]
            lines.append(f"### {ptitle}\n")
            lines.append(f"> **Total:** `{len(dict_ops)}` handlers oficiales.\n")
            lines.append("| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Función Callback (WorldManager.exe) | Estructura Oficial | Dirección |")
            lines.append("| :--- | :--- | :--- | :--- | :--- | :--- | :--- |")
            for op in sorted(dict_ops.keys()):
                dept, cmd, sym, fva = dict_ops[op]
                sname = sym.replace('fc_', 'PROTO_') if sym.startswith('fc_') else sym
                lines.append(f"| `0x{op:04X}` | {op:5d} | {dept:2d} | {cmd:3d} | `{sym}` | `{sname}` | {pdir} |")
            lines.append("")

    content = "\n".join(lines) + "\n\n"

    with open('docs/Servicios/WorldManager.md', 'r', encoding='utf-8', errors='ignore') as f:
        wm_md = f.read()

    marker = "## Catálogo Exhaustivo de Estructuras de Paquetes"
    if marker in wm_md:
        parts = wm_md.split(marker, 1)
        new_wm_md = parts[0] + content + marker + parts[1]
    else:
        new_wm_md = wm_md + "\n\n" + content

    with open('docs/Servicios/WorldManager.md', 'w', encoding='utf-8') as f:
        f.write(new_wm_md)

    print(f"Successfully integrated {len(records)} official WorldManager protocols into docs/Servicios/WorldManager.md")

if __name__ == '__main__':
    main()
