import struct
import math
import os
import sys

PRIMITIVE_TYPES = {
    0x0003: "void",
    0x0010: "char",
    0x0020: "unsigned char",
    0x0011: "short",
    0x0021: "unsigned short",
    0x0012: "int",
    0x0022: "unsigned int",
    0x0013: "__int64",
    0x0023: "unsigned __int64",
    0x0040: "float",
    0x0041: "double",
    0x0070: "char",
    0x0071: "wchar_t",
    0x0030: "bool",
    0x0200: "int",
    0x0074: "int",
    0x0075: "unsigned int",
}

class MSFReader:
    def __init__(self, filename):
        with open(filename, "rb") as f: self.data = f.read()
        self.block_size, _, _, self.num_dir_bytes, _ = struct.unpack_from("<IIIII", self.data, 0x20)
        self.block_map_addr = struct.unpack_from("<I", self.data, 0x34)[0]
        num_dir_blocks = math.ceil(self.num_dir_bytes / self.block_size)
        dir_blocks = struct.unpack_from(f"<{num_dir_blocks}I", self.data, self.block_map_addr * self.block_size)
        dir_data = bytearray()
        for b in dir_blocks: dir_data += self.data[b*self.block_size : (b+1)*self.block_size]
        dir_data = dir_data[:self.num_dir_bytes]
        self.num_streams = struct.unpack_from("<I", dir_data, 0)[0]
        self.stream_sizes = struct.unpack_from(f"<{self.num_streams}i", dir_data, 4)
        pos = 4 + self.num_streams * 4
        self.stream_blocks = []
        for size in self.stream_sizes:
            if size <= 0: self.stream_blocks.append([]); continue
            n = math.ceil(size / self.block_size)
            self.stream_blocks.append(struct.unpack_from(f"<{n}I", dir_data, pos))
            pos += n * 4
    def get_stream(self, idx):
        if idx >= self.num_streams or self.stream_sizes[idx] <= 0: return b""
        res = bytearray()
        for b in self.stream_blocks[idx]: res += self.data[b*self.block_size : (b+1)*self.block_size]
        return bytes(res[:self.stream_sizes[idx]])

def parse_numeric(buf, offset):
    tag = struct.unpack_from("<H", buf, offset)[0]
    if tag < 0x8000: return tag, offset + 2
    elif tag == 0x8000: return struct.unpack_from("<b", buf, offset+2)[0], offset + 3
    elif tag == 0x8001: return struct.unpack_from("<h", buf, offset+2)[0], offset + 4
    elif tag == 0x8002: return struct.unpack_from("<H", buf, offset+2)[0], offset + 4
    elif tag == 0x8003: return struct.unpack_from("<i", buf, offset+2)[0], offset + 6
    elif tag == 0x8004: return struct.unpack_from("<I", buf, offset+2)[0], offset + 6
    return tag, offset + 2

class PDBParser:
    def __init__(self, pdb_path):
        self.msf = MSFReader(pdb_path)
        tpi = self.msf.get_stream(2)
        ver, hdr_size, self.ti_begin, self.ti_end, rec_bytes = struct.unpack_from("<IIIII", tpi, 0)
        self.type_records = {}
        pos = hdr_size
        ti = self.ti_begin
        while pos < hdr_size + rec_bytes:
            rec_len = struct.unpack_from("<H", tpi, pos)[0]
            leaf = struct.unpack_from("<H", tpi, pos + 2)[0]
            self.type_records[ti] = (leaf, tpi[pos+4 : pos+2+rec_len])
            ti += 1
            pos += 2 + rec_len

        self.structs = {}
        for ti, (leaf, rdata) in self.type_records.items():
            if leaf in (0x1505, 0x1504): # LF_STRUCTURE, LF_CLASS
                fcount, prop, field_list = struct.unpack_from("<HHI", rdata, 0)
                sz_val, sz_pos = parse_numeric(rdata, 16)
                name = rdata[sz_pos:].split(b"\x00")[0].decode("latin1", errors="ignore")
                is_fwd = bool(prop & 0x80)
                if not is_fwd or name not in self.structs:
                    self.structs[name] = {"ti": ti, "size": sz_val, "field_list_ti": field_list, "fcount": fcount, "leaf": leaf, "is_fwd": is_fwd}
            elif leaf == 0x1506: # LF_UNION
                fcount, prop, field_list = struct.unpack_from("<HHI", rdata, 0)
                sz_val, sz_pos = parse_numeric(rdata, 8)
                name = rdata[sz_pos:].split(b"\x00")[0].decode("latin1", errors="ignore")
                is_fwd = bool(prop & 0x80)
                if not is_fwd or name not in self.structs:
                    self.structs[name] = {"ti": ti, "size": sz_val, "field_list_ti": field_list, "fcount": fcount, "leaf": leaf, "is_fwd": is_fwd}

    def resolve_type(self, ti):
        if ti in PRIMITIVE_TYPES: return PRIMITIVE_TYPES[ti]
        if ti not in self.type_records: return f"type_0x{ti:x}"
        leaf, rdata = self.type_records[ti]
        if leaf == 0x1001:
            mod_ti, attr = struct.unpack_from("<IH", rdata, 0)
            prefix = "const " if (attr & 1) else ""
            return f"{prefix}{self.resolve_type(mod_ti)}"
        elif leaf in (0x1002, 0x1500):
            target_ti = struct.unpack_from("<I", rdata, 0)[0]
            return f"{self.resolve_type(target_ti)}*"
        elif leaf == 0x1503: # LF_ARRAY
            elem_ti, idx_ti = struct.unpack_from("<II", rdata, 0)
            elem_name = self.resolve_type(elem_ti)
            sz, _ = parse_numeric(rdata, 8)
            elem_size = 1
            if elem_name in ("char", "unsigned char", "bool"): elem_size = 1
            elif elem_name in ("short", "unsigned short"): elem_size = 2
            elif elem_name in ("int", "unsigned int", "float", "long", "unsigned long"): elem_size = 4
            elif elem_name in ("__int64", "unsigned __int64", "double"): elem_size = 8
            elif elem_name in self.structs and self.structs[elem_name]["size"] > 0:
                elem_size = self.structs[elem_name]["size"]
            count = sz // elem_size if elem_size > 0 else sz
            return f"{elem_name}[{count}]"
        elif leaf in (0x1505, 0x1504, 0x1506):
            sz_pos = 8 if leaf == 0x1506 else 16
            sz_val, sz_pos = parse_numeric(rdata, sz_pos)
            sname = rdata[sz_pos:].split(b"\x00")[0].decode("latin1", errors="ignore")
            # Special friendly representation for known unions
            if sname == "Name256Byte": return "char[256] /* Name256Byte */"
            if sname == "Name4": return "char[16] /* Name4 */"
            if sname == "Name18Byte": return "char[18] /* Name18Byte */"
            if sname == "Name50Byte": return "char[50] /* Name50Byte */"
            if sname == "Name60Byte": return "char[60] /* Name60Byte */"
            if sname == "ORToken": return "char[20] /* ORToken */"
            if sname == "NETCOMMAND": return "NETCOMMAND /* union (command:10, dept:6) */"
            return sname
        elif leaf == 0x1507:
            ename = rdata[12:].split(b"\x00")[0].decode("latin1", errors="ignore")
            return f"enum {ename}"
        elif leaf == 0x1206:
            base_ti, length, pos = struct.unpack_from("<IBB", rdata, 0)
            return f"{self.resolve_type(base_ti)} : {length}"
        return f"type_0x{leaf:x}"

    def get_members(self, sname):
        if sname not in self.structs: return []
        flist_ti = self.structs[sname]["field_list_ti"]
        if flist_ti not in self.type_records: return []
        leaf, rdata = self.type_records[flist_ti]
        if leaf != 0x1203: return []
        members = []
        pos = 0
        while pos < len(rdata):
            if rdata[pos] >= 0xf0:
                pos += 1; continue
            sleaf = struct.unpack_from("<H", rdata, pos)[0]
            if sleaf == 0x150d: # LF_MEMBER
                attr, mtype_ti = struct.unpack_from("<HI", rdata, pos + 2)
                offset, cur = parse_numeric(rdata, pos + 8)
                mname = rdata[cur:].split(b"\x00")[0].decode("latin1", errors="ignore")
                cur += len(mname) + 1
                members.append({
                    "name": mname,
                    "type": self.resolve_type(mtype_ti),
                    "offset": offset
                })
                pos = cur
            elif sleaf == 0x150e: # LF_STMEMBER
                cur = pos + 8
                mname = rdata[cur:].split(b"\x00")[0].decode("latin1", errors="ignore")
                pos = cur + len(mname) + 1
            elif sleaf == 0x1400: # LF_BCLASS
                attr, btype_ti = struct.unpack_from("<HI", rdata, pos + 2)
                offset, cur = parse_numeric(rdata, pos + 8)
                bname = self.resolve_type(btype_ti)
                members.append({
                    "name": f"/* base */ {bname}",
                    "type": bname,
                    "offset": offset
                })
                pos = cur
            elif sleaf in (0x150b, 0x1510):
                cur = pos + (8 if sleaf == 0x1510 else 8)
                mname = rdata[cur:].split(b"\x00")[0]
                pos = cur + len(mname) + 1
            elif sleaf == 0x150c:
                cur = pos + 6
                tname = rdata[cur:].split(b"\x00")[0]
                pos = cur + len(tname) + 1
            else:
                break
        return members

def format_struct_markdown(sname, sinfo, members, opcode_info=""):
    md = []
    size = sinfo["size"]
    md.append(f"### `{sname}`\n")
    if opcode_info:
        md.append(f"- **OpCode Asociado:** `{opcode_info}`")
    md.append(f"- **Tamaño Total en Memoria:** `{size} bytes` (`0x{size:X}`)\n")
    
    # C++ declaration
    md.append("```cpp")
    md.append("#pragma pack(push, 1)")
    md.append(f"struct {sname} {{")
    if not members:
        md.append("    // Estructura vacía / Payload de 0 bytes")
    for m in members:
        # Check if type is array like char[256]
        t = m["type"]
        n = m["name"]
        if "[" in t:
            base, arr = t.split("[", 1)
            arr = "[" + arr
            # clean comment
            comment = ""
            if "/*" in arr:
                arr, comment = arr.split("/*", 1)
                arr = arr.strip()
                comment = " /*" + comment
            md.append(f"    {base.strip()} {n}{arr};{comment}")
        else:
            md.append(f"    {t} {n};")
    md.append("};")
    md.append("#pragma pack(pop)")
    md.append("```\n")

    # Field table
    if members:
        md.append("| Offset (Hex) | Offset (Dec) | Tipo de Dato | Campo / Variable | Descripción Oficial |")
        md.append("| :--- | :--- | :--- | :--- | :--- |")
        for m in members:
            off_h = f"0x{m['offset']:03X}"
            off_d = f"{m['offset']}"
            t = m['type']
            n = m['name']
            desc = "Campo de paquete oficial extraído de PDB TPI"
            if "netpacketheader" in n.lower(): desc = "Cabecera interna de sesión (ClientHandle uint16)"
            elif "netcmd" in n.lower() or "protocol" in n.lower(): desc = "OpCode del protocolo (Dept:6 bits, Cmd:10 bits)"
            elif "size" == n.lower(): desc = "Longitud del payload del paquete en bytes"
            elif "user" in n.lower() and "name" in n.lower(): desc = "Nombre de cuenta o usuario (string)"
            elif "password" in n.lower(): desc = "Contraseña en texto plano o hash del cliente"
            elif "ip" in n.lower(): desc = "Dirección IP del cliente o del servidor"
            elif "port" in n.lower(): desc = "Puerto TCP de escucha"
            elif "world" in n.lower(): desc = "Identificador o información del WorldServer"
            elif "block" in n.lower(): desc = "Estado de bloqueo / baneo de cuenta o IP"
            elif "auth" in n.lower(): desc = "Nivel de autorización (Player, GM, Admin)"
            elif "otp" in n.lower(): desc = "One-Time Password token"
            elif "token" in n.lower(): desc = "Token de autenticación externa"
            md.append(f"| `{off_h}` | `{off_d}` | `{t}` | `{n}` | {desc} |")
        md.append("")
    else:
        md.append("> *Paquete de señalización sin campos de datos adicionales (0 bytes de payload exclusivo).*\n")
    return "\n".join(md)

print("PDB Parser module loaded successfully.")
