import os
import struct
import json

base_dir = "./data/AbState" if os.path.exists("./data/AbState") else "../AstraServer/Server/9Data/Shine/AbState"
info_path = os.path.join(base_dir, "AbStateInfo.dat")
output_path = "./data/abstate.json"

STATE_METADATA = {
    1: {"name": "StaHardSkin", "description": "Piel de Piedra (Defensa física de Fighter)", "fields": ["level", "duration_ms", "min_def", "max_def"]},
    11: {"name": "StaCleHold", "description": "Inmovilización / Hold de Cleric", "fields": ["level", "duration_ms"]},
    13: {"name": "StaShield", "description": "Escudo Sagrado protector de Cleric", "fields": ["level", "duration_ms", "min_shield", "max_shield"]},
    33: {"name": "StaArcDmgDown", "description": "Reducción de daño de arquero", "fields": ["level", "duration_ms", "min_dmg_down", "max_dmg_down", "rate"]},
    35: {"name": "StaDodge", "description": "Incremento de evasión de Archer/Joker", "fields": ["level", "duration_ms", "min_dodge", "max_dodge", "evade_rate"]},
    41: {"name": "StaVine", "description": "Enredaderas de la naturaleza (Retención)", "fields": ["level", "duration_ms"]},
    45: {"name": "StaMagSpdDown", "description": "Ralentización de velocidad mágica", "fields": ["level", "duration_ms", "speed_down_rate"]},
    46: {"name": "StaLightDmg", "description": "Daño periódico de relámpago (DoT)", "fields": ["level", "duration_ms", "dot_damage", "dot_interval_ms"]},
    48: {"name": "StaRingShield", "description": "Escudo de anillo protector", "fields": ["level", "duration_ms", "min_shield", "max_shield"]},
    57: {"name": "StaCleDefDown", "description": "Debuff de bajada de defensa de Cleric", "fields": ["level", "duration_ms", "min_def_down", "max_def_down", "rate"]}
}

with open(info_path, "rb") as f:
    count, = struct.unpack("<I", f.read(4))
    entries = []
    for _ in range(count):
        raw = f.read(104)
        state_id, = struct.unpack("<H", raw[:2])
        name = raw[2:34].split(b"\x00")[0].decode("ascii", errors="ignore")
        entries.append((state_id, name, raw))

abstate_db = {}

for state_id, name, raw_entry in entries:
    dat_file = os.path.join(base_dir, f"{name}.dat")
    if not os.path.exists(dat_file):
        print(f"Warning: {name}.dat not found!")
        continue

    meta = STATE_METADATA.get(state_id, {"name": name, "description": name, "fields": []})
    fields = meta["fields"]

    with open(dat_file, "rb") as df:
        max_lv, = struct.unpack("<I", df.read(4))
        rec_size = (os.path.getsize(dat_file) - 4) // max_lv
        num_ints = rec_size // 4
        
        levels = []
        for lv in range(max_lv):
            raw = df.read(rec_size)
            ints = list(struct.unpack(f"<{num_ints}I", raw))
            
            # Map to semantic dict
            entry_dict = {}
            for idx, val in enumerate(ints):
                if idx < len(fields):
                    entry_dict[fields[idx]] = val
                else:
                    entry_dict[f"param_{idx}"] = val
            entry_dict["raw"] = ints
            levels.append(entry_dict)

    abstate_db[name] = {
        "id": state_id,
        "name": name,
        "description": meta["description"],
        "max_levels": max_lv,
        "record_size_bytes": rec_size,
        "meta_hex": raw_entry.hex(),
        "levels": levels
    }

with open(output_path, "w", encoding="utf-8") as out:
    json.dump(abstate_db, out, indent=2, ensure_ascii=False)

print(f"Generated {output_path} successfully ({len(abstate_db)} states, {os.path.getsize(output_path)} bytes)")
