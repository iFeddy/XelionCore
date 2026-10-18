#!/usr/bin/env python3
import os
import glob
import json
from convert_mobregen_to_json import parse_txt_file

def verify_all():
    base_dir = "data/MobRegen"
    txt_files = sorted(glob.glob(f"{base_dir}/**/*.txt", recursive=True))
    
    total_files = len(txt_files)
    total_groups_checked = 0
    total_mobs_checked = 0
    mismatches = []
    
    for txt_path in txt_files:
        json_path = os.path.splitext(txt_path)[0] + ".json"
        if not os.path.exists(json_path):
            mismatches.append(f"Falta archivo JSON: {json_path}")
            continue
            
        txt_data = parse_txt_file(txt_path)
        with open(json_path, "r", encoding="utf-8") as jf:
            json_data = json.load(jf)
            
        # 1. Comparar cantidad de grupos
        if len(txt_data) != len(json_data):
            mismatches.append(f"[{txt_path}] Cantidad de grupos diferente: TXT={len(txt_data)} vs JSON={len(json_data)}")
            
        # 2. Comparar cada grupo y cada mob
        for gid, t_group in txt_data.items():
            if gid not in json_data:
                mismatches.append(f"[{txt_path}] Grupo {gid} no encontrado en JSON")
                continue
                
            j_group = json_data[gid]
            total_groups_checked += 1
            
            # Chequear atributos del grupo
            for attr in ["enabled", "is_family", "center_x", "center_y", "width", "height", "range_degree"]:
                if t_group[attr] != j_group[attr]:
                    mismatches.append(f"[{txt_path}][{gid}] Mismatch en {attr}: TXT={t_group[attr]} vs JSON={j_group[attr]}")
                    
            # Chequear mobs
            t_mobs = t_group["mobs"]
            j_mobs = j_group["mobs"]
            if len(t_mobs) != len(j_mobs):
                mismatches.append(f"[{txt_path}][{gid}] Cantidad de mobs diferente: TXT={len(t_mobs)} vs JSON={len(j_mobs)}")
                continue
                
            for idx in range(len(t_mobs)):
                tm = t_mobs[idx]
                jm = j_mobs[idx]
                total_mobs_checked += 1
                
                for field in ["mob", "count", "kill_num", "reg_standard", "reg_min", "reg_max", "deltas", "delta_secs", "enabled"]:
                    if tm[field] != jm[field]:
                        mismatches.append(f"[{txt_path}][{gid}][mob #{idx}] Mismatch en {field}: TXT={tm[field]} vs JSON={jm[field]}")
                        
    print("=" * 70)
    print("REPORTE DE VERIFICACIÓN DE PARIDAD MOBREGEN (.txt vs .json)")
    print("=" * 70)
    print(f"Archivos verificados: {total_files}/126")
    print(f"Grupos verificados campo por campo: {total_groups_checked}")
    print(f"Mobs verificados campo por campo: {total_mobs_checked}")
    print(f"Total discrepancias encontradas: {len(mismatches)}")
    
    if mismatches:
        print("\nPrimeras discrepancias encontradas:")
        for m in mismatches[:10]:
            print("  -", m)
        return False
    else:
        print("\n>>> RESULTADO: 100% PARIDAD CONFIRMADA. CERO DISCREPANCIAS. <<<")
        return True

if __name__ == "__main__":
    verify_all()
