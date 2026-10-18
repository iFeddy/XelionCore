#!/usr/bin/env python3
import os
import glob
import json
import re

def parse_txt_file(filepath):
    groups = {}
    
    with open(filepath, "r", encoding="cp949", errors="replace") as f:
        lines = f.readlines()
        
    cur_table = None
    group_cols = []
    regen_cols = []
    
    for line in lines:
        stripped = line.strip()
        if not stripped or stripped.startswith(";	"): # tab-comment
            continue
            
        # Detectar tabla
        m_tab = re.match(r"^#table\s+([a-zA-Z0-9_]+)", stripped, re.I)
        if m_tab:
            cur_table = m_tab.group(1).lower()
            continue
            
        # Detectar columnas
        if stripped.lower().startswith("#columnname"):
            parts = stripped.split()[1:]
            if cur_table == "mobregengroup":
                group_cols = parts
            elif cur_table == "mobregen":
                regen_cols = parts
            continue
            
        # Detectar registros (activos o comentados)
        is_active = stripped.lower().startswith("#record")
        is_commented = stripped.lower().startswith(";#record")
        
        if not (is_active or is_commented):
            continue
            
        # Limpiar inline comments
        clean_line = stripped
        if ";" in clean_line and not is_commented:
            clean_line = clean_line.split(";", 1)[0].strip()
        elif is_commented and ";" in clean_line[1:]:
            # Comentario dentro de un registro ya comentado
            sub = clean_line[1:]
            if ";" in sub:
                clean_line = ";" + sub.split(";", 1)[0].strip()
                
        # Extraer tokens separados por \t
        raw_parts = [p.strip() for p in clean_line.split("\t") if p.strip()]
        if not raw_parts:
            continue
        tokens = raw_parts[1:] # quitar #record o ;#record
        if not tokens:
            continue
            
        if cur_table == "mobregengroup":
            # Columnas: GroupIndex, IsFamily, CenterX, CenterY, Width, Height, RangeDegree
            gid = tokens[0]
            is_family = tokens[1].upper() == "Y" if len(tokens) > 1 else False
            cx = int(tokens[2]) if len(tokens) > 2 and tokens[2].lstrip("-").isdigit() else 0
            cy = int(tokens[3]) if len(tokens) > 3 and tokens[3].lstrip("-").isdigit() else 0
            w = int(tokens[4]) if len(tokens) > 4 and tokens[4].lstrip("-").isdigit() else 0
            h = int(tokens[5]) if len(tokens) > 5 and tokens[5].lstrip("-").isdigit() else 0
            deg = int(tokens[6]) if len(tokens) > 6 and tokens[6].lstrip("-").isdigit() else 0
            
            groups[gid] = {
                "enabled": is_active,
                "is_family": is_family,
                "center_x": cx,
                "center_y": cy,
                "width": w,
                "height": h,
                "range_degree": deg,
                "mobs": []
            }
            
        elif cur_table == "mobregen":
            gid = tokens[0]
            if gid not in groups:
                # Caso borde: si el grupo no fue definido antes, crearlo
                groups[gid] = {
                    "enabled": is_active,
                    "is_family": False,
                    "center_x": 0,
                    "center_y": 0,
                    "width": 0,
                    "height": 0,
                    "range_degree": 0,
                    "mobs": []
                }
                
            has_kill_num = "KillNum" in regen_cols
            
            mob_index = tokens[1] if len(tokens) > 1 else ""
            mob_num = int(tokens[2]) if len(tokens) > 2 and tokens[2].isdigit() else 0
            
            offset = 3
            kill_num = 0
            if has_kill_num and len(tokens) > 3:
                kill_num = int(tokens[3]) if tokens[3].isdigit() else 0
                offset = 4
                
            reg_std = int(tokens[offset]) if len(tokens) > offset and tokens[offset].isdigit() else 0
            reg_min = int(tokens[offset+1]) if len(tokens) > offset+1 and tokens[offset+1].isdigit() else 0
            reg_max = int(tokens[offset+2]) if len(tokens) > offset+2 and tokens[offset+2].isdigit() else 0
            
            # Deltas y segundos de aceleración
            deltas = []
            delta_secs = []
            
            # Formato estándar de curva: Delta0, Sec0, Delta1, Sec1, Delta2, Sec2, Delta3, Sec3, Delta4
            curv_idx = offset + 3
            if len(tokens) >= curv_idx + 9:
                for k in range(5):
                    d_pos = curv_idx + (k * 2)
                    if d_pos < len(tokens):
                        val_str = tokens[d_pos].lstrip("-")
                        d_val = int(tokens[d_pos]) if val_str.isdigit() else 0
                        deltas.append(d_val)
                    if k < 4:
                        s_pos = curv_idx + (k * 2) + 1
                        if s_pos < len(tokens):
                            s_val = int(tokens[s_pos]) if tokens[s_pos].isdigit() else 0
                            delta_secs.append(s_val)
                            
            mob_entry = {
                "mob": mob_index,
                "count": mob_num,
                "kill_num": kill_num,
                "reg_standard": reg_std,
                "reg_min": reg_min,
                "reg_max": reg_max,
                "deltas": deltas,
                "delta_secs": delta_secs,
                "enabled": is_active
            }
            
            groups[gid]["mobs"].append(mob_entry)
            
    return groups

def convert_all():
    base_dir = "data/MobRegen"
    txt_files = sorted(glob.glob(f"{base_dir}/**/*.txt", recursive=True))
    if not txt_files:
        base_dir = "../AstraServer/Server/9Data/Shine/MobRegen"
        txt_files = sorted(glob.glob(f"{base_dir}/**/*.txt", recursive=True))

    total_files = len(txt_files)
    total_groups = 0
    total_mobs = 0

    for f in txt_files:
        rel = os.path.relpath(f, base_dir)
        json_path = os.path.join("data/MobRegen", os.path.splitext(rel)[0] + ".json")
        os.makedirs(os.path.dirname(json_path), exist_ok=True)
        data = parse_txt_file(f)

        total_groups += len(data)
        for g in data.values():
            total_mobs += len(g["mobs"])

        with open(json_path, "w", encoding="utf-8") as out:
            json.dump(data, out, indent=2, ensure_ascii=False)

    print(f"Convertidos exitosamente {total_files} archivos .txt a .json")
    print(f"Total Grupos procesados: {total_groups}")
    print(f"Total Mobs procesados: {total_mobs}")

if __name__ == "__main__":
    convert_all()
