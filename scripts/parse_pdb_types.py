import struct
import math
import sys

class MSFReader:
    def __init__(self, filename):
        with open(filename, "rb") as f:
            self.data = f.read()
        self.block_size, _, _, self.num_dir_bytes, _ = struct.unpack_from("<IIIII", self.data, 0x20)
        self.block_map_addr = struct.unpack_from("<I", self.data, 0x34)[0]
        num_dir_blocks = math.ceil(self.num_dir_bytes / self.block_size)
        dir_blocks = struct.unpack_from(f"<{num_dir_blocks}I", self.data, self.block_map_addr * self.block_size)
        dir_data = bytearray()
        for b in dir_blocks:
            dir_data += self.data[b*self.block_size : (b+1)*self.block_size]
        dir_data = dir_data[:self.num_dir_bytes]
        self.num_streams = struct.unpack_from("<I", dir_data, 0)[0]
        self.stream_sizes = struct.unpack_from(f"<{self.num_streams}i", dir_data, 4)
        pos = 4 + self.num_streams * 4
        self.stream_blocks = []
        for size in self.stream_sizes:
            if size <= 0:
                self.stream_blocks.append([])
                continue
            n = math.ceil(size / self.block_size)
            self.stream_blocks.append(struct.unpack_from(f"<{n}I", dir_data, pos))
            pos += n * 4

    def get_stream(self, idx):
        if idx >= self.num_streams or self.stream_sizes[idx] <= 0:
            return b""
        res = bytearray()
        for b in self.stream_blocks[idx]:
            res += self.data[b*self.block_size : (b+1)*self.block_size]
        return bytes(res[:self.stream_sizes[idx]])

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

def parse_numeric(buf, offset):
    tag = struct.unpack_from("<H", buf, offset)[0]
    if tag < 0x8000:
        return tag, offset + 2
    elif tag == 0x8000:
        return struct.unpack_from("<b", buf, offset+2)[0], offset + 3
    elif tag == 0x8001:
        return struct.unpack_from("<h", buf, offset+2)[0], offset + 4
    elif tag == 0x8002:
        return struct.unpack_from("<H", buf, offset+2)[0], offset + 4
    elif tag == 0x8003:
        return struct.unpack_from("<i", buf, offset+2)[0], offset + 6
    elif tag == 0x8004:
        return struct.unpack_from("<I", buf, offset+2)[0], offset + 6
    return tag, offset + 2

class PDBTypeParser:
    def __init__(self, pdb_path):
        msf = MSFReader(pdb_path)
        tpi = msf.get_stream(2)
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
                name_bytes = rdata[sz_pos:].split(b"\x00")[0]
                name = name_bytes.decode("latin1", errors="ignore")
                is_fwd = bool(prop & 0x80)
                if not is_fwd or name not in self.structs:
                    self.structs[name] = {
                        "ti": ti,
                        "size": sz_val,
                        "field_list_ti": field_list,
                        "fcount": fcount,
                        "leaf": leaf,
                        "is_fwd": is_fwd
                    }
                elif is_fwd and self.structs[name]["is_fwd"] and sz_val > 0:
                    self.structs[name] = {
                        "ti": ti,
                        "size": sz_val,
                        "field_list_ti": field_list,
                        "fcount": fcount,
                        "leaf": leaf,
                        "is_fwd": is_fwd
                    }
            elif leaf == 0x1506: # LF_UNION
                fcount, prop, field_list = struct.unpack_from("<HHI", rdata, 0)
                sz_val, sz_pos = parse_numeric(rdata, 8)
                name_bytes = rdata[sz_pos:].split(b"\x00")[0]
                name = name_bytes.decode("latin1", errors="ignore")
                is_fwd = bool(prop & 0x80)
                if not is_fwd or name not in self.structs:
                    self.structs[name] = {
                        "ti": ti,
                        "size": sz_val,
                        "field_list_ti": field_list,
                        "fcount": fcount,
                        "leaf": leaf,
                        "is_fwd": is_fwd
                    }

    def resolve_type_str(self, ti):
        if ti in PRIMITIVE_TYPES:
            return PRIMITIVE_TYPES[ti]
        if ti not in self.type_records:
            return f"type_0x{ti:x}"
        leaf, rdata = self.type_records[ti]
        if leaf == 0x1001: # LF_MODIFIER
            mod_ti, attr = struct.unpack_from("<IH", rdata, 0)
            prefix = ""
            if attr & 1: prefix += "const "
            if attr & 2: prefix += "volatile "
            return f"{prefix}{self.resolve_type_str(mod_ti)}"
        elif leaf in (0x1002, 0x1500): # LF_POINTER
            target_ti = struct.unpack_from("<I", rdata, 0)[0]
            return f"{self.resolve_type_str(target_ti)}*"
        elif leaf == 0x1503: # LF_ARRAY
            elem_ti, idx_ti = struct.unpack_from("<II", rdata, 0)
            elem_name = self.resolve_type_str(elem_ti)
            sz, _ = parse_numeric(rdata, 8)
            elem_size = 1
            if elem_name in ("char", "unsigned char", "bool"):
                elem_size = 1
            elif elem_name in ("short", "unsigned short"):
                elem_size = 2
            elif elem_name in ("int", "unsigned int", "float", "long", "unsigned long"):
                elem_size = 4
            elif elem_name in ("__int64", "unsigned __int64", "double"):
                elem_size = 8
            elif elem_name in self.structs and self.structs[elem_name]["size"] > 0:
                elem_size = self.structs[elem_name]["size"]
            count = sz // elem_size if elem_size > 0 else sz
            return f"{elem_name}[{count}]"
        elif leaf in (0x1505, 0x1504): # LF_STRUCTURE, LF_CLASS
            sz_val, sz_pos = parse_numeric(rdata, 16)
            sname = rdata[sz_pos:].split(b"\x00")[0].decode("latin1", errors="ignore")
            return sname
        elif leaf == 0x1506: # LF_UNION
            sz_val, sz_pos = parse_numeric(rdata, 8)
            uname = rdata[sz_pos:].split(b"\x00")[0].decode("latin1", errors="ignore")
            return uname
        elif leaf == 0x1507: # LF_ENUM
            ename = rdata[12:].split(b"\x00")[0].decode("latin1", errors="ignore")
            return f"enum {ename}"
        elif leaf == 0x1206: # LF_BITFIELD
            base_ti, length, pos = struct.unpack_from("<IBB", rdata, 0)
            return f"{self.resolve_type_str(base_ti)} : {length}"
        return f"type_0x{leaf:x}"

    def get_struct_members(self, sname):
        if sname not in self.structs:
            return None
        sinfo = self.structs[sname]
        flist_ti = sinfo["field_list_ti"]
        if flist_ti not in self.type_records:
            return []
        leaf, rdata = self.type_records[flist_ti]
        if leaf != 0x1203:
            return []
        members = []
        pos = 0
        while pos < len(rdata):
            if rdata[pos] >= 0xf0:
                pos += 1
                continue
            sub_leaf = struct.unpack_from("<H", rdata, pos)[0]
            if sub_leaf == 0x150d: # LF_MEMBER
                attr, mtype_ti = struct.unpack_from("<HI", rdata, pos + 2)
                offset, cur = parse_numeric(rdata, pos + 8)
                mname_bytes = rdata[cur:].split(b"\x00")[0]
                mname = mname_bytes.decode("latin1", errors="ignore")
                cur += len(mname_bytes) + 1
                members.append({
                    "name": mname,
                    "type": self.resolve_type_str(mtype_ti),
                    "type_ti": mtype_ti,
                    "offset": offset
                })
                pos = cur
            elif sub_leaf == 0x150e: # LF_STMEMBER
                attr, mtype_ti = struct.unpack_from("<HI", rdata, pos + 2)
                cur = pos + 8
                mname_bytes = rdata[cur:].split(b"\x00")[0]
                mname = mname_bytes.decode("latin1", errors="ignore")
                pos = cur + len(mname_bytes) + 1
            elif sub_leaf in (0x150b, 0x1510):
                if sub_leaf == 0x1510:
                    cur = pos + 8
                    mname_bytes = rdata[cur:].split(b"\x00")[0]
                    pos = cur + len(mname_bytes) + 1
                else:
                    count, mlist = struct.unpack_from("<HI", rdata, pos + 2)
                    cur = pos + 8
                    mname_bytes = rdata[cur:].split(b"\x00")[0]
                    pos = cur + len(mname_bytes) + 1
            elif sub_leaf == 0x150c:
                cur = pos + 6
                tname = rdata[cur:].split(b"\x00")[0]
                pos = cur + len(tname) + 1
            elif sub_leaf == 0x1400:
                attr, btype_ti = struct.unpack_from("<HI", rdata, pos + 2)
                offset, cur = parse_numeric(rdata, pos + 8)
                bname = self.resolve_type_str(btype_ti)
                members.append({
                    "name": f"::base_class_{bname}",
                    "type": bname,
                    "type_ti": btype_ti,
                    "offset": offset
                })
                pos = cur
            else:
                break
        return members
