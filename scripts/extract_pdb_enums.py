import struct
import math
import sys
import os
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from parse_pdb_types import MSFReader, parse_numeric

def get_enums(pdb_path):
    msf = MSFReader(pdb_path)
    tpi = msf.get_stream(2)
    ver, hdr_size, ti_begin, ti_end, rec_bytes = struct.unpack_from("<IIIII", tpi, 0)
    
    type_records = {}
    pos = hdr_size
    ti = ti_begin
    while pos < hdr_size + rec_bytes:
        rec_len = struct.unpack_from("<H", tpi, pos)[0]
        leaf = struct.unpack_from("<H", tpi, pos + 2)[0]
        type_records[ti] = (leaf, tpi[pos+4 : pos+2+rec_len])
        ti += 1
        pos += 2 + rec_len
        
    enums = {}
    for ti, (leaf, rdata) in type_records.items():
        if leaf == 0x1507: # LF_ENUM
            fcount, prop, utype, field_list = struct.unpack_from("<HHII", rdata, 0)
            name = rdata[12:].split(b"\x00")[0].decode("latin1", errors="ignore")
            # Parse field_list (LF_FIELDLIST = 0x1203)
            fields = []
            if field_list in type_records and type_records[field_list][0] == 0x1203:
                fdata = type_records[field_list][1]
                fpos = 0
                while fpos < len(fdata):
                    fleaf = struct.unpack_from("<H", fdata, fpos)[0]
                    if fleaf == 0x1502: # LF_ENUMERATE
                        attr = struct.unpack_from("<H", fdata, fpos+2)[0]
                        val, vpos = parse_numeric(fdata, fpos+4)
                        fname = fdata[vpos:].split(b"\x00")[0].decode("latin1", errors="ignore")
                        fields.append((fname, val))
                        fpos = vpos + len(fname.encode("latin1")) + 1
                        # align to 4 bytes
                        fpos = (fpos + 3) & ~3
                    elif fleaf in (0x1512, 0x1513): # pad
                        fpos += 2
                    else:
                        break
            enums[name] = fields
    return enums

if __name__ == '__main__':
    pdb_path = sys.argv[1] if len(sys.argv) > 1 else 'pdo/Zone.pdb'
    enums = get_enums(pdb_path)
    print(f"Total enums in {pdb_path}: {len(enums)}")
    nc_enums = {k: v for k, v in enums.items() if 'NC_' in k or any('NC_' in f[0] for f in v)}
    print(f"Enums with NC_: {len(nc_enums)}")
    for name, fields in nc_enums.items():
        print(f"Enum {name}: {len(fields)} values")
        for fn, fv in fields[:10]:
            print(f"  {fn} = {fv} (0x{fv:X})")
        if len(fields) > 10:
            print(f"  ... and {len(fields)-10} more")
