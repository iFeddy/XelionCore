import os
import glob
import re
import sys

def parse_doc(doc_path):
    with open(doc_path) as f:
        text = f.read()

    sections = re.split(r"\n### `", text)[1:]
    structs = {}
    for s in sections:
        name = s.split("`")[0]
        if name.startswith("std::") or "OPTOOL" in name:
            continue
        sz_m = re.search(r"Tamaño Total en Memoria:\*\* `(\d+) bytes`", s)
        if not sz_m:
            continue
        total_size = int(sz_m.group(1))

        rows = []
        for line in s.split("\n"):
            line = line.strip()
            if line.startswith("### ") or line.startswith("## "):
                break
            if not line.startswith("| `0x"):
                continue
            cols = [c.strip().strip("`") for c in line.split("|")[1:-1]]
            if len(cols) >= 4:
                hex_off, dec_off, type_col, name_col = cols[0], cols[1], cols[2], cols[3]
                try:
                    dec_val = int(dec_off)
                    rows.append((dec_val, type_col, name_col))
                except ValueError:
                    pass
        structs[name] = (total_size, rows)
    return structs

# Read opcodes
with open("pkg/protocol/opcodes.go") as f:
    opcodes_content = f.read()
opcodes = set(re.findall(r"(\bPROTO_NC_[A-Z0-9_]+)\s*(?:uint16)?\s*=", opcodes_content))

# Collect all types declared in original files
generated_filenames = {
    "gamelog_structs.go",
    "zone_gamelog_registry.go",
    "zone_structs_act.go",
    "zone_structs_bat.go",
    "zone_structs_briefinfo.go",
    "zone_structs_casino.go",
    "zone_structs_char.go",
    "zone_structs_guild.go",
    "zone_structs_item.go",
    "zone_structs_map.go",
    "zone_structs_misc.go",
    "zone_structs_shine.go",
    "zone_structs_subsystems.go",
    "zone_structs_unnamed.go",
}

existing_types = set()
for path in glob.glob("pkg/protocol/*.go"):
    if os.path.basename(path) in generated_filenames:
        continue
    with open(path) as f:
        c = f.read()
    for m in re.finditer(r"type\s+([A-Za-z0-9_]+)\s+struct\b", c):
        existing_types.add(m.group(1))

print(f"Collected {len(existing_types)} existing types from base protocol files.")

RESERVED_WORDS = {
    "break", "default", "func", "interface", "select",
    "case", "defer", "go", "map", "struct",
    "chan", "else", "goto", "package", "switch",
    "const", "fallthrough", "if", "range", "type",
    "continue", "for", "import", "return", "var",
    "string", "int", "uint", "byte", "rune", "float32", "float64",
}

def clean_field_name(name, seen):
    clean = re.sub(r"[^A-Za-z0-9_]", "", name)
    if not clean:
        clean = "Field"
    if clean[0].isdigit():
        clean = "Field_" + clean
    pascal = clean[0].upper() + clean[1:]
    if pascal.lower() in RESERVED_WORDS or pascal in RESERVED_WORDS:
        pascal = pascal + "Field"
    base = pascal
    cnt = 1
    while pascal in seen:
        cnt += 1
        pascal = f"{base}_{cnt}"
    seen.add(pascal)
    return pascal

def resolve_field_type(cpp_type, field_size):
    t = cpp_type.strip()
    t = re.sub(r"/\*.*?\*/", "", t).strip()
    if "[0]" in t:
        return "[0]byte"
    arr_m = re.search(r"\[(\d+)\]", t)
    if arr_m:
        count = int(arr_m.group(1))
        if "char" in t and field_size == count:
            return f"[{count}]byte"
        elif "short" in t and field_size == count * 2:
            return f"[{count}]uint16"
        elif "int" in t and field_size == count * 4:
            return f"[{count}]uint32"
        elif field_size > 0:
            return f"[{field_size}]byte"
    if t in ("unsigned char", "uint8", "bool") and field_size == 1:
        return "uint8"
    if t in ("char", "int8") and field_size == 1:
        return "int8"
    if t in ("unsigned short", "uint16", "NETCOMMAND") and field_size == 2:
        return "uint16"
    if t in ("short", "int16") and field_size == 2:
        return "int16"
    if t in ("unsigned int", "unsigned long", "uint32") and field_size == 4:
        return "uint32"
    if t in ("int", "long", "int32") and field_size == 4:
        return "int32"
    if t in ("unsigned __int64", "uint64", "SHINE_ITEM_REGISTNUMBER") and field_size == 8:
        return "uint64"
    if t in ("__int64", "int64") and field_size == 8:
        return "int64"
    if t == "float" and field_size == 4:
        return "float32"
    if t == "double" and field_size == 8:
        return "float64"
    if t == "NETPACKETHEADER" and field_size == 2:
        return "NETPACKETHEADER"
    if t == "NETPACKETZONEHEADER" and field_size == 6:
        return "NETPACKETZONEHEADER"
    if t == "Name3" and field_size == 12:
        return "Name3"
    if t == "Name4" and field_size == 16:
        return "Name4"
    if t == "Name5" and field_size == 20:
        return "Name5"
    if t == "Name8" and field_size == 32:
        return "Name8"
    if t == "Name256Byte" and field_size == 256:
        return "Name256Byte"
    if t == "SHINE_XY_TYPE" and field_size == 8:
        return "SHINE_XY_TYPE"
    if t == "MAPPOS" and field_size == 20:
        return "MAPPOS"
    if t == "ITEM_INVEN" and field_size == 2:
        return "ITEM_INVEN"
    if field_size == 1:
        return "uint8"
    elif field_size == 2:
        return "uint16"
    elif field_size == 4:
        return "uint32"
    elif field_size == 8:
        return "uint64"
    elif field_size > 0:
        return f"[{field_size}]byte"
    return "[0]byte"

def sanitize_type_name(name):
    clean = name.replace("::", "_").replace("<", "_").replace(">", "_").replace("-", "_")
    clean = re.sub(r"_+", "_", clean).strip("_")
    return clean

def classify_struct(name, is_gamelog):
    if is_gamelog:
        return "gamelog_structs.go"
    if "::" in name:
        return "zone_structs_unnamed.go"
    if name.startswith("ShineSendType"):
        return "zone_structs_shine.go"
    if "PROTO_NC_ACT_" in name:
        return "zone_structs_act.go"
    if "PROTO_NC_BAT_" in name:
        return "zone_structs_bat.go"
    if "PROTO_NC_BRIEFINFO_" in name:
        return "zone_structs_briefinfo.go"
    if "PROTO_NC_CHAR" in name:
        return "zone_structs_char.go"
    if "PROTO_NC_ITEM" in name:
        return "zone_structs_item.go"
    if "PROTO_NC_MAP_" in name:
        return "zone_structs_map.go"
    if "PROTO_NC_MISC_" in name or "PROTO_NC_USER_" in name:
        return "zone_structs_misc.go"
    if "PROTO_NC_GUILD" in name:
        return "zone_structs_guild.go"
    if "PROTO_NC_DICE" in name or "PROTO_NC_GAMBLE" in name:
        return "zone_structs_casino.go"
    return "zone_structs_subsystems.go"

def main():
    gl_structs = parse_doc("docs/Servicios/GameLog.md")
    zn_structs = parse_doc("docs/Servicios/Zone.md")

    file_buckets = {}
    registry_entries = []

    # 1. Process GameLog structs
    for name, (total_size, rows) in sorted(gl_structs.items()):
        go_type = sanitize_type_name(name)
        if name in opcodes:
            go_type += "_DATA"
        target_file = classify_struct(name, is_gamelog=True)
        file_buckets.setdefault(target_file, []).append((name, go_type, total_size, rows, False))
        registry_entries.append((name, go_type, total_size))

    # 2. Process Zone structs
    for name, (total_size, rows) in sorted(zn_structs.items()):
        # If already registered as GameLog struct, skip generating again
        if name in gl_structs:
            continue

        base_go_type = sanitize_type_name(name)
        candidate_data_type = base_go_type + "_DATA"

        # Check if already declared in existing files
        is_already_declared = False
        go_type = base_go_type

        if candidate_data_type in existing_types:
            is_already_declared = True
            go_type = candidate_data_type
        elif base_go_type in existing_types:
            is_already_declared = True
            go_type = base_go_type
        else:
            if name in opcodes:
                go_type = candidate_data_type
            else:
                go_type = base_go_type

        target_file = classify_struct(name, is_gamelog=False)
        file_buckets.setdefault(target_file, []).append((name, go_type, total_size, rows, is_already_declared))
        registry_entries.append((name, go_type, total_size))

    total_structs = sum(len(v) for v in file_buckets.values())
    declared_count = sum(sum(1 for item in v if not item[4]) for v in file_buckets.values())
    reused_count = sum(sum(1 for item in v if item[4]) for v in file_buckets.values())
    print(f"Total struct definitions mapped: {total_structs} (New: {declared_count}, Reused from existing: {reused_count})")

    for target_file, items in sorted(file_buckets.items()):
        target_path = os.path.join("pkg/protocol", target_file)
        lines = []
        lines.append("package protocol")
        lines.append("")
        if target_file == "gamelog_structs.go":
            lines.append("import (")
            lines.append('\t"bytes"')
            lines.append('\t"encoding/binary"')
            lines.append(")")
            lines.append("")

        lines.append("// ============================================================================")
        lines.append(f"// Estructuras Oficiales C++ (#pragma pack(push, 1)) — {target_file}")
        lines.append("// Generadas automáticamente desde Zone.pdb y GameLog.md con paridad exacta de sizeof.")
        lines.append("// ============================================================================")
        lines.append("")

        emitted_count = 0
        for orig_name, go_name, total_size, rows, already_decl in items:
            lines.append(f"// {orig_name} (sizeof = {total_size} bytes)")
            if already_decl:
                lines.append(f"// Declarada previamente en el paquete protocol ({go_name})")
                lines.append("")
                continue

            emitted_count += 1
            fixed_rows = [r for r in rows if "[0]" not in r[1]]
            if not fixed_rows:
                lines.append(f"type {go_name} struct {{}}")
                lines.append("")
                lines.append(f"func (s *{go_name}) FixedWireSize() int {{ return {total_size} }}")
                lines.append("")
                continue

            lines.append(f"type {go_name} struct {{")
            seen_fields = set()
            for idx, (off, cpp_type, fname) in enumerate(fixed_rows):
                if idx + 1 < len(fixed_rows):
                    f_size = fixed_rows[idx+1][0] - off
                else:
                    f_size = total_size - off
                go_field = clean_field_name(fname, seen_fields)
                go_type_str = resolve_field_type(cpp_type, f_size)
                lines.append(f"\t{go_field:<32} {go_type_str:<20} // offset: {off}, size: {f_size} bytes")

            flex_rows = [r for r in rows if "[0]" in r[1]]
            for off, cpp_type, fname in flex_rows:
                go_field = clean_field_name(fname, seen_fields)
                lines.append(f"\t{go_field:<32} [0]byte              // offset: {off}, flexible array member")

            lines.append("}")
            lines.append("")
            lines.append(f"func (s *{go_name}) FixedWireSize() int {{ return {total_size} }}")
            lines.append("")

            if target_file == "gamelog_structs.go":
                lines.append(f"func (s *{go_name}) Serialize() ([]byte, error) {{")
                lines.append("\tbuf := new(bytes.Buffer)")
                lines.append("\tif err := binary.Write(buf, binary.LittleEndian, s); err != nil {")
                lines.append("\t\treturn nil, err")
                lines.append("\t}")
                lines.append("\treturn buf.Bytes(), nil")
                lines.append("}")
                lines.append("")
                lines.append(f"func (s *{go_name}) Deserialize(data []byte) error {{")
                lines.append("\tbuf := bytes.NewReader(data)")
                lines.append("\treturn binary.Read(buf, binary.LittleEndian, s)")
                lines.append("}")
                lines.append("")

        with open(target_path, "w") as f:
            f.write("\n".join(lines))
        print(f"Written {target_path} ({len(items)} structs, {emitted_count} emitted)")

    # Generate zone_gamelog_registry.go
    reg_path = "pkg/protocol/zone_gamelog_registry.go"
    reg_lines = [
        "package protocol",
        "",
        "// ZoneGameLogPDBStructs contiene una instancia de cada estructura oficial C++",
        "// de Zone.md y GameLog.md, indexada por su nombre oficial exacto del PDB TPI.",
        "var ZoneGameLogPDBStructs = map[string]interface{}{",
    ]
    seen_reg = set()
    for orig_name, go_name, total_size in sorted(registry_entries):
        # Register original PDB name
        if orig_name not in seen_reg:
            reg_lines.append(f'\t"{orig_name}": {go_name}{{}},')
            seen_reg.add(orig_name)
        # If go_name differs (e.g. _DATA or sanitized unnamed), register it too
        if go_name != orig_name and go_name not in seen_reg:
            reg_lines.append(f'\t"{go_name}": {go_name}{{}},')
            seen_reg.add(go_name)

    reg_lines.append("}")
    reg_lines.append("")
    with open(reg_path, "w") as f:
        f.write("\n".join(reg_lines))
    print(f"Written {reg_path} ({len(registry_entries)} official entries, {len(seen_reg)} total keys)")

if __name__ == "__main__":
    main()
