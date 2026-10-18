import os
import struct
import json
import argparse
import hashlib

def json_to_dat(json_path="./data/abstate.json", output_dir="./data/AbState"):
    with open(json_path, "r", encoding="utf-8") as f:
        data = json.load(f)

    os.makedirs(output_dir, exist_ok=True)
    states = data.get("states", data)

    # 1. Generar AbStateInfo.dat
    info_path = os.path.join(output_dir, "AbStateInfo.dat")
    with open(info_path, "wb") as f_info:
        # Header: total states (uint32)
        f_info.write(struct.pack("<I", len(states)))
        for name, sinfo in sorted(states.items(), key=lambda x: x[1]["id"]):
            if "meta_hex" in sinfo:
                f_info.write(bytes.fromhex(sinfo["meta_hex"]))
            else:
                state_id = sinfo["id"]
                record = bytearray(104)
                struct.pack_into("<H", record, 0, state_id)
                name_bytes = name.encode("ascii")[:31]
                record[2:2+len(name_bytes)] = name_bytes
                f_info.write(record)

    # 2. Generar cada Sta*.dat
    for name, sinfo in states.items():
        dat_path = os.path.join(output_dir, f"{name}.dat")
        levels = sinfo["levels"]
        max_levels = sinfo["max_levels"]

        with open(dat_path, "wb") as f_dat:
            # Header: max_levels (uint32)
            f_dat.write(struct.pack("<I", max_levels))
            for lv_data in levels:
                raw_ints = lv_data.get("raw")
                if not raw_ints:
                    # fallback to reconstructing from semantic fields
                    print(f"Warning: no raw array found for {name}")
                    continue
                f_dat.write(struct.pack(f"<{len(raw_ints)}I", *raw_ints))

    print(f"Exported {len(states)} .dat files to {output_dir}")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Convierte data/abstate.json de vuelta a archivos binarios .dat")
    parser.add_argument("--json", default="./data/abstate.json", help="Ruta al archivo abstate.json")
    parser.add_argument("--out", default="./data/AbState", help="Directorio de salida para los archivos .dat")
    args = parser.parse_args()

    json_to_dat(args.json, args.out)
