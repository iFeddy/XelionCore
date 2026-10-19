import glob
import os
import re

def parse_all_docs():
    all_opcodes = {}
    for doc in sorted(glob.glob('docs/Servicios/*.md')):
        fname = os.path.basename(doc)
        with open(doc, 'r', encoding='utf-8', errors='ignore') as f:
            for line in f:
                line = line.strip()
                if not (line.startswith('| `0x') or line.startswith('| `0X')):
                    continue
                parts = [p.strip().strip('`') for p in line.split('|')[1:-1]]
                if len(parts) >= 5:
                    hex_str = parts[0]
                    try:
                        op = int(hex_str, 16)
                    except ValueError:
                        continue
                    if op == 0 or op > 0xFFFF:
                        continue
                    dept = op >> 10
                    cmd = op & 0x3FF
                    
                    # Look for PROTO_NC_* or NC_*
                    proto_name = ""
                    for p in parts[3:]:
                        if 'PROTO_NC_' in p:
                            proto_name = p
                            break
                        elif p.startswith('sp_NC_') or p.startswith('fc_NC_') or p.startswith('wms_NC_') or p.startswith('wc_NC_') or p.startswith('cs_NC_') or p.startswith('as_NC_'):
                            prefix = p.split('_')[0]
                            proto_name = 'PROTO_' + p[len(prefix)+1:]
                            break
                        elif p.startswith('NC_'):
                            proto_name = 'PROTO_' + p
                            break
                    
                    if not proto_name and len(parts) >= 5:
                        proto_name = parts[4]
                    
                    if op not in all_opcodes:
                        all_opcodes[op] = []
                    all_opcodes[op].append((fname, dept, cmd, proto_name, parts))
    return all_opcodes

def main():
    all_ops = parse_all_docs()
    print(f"Total unique opcodes collected across all docs: {len(all_ops)}")
    
    proto_named = sum(1 for op, entries in all_ops.items() if any('PROTO_NC_' in e[3] for e in entries))
    print(f"Opcodes with official PROTO_NC_ name: {proto_named}")

if __name__ == '__main__':
    main()
