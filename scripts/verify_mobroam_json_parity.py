#!/usr/bin/env python3
import os
import glob
import json
from convert_mobroam_to_json import parse_roam_file

def verify():
    base_dir = "data/MobRoam"
    txt_files = sorted(glob.glob(f"{base_dir}/*.txt"))
    json_path = "data/MobRoam/mob_roam.json"

    with open(json_path, "r", encoding="utf-8") as f:
        json_data = json.load(f)

    mismatches = []
    total_waypoints_checked = 0

    if len(txt_files) != len(json_data):
        mismatches.append(f"Cantidad de rutas diferente: TXT={len(txt_files)} vs JSON={len(json_data)}")

    for txt_path in txt_files:
        name, expected_route = parse_roam_file(txt_path)
        if name not in json_data:
            mismatches.append(f"Ruta {name} no encontrada en JSON")
            continue

        actual_route = json_data[name]

        if expected_route["loop_mode"] != actual_route["loop_mode"]:
            mismatches.append(f"[{name}] Loop mode mismatch: TXT={expected_route['loop_mode']} vs JSON={actual_route['loop_mode']}")

        if len(expected_route["points"]) != len(actual_route["points"]):
            mismatches.append(f"[{name}] Puntos mismatch: TXT={len(expected_route['points'])} vs JSON={len(actual_route['points'])}")
            continue

        for idx, (p_exp, p_act) in enumerate(zip(expected_route["points"], actual_route["points"])):
            total_waypoints_checked += 1
            if p_exp["id"] != p_act["id"] or p_exp["x"] != p_act["x"] or p_exp["y"] != p_act["y"]:
                mismatches.append(f"[{name}][pt #{idx}] Mismatch: TXT={p_exp} vs JSON={p_act}")

    print("=" * 60)
    print("REPORTE DE VERIFICACIÓN DE MOBROAM (.txt vs .json)")
    print("=" * 60)
    print(f"Rutas verificadas: {len(txt_files)}/28")
    print(f"Waypoints verificados campo por campo: {total_waypoints_checked}/470")
    print(f"Total discrepancias: {len(mismatches)}")

    if mismatches:
        print("Errores encontrados:")
        for m in mismatches[:10]:
            print("  -", m)
        return False
    else:
        print(">>> RESULTADO: 100% PARIDAD CONFIRMADA. CERO DISCREPANCIAS. <<<")
        return True

if __name__ == "__main__":
    verify()
