import struct

def main():
    with open('pdo/Zone.pdb', 'rb') as f:
        pdb = f.read()
    with open('../NA2016/Server/Zone00/Zone.exe', 'rb') as f:
        exe = f.read()

    symbols = {}
    pos = 0
    while True:
        idx = pdb.find(b'\x0e\x11', pos)
        if idx == -1:
            break
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
                    except:
                        pass
        pos = idx + 2

    raw_start = 0x400 + 0xd3510
    raw_end = 0xd48c0
    data = exe[raw_start:raw_end]

    records = []
    p = 0
    while p < len(data) - 15:
        if data[p] == 0x68:
            func_va = struct.unpack_from('<I', data, p+1)[0]
            cur = p + 5
            cmd = None
            if data[cur] == 0x6a:
                cmd = data[cur+1]
                cur += 2
            elif data[cur] == 0x68:
                cmd = struct.unpack_from('<I', data, cur+1)[0]
                cur += 5
            dept = None
            if cur < len(data) and data[cur] == 0x6a:
                dept = data[cur+1]
                cur += 2
            elif cur < len(data) and data[cur] == 0x68:
                dept = struct.unpack_from('<I', data, cur+1)[0]
                cur += 5

            if cmd is not None and dept is not None and cur + 7 <= len(data):
                if data[cur:cur+2] == b'\x8b\xce' and data[cur+2] == 0xe8:
                    opcode = (dept << 10) | cmd
                    sym = symbols.get(func_va, f'unk_{hex(func_va)}')
                    clean_sym = sym.split('@')[0].replace('?', '')
                    records.append((opcode, dept, cmd, clean_sym, func_va))
                    p = cur + 7
                    continue
        p += 1

    unique = {}
    for op, dept, cmd, sym, va in records:
        if op not in unique:
            unique[op] = (dept, cmd, sym, va)

    md_lines = []
    md_lines.append("## Dispatcher de Protocolos Oficial (C2S - Zone.exe)\n")
    md_lines.append("Mapeo exacto de la tabla de despacho `protocolstore(ShinePlayer)` (`0x004D3510`) del binario oficial `Zone.exe` contra los símbolos de `Zone.pdb`.\n")
    md_lines.append(f"> **Total de Opcodes Oficiales Registrados:** `{len(unique)}`\n")
    md_lines.append("| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Función de Callback (Zone.exe) | Estructura Oficial | Dirección |")
    md_lines.append("| :--- | :--- | :--- | :--- | :--- | :--- | :--- |")

    for op in sorted(unique.keys()):
        dept, cmd, sym, va = unique[op]
        struct_name = sym.replace('sp_', 'PROTO_') if sym.startswith('sp_') else sym
        md_lines.append(f"| `0x{op:04X}` | {op:5d} | {dept:2d} | {cmd:3d} | `{sym}` | `{struct_name}` | C2S |")

    content = "\n".join(md_lines) + "\n\n"

    # Now read Zone.md and insert before "## Catálogo Exhaustivo"
    with open('docs/Servicios/Zone.md', 'r', encoding='utf-8', errors='ignore') as f:
        zone_md = f.read()

    marker = "## Catálogo Exhaustivo de Estructuras de Paquetes"
    if marker in zone_md:
        parts = zone_md.split(marker, 1)
        new_zone_md = parts[0] + content + marker + parts[1]
    else:
        # Try finding Cat
        idx = zone_md.find("## Cat")
        if idx != -1:
            new_zone_md = zone_md[:idx] + content + zone_md[idx:]
        else:
            new_zone_md = zone_md + "\n\n" + content

    with open('docs/Servicios/Zone.md', 'w', encoding='utf-8') as f:
        f.write(new_zone_md)

    print(f"Successfully integrated {len(unique)} official Zone protocols into docs/Servicios/Zone.md")

if __name__ == '__main__':
    main()
