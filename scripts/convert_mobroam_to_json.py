#!/usr/bin/env python3
import os
import glob
import json

def parse_roam_file(filepath):
    route_name = os.path.splitext(os.path.basename(filepath))[0]
    points = []
    loop_mode = "none"

    with open(filepath, "r", encoding="cp949", errors="replace") as f:
        for line in f:
            stripped = line.strip()
            if not stripped or stripped.startswith(";"):
                continue

            if stripped.lower().startswith("#record"):
                parts = stripped.split()
                # parts: [#record, ID, X, Y, EventIndex]
                if len(parts) >= 5:
                    step_id = int(parts[1]) if parts[1].isdigit() else 0
                    x = int(parts[2]) if parts[2].lstrip("-").isdigit() else 0
                    y = int(parts[3]) if parts[3].lstrip("-").isdigit() else 0
                    event = parts[4]

                    if event in ["rotate", "return"]:
                        loop_mode = event
                    else:
                        points.append({
                            "id": step_id,
                            "x": x,
                            "y": y
                        })

    return route_name, {
        "loop_mode": loop_mode,
        "total_points": len(points),
        "points": points
    }

def convert_all():
    base_dir = "data/MobRoam"
    txt_files = sorted(glob.glob(f"{base_dir}/*.txt"))
    if not txt_files:
        base_dir = "../AstraServer/Server/9Data/Shine/MobRoam"
        txt_files = sorted(glob.glob(f"{base_dir}/*.txt"))

    output_path = "data/MobRoam/mob_roam.json"
    os.makedirs(os.path.dirname(output_path), exist_ok=True)

    data = {}
    total_waypoints = 0

    for f in txt_files:
        name, route = parse_roam_file(f)
        data[name] = route
        total_waypoints += route["total_points"]

    with open(output_path, "w", encoding="utf-8") as out:
        json.dump(data, out, indent=2, ensure_ascii=False)

    size = os.path.getsize(output_path)
    print(f"Generated {output_path} successfully!")
    print(f"Total Routes: {len(data)}/28")
    print(f"Total Waypoints: {total_waypoints}")
    print(f"File Size: {size} bytes ({size / 1024:.2f} KB)")

if __name__ == "__main__":
    convert_all()
