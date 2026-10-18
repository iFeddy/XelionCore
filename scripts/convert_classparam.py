#!/usr/bin/env python3
"""
Convert legacy Param<Class>Server.txt into modular <Class>.json with clean naming and 100% data fidelity.
"""
import glob
import json
import os

COLUMNS = [
    ("level", int),
    ("strength", int),
    ("constitution", int),
    ("intelligence", int),
    ("wizdom", int),
    ("dexterity", int),
    ("mentalPower", int),
    ("soulHp", int),
    ("maxSoulHp", int),
    ("priceHpStone", int),
    ("soulSp", int),
    ("maxSoulSp", int),
    ("priceSpStone", int),
    ("atkPerAp", int),
    ("dmgPerAp", int),
    ("maxPwrStone", int),
    ("numPwrStone", int),
    ("pricePwrStone", int),
    ("pwrStoneWc", int),
    ("pwrStoneMa", int),
    ("maxGrdStone", int),
    ("numGrdStone", int),
    ("priceGrdStone", int),
    ("grdStoneAc", int),
    ("grdStoneMr", int),
    ("painRes", int),
    ("restraintRes", int),
    ("curseRes", int),
    ("shockRes", int),
    ("maxHp", int),
    ("maxSp", int),
    ("charTitlePt", int),
    ("skillPwrPt", int),
    ("jobChangeDmgUp", int),
]

def convert():
    base_dir = "data/World/ClassParam"
    txt_files = sorted(glob.glob(os.path.join(base_dir, "Param*Server.txt")))
    print(f"Found {len(txt_files)} class param txt files.")

    for fpath in txt_files:
        base = os.path.basename(fpath)
        class_name = base[5:-10] # Param<Class>Server.txt -> <Class>
        out_path = os.path.join(base_dir, f"{class_name}.json")

        records = []
        with open(fpath, "r", encoding="latin1") as fp:
            for line in fp:
                line = line.strip()
                if not line.startswith("#record"):
                    continue
                parts = line.split()[1:]
                if len(parts) != len(COLUMNS):
                    raise ValueError(f"Column count mismatch in {base}: expected {len(COLUMNS)}, got {len(parts)}")
                
                entry = {}
                for (col_name, col_type), val in zip(COLUMNS, parts):
                    entry[col_name] = col_type(val)
                records.append(entry)

        with open(out_path, "w", encoding="utf-8") as out:
            json.dump(records, out, indent=2)
        print(f"Converted {base} -> {class_name}.json ({len(records)} levels)")

if __name__ == "__main__":
    convert()
