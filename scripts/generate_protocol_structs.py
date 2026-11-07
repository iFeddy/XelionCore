import os
import re

def parse_markdown_structs_with_tables(path):
    structs = {}
    with open(path) as f:
        content = f.read()
    sections = re.split(r"\n###\s+`?(PROTO_[A-Za-z0-9_]+)`?\s*\n", content)
    for i in range(1, len(sections), 2):
        name = sections[i]
        body = sections[i+1]
        if "OPTOOL" in name:
            continue
        size_m = re.search(r"-\s*\*\*Tamaño Total en Memoria:\*\*\s*`?(\d+)\s*bytes`?", body)
        if not size_m:
            continue
        total_size = int(size_m.group(1))
        rows = []
        for line in body.split("\n"):
            line = line.strip()
            if line.startswith("### ") or line.startswith("## "):
                break
            if not line.startswith("|"):
                continue
            cols = [c.strip().strip("`") for c in line.split("|")[1:-1]]
            if len(cols) >= 4 and cols[0] != "Offset (Hex)" and cols[0] != ":---":
                hex_off, dec_off, type_col, name_col = cols[0], cols[1], cols[2], cols[3]
                try:
                    dec_val = int(dec_off)
                    rows.append((dec_val, type_col, name_col))
                except ValueError:
                    pass
        structs[name] = (total_size, rows)
    return structs

wm_t = parse_markdown_structs_with_tables("docs/Servicios/WorldManager.md")
char_t = parse_markdown_structs_with_tables("docs/Servicios/Character.md")
all_structs = {}
all_structs.update(char_t)
all_structs.update(wm_t)

with open("pkg/protocol/opcodes.go") as f:
    go_code = f.read()
opcodes = set(re.findall(r"(\bPROTO_NC_[A-Z0-9_]+)\s*(?:uint16)?\s*=", go_code))

def clean_field_name(name, seen):
    reserved = {"type", "range", "map", "select", "default", "var", "const", "func", "interface", "package", "import", "go", "defer", "chan", "case"}
    name = re.sub(r"[^A-Za-z0-9_]", "", name)
    if not name:
        name = "Field"
    pascal = name[0].upper() + name[1:]
    if pascal.lower() in reserved:
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
        if "char" in t:
            return f"[{count}]byte"
        elif "short" in t and field_size == count * 2:
            return f"[{count}]uint16"
        elif "int" in t and field_size == count * 4:
            return f"[{count}]uint32"
        elif field_size > 0:
            return f"[{field_size}]byte"
    if t in ("unsigned char", "uint8") and field_size == 1:
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
    if t in ("unsigned __int64", "uint64") and field_size == 8:
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

def classify_struct(name):
    n = name.upper()
    if "GUILD_ACADEMY" in n:
        return "proto_nc_guild_academy.go"
    elif "GUILD" in n:
        return "proto_nc_guild.go"
    elif "AVATAR" in n:
        return "proto_nc_avatar.go"
    elif "CHAR_OPTION" in n or "SHORTCUT" in n or "KEYMAP" in n:
        return "proto_nc_char_option.go"
    elif "RESTRICT" in n or "SPAM" in n or "CHAT_BLOCK" in n:
        return "proto_nc_chat_restrict.go"
    elif "RAID" in n:
        return "proto_nc_raid.go"
    elif "PARTY" in n:
        return "proto_nc_party.go"
    elif "FRIEND" in n:
        return "proto_nc_friend.go"
    elif "KQ" in n or "KINGDOM" in n:
        return "proto_nc_kq.go"
    elif "HOLY_PROMISE" in n or "MASPUP" in n:
        return "proto_nc_holy_promise.go"
    elif "PRISON" in n:
        return "proto_nc_prison.go"
    elif "ITEM" in n:
        return "proto_nc_itemdb.go"
    elif "QUEST" in n or "SKILL" in n:
        return "proto_nc_quest_skill.go"
    elif "MINIHOUSE" in n or "PET" in n or "MINIMON" in n or "FURNITURE" in n:
        return "proto_nc_minihouse_pet.go"
    elif "INSTANCE" in n or "RANK" in n or "DUNGEON" in n or "INDUN" in n:
        return "proto_nc_instance_rank.go"
    elif "ACT" in n or "CHAT" in n or "WHISPER" in n or "NOTICE" in n or "EMOT" in n:
        return "proto_nc_act_chat.go"
    elif "CHAR" in n:
        return "proto_nc_char.go"
    elif "USER" in n or "MISC" in n or "CONNECT" in n:
        return "proto_nc_user_world.go"
    else:
        return "proto_nc_char.go"

# Map generated structs to files
file_structs = {}
registry_entries = []

for name, (total_size, rows) in sorted(all_structs.items()):
    go_type_name = name + "_DATA" if name in opcodes else name
    target_file = classify_struct(name)
    file_structs.setdefault(target_file, []).append((name, go_type_name, total_size, rows))
    registry_entries.append((name, go_type_name, total_size))

# Existing files where we append
existing_files_to_append = {
    "proto_nc_char.go", "proto_nc_friend.go", "proto_nc_guild.go", "proto_nc_kq.go", "proto_nc_party.go"
}

for fname, s_list in sorted(file_structs.items()):
    fpath = os.path.join("pkg/protocol", fname)
    lines = []
    
    if fname in existing_files_to_append and os.path.exists(fpath):
        with open(fpath) as f:
            existing_content = f.read()
        marker = "// ============================================================================\n// Estructuras Oficiales C++ (#pragma pack(push, 1))"
        if marker in existing_content:
            existing_content = existing_content.split(marker)[0].strip()
        lines.append(existing_content.strip())
        lines.append("")
        lines.append(marker)
        lines.append("// ============================================================================")
        lines.append("")
    else:
        lines.append("package protocol")
        lines.append("")
        lines.append("// ============================================================================")
        lines.append("// Estructuras Oficiales C++ (#pragma pack(push, 1)) de WorldManager y Character")
        lines.append("// ============================================================================")
        lines.append("")

    for orig_name, go_name, total_size, rows in s_list:
        lines.append(f"// {orig_name} (size = {total_size} bytes)")
        if not rows:
            lines.append(f"type {go_name} struct {{}}")
            lines.append("")
            continue
        lines.append(f"type {go_name} struct {{")
        seen_fields = set()
        fixed_rows = [r for r in rows if "[0]" not in r[1]]
        for idx, (off, cpp_type, field_name) in enumerate(fixed_rows):
            if idx + 1 < len(fixed_rows):
                f_size = fixed_rows[idx+1][0] - off
            else:
                f_size = total_size - off
            go_f = clean_field_name(field_name, seen_fields)
            go_t = resolve_field_type(cpp_type, f_size)
            lines.append(f"\t{go_f:<32} {go_t:<20} // offset: {off}, size: {f_size} bytes")
        flex_rows = [r for r in rows if "[0]" in r[1]]
        for off, cpp_type, field_name in flex_rows:
            go_f = clean_field_name(field_name, seen_fields)
            lines.append(f"\t{go_f:<32} [0]byte              // offset: {off}, flexible array member")
        lines.append("}")
        lines.append("")

    with open(fpath, "w") as f:
        f.write("\n".join(lines))
    print(f"Written {fpath} with {len(s_list)} structs")

# Write world_character_registry.go
reg_lines = [
    "package protocol",
    "",
    "// WorldCharacterPDBStructs contiene una instancia de cada estructura oficial C++",
    "// no-OpTool de WorldManager.md y Character.md, indexada por su nombre PDB oficial.",
    "var WorldCharacterPDBStructs = map[string]interface{}{",
]
for orig_name, go_name, total_size in sorted(registry_entries):
    reg_lines.append(f'\t"{orig_name}": {go_name}{{}},')
    if go_name != orig_name:
        reg_lines.append(f'\t"{go_name}": {go_name}{{}},')
reg_lines.append("}")
reg_lines.append("")

with open("pkg/protocol/world_character_registry.go", "w") as f:
    f.write("\n".join(reg_lines))

print(f"Written pkg/protocol/world_character_registry.go with {len(registry_entries)} entries")
