import struct
import re

def main():
    with open('../NA2016/Server/Character/Character.exe', 'rb') as f:
        exe = f.read()

    target_rva = 0x31160
    records = []
    for i in range(len(exe) - 0x400 - 5):
        if exe[0x400 + i] == 0xe8:
            rel = struct.unpack_from('<i', exe, 0x400 + i + 1)[0]
            target = (0x1000 + i + 5 + rel) & 0xFFFFFFFF
            if target == target_rva:
                call_raw = 0x400 + i
                call_rva = 0x1000 + i
                chunk = exe[call_raw-25:call_raw]
                
                p = len(chunk) - 5
                while p >= 0:
                    if chunk[p] == 0x68:
                        fva = struct.unpack_from('<I', chunk, p+1)[0]
                        if 0x401000 <= fva <= 0x4b0000:
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
                                records.append((call_rva, opcode, dept, cmd, fva))
                                break
                    p -= 1

    # Load struct names from Character.md
    with open('docs/Servicios/Character.md', 'r', encoding='utf-8', errors='ignore') as f:
        char_md = f.read()
    known_structs = set(re.findall(r'### `(PROTO_NC_[A-Za-z0-9_]+)`', char_md))

    # Resolve handler name from function body
    unique = {}
    for call_rva, op, dept, cmd, fva in records:
        if op in unique: continue
        fraw = 0x400 + (fva - 0x401000)
        hname = ""
        # Search for .rdata strings in fraw..fraw+400
        for offset in range(0, min(400, len(exe)-fraw-4)):
            val = struct.unpack_from('<I', exe, fraw + offset)[0]
            if 0x4ad000 <= val < 0x4d6a00:
                s_raw = 0xac000 + (val - 0x4ad000)
                s_data = exe[s_raw:s_raw+120].split(b'\x00')[0]
                if b'NC_' in s_data:
                    try:
                        s = s_data.decode('latin1')
                        # Clean string
                        cand = s.split('::')[-1].split(',')[0].split(' - ')[-1].split(':')[-1].strip()
                        if 'NC_' in cand:
                            m = re.search(r'(?:fc_)?(NC_[A-Za-z0-9_]+)', cand)
                            if m:
                                hname = 'fc_' + m.group(1).replace('fc_', '')
                                break
                    except: pass
        
        # If not found, try matching against known structs
        sname = hname.replace('fc_', 'PROTO_') if hname else ""
        if not sname:
            # Fallback to fva
            hname = f"fc_NC_DEPT{dept}_CMD{cmd}"
            sname = f"PROTO_NC_DEPT{dept}_CMD{cmd}"

        unique[op] = (dept, cmd, hname, sname, fva)

    lines = []
    lines.append("## Dispatcher de Protocolos Oficial (Character.exe)\n")
    lines.append("Mapeo exacto de la tabla de despacho `pft_Store` (`0x00431160`) del binario oficial `Character.exe`.\n")
    lines.append(f"> **Total de Protocolos S2S Registrados:** `{len(unique)}`\n")
    lines.append("| OpCode (Hex) | OpCode (Dec) | Dept (Head) | Cmd (Tail) | Función de Callback (Character.exe) | Estructura Oficial | Dirección |")
    lines.append("| :--- | :--- | :--- | :--- | :--- | :--- | :--- |")

    for op in sorted(unique.keys()):
        dept, cmd, hname, sname, fva = unique[op]
        lines.append(f"| `0x{op:04X}` | {op:5d} | {dept:2d} | {cmd:3d} | `{hname}` | `{sname}` | S2S |")

    content = "\n".join(lines) + "\n\n"

    marker = "## Procedimientos Almacenados SQL Invocados"
    if marker in char_md:
        parts = char_md.split(marker, 1)
        new_char_md = parts[0] + content + marker + parts[1]
    else:
        new_char_md = char_md + "\n\n" + content

    with open('docs/Servicios/Character.md', 'w', encoding='utf-8') as f:
        f.write(new_char_md)

    print(f"Successfully integrated {len(unique)} official Character protocols into docs/Servicios/Character.md")

if __name__ == '__main__':
    main()
