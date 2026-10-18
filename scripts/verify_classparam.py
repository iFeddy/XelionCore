#!/usr/bin/env python3
"""
Verify 100% numerical and structural parity between Param*Server.txt and <Class>.json.
"""
import glob
import json
import os
from convert_classparam import COLUMNS

def verify():
    base_dir = "data/World/ClassParam"
    txt_files = sorted(glob.glob(os.path.join(base_dir, "Param*Server.txt")))
    assert len(txt_files) == 27, f"Expected 27 txt files, got {len(txt_files)}"

    total_records = 0
    total_fields = 0

    for fpath in txt_files:
        base = os.path.basename(fpath)
        class_name = base[5:-10]
        json_path = os.path.join(base_dir, f"{class_name}.json")
        assert os.path.exists(json_path), f"Missing JSON file: {json_path}"

        with open(json_path, "r", encoding="utf-8") as jf:
            json_data = json.load(jf)

        txt_records = []
        with open(fpath, "r", encoding="latin1") as fp:
            for line in fp:
                line = line.strip()
                if line.startswith("#record"):
                    parts = line.split()[1:]
                    txt_records.append(parts)

        assert len(txt_records) == len(json_data), f"{class_name}: count mismatch ({len(txt_records)} vs {len(json_data)})"
        assert len(json_data) == 150, f"{class_name}: expected 150 levels, got {len(json_data)}"

        for i, (txt_row, json_obj) in enumerate(zip(txt_records, json_data)):
            for (col_name, col_type), raw_val in zip(COLUMNS, txt_row):
                expected = col_type(raw_val)
                actual = json_obj[col_name]
                if expected != actual:
                    raise AssertionError(f"Mismatch in {class_name} level {json_obj['level']} col {col_name}: expected {expected}, got {actual}")
                total_fields += 1
            total_records += 1

    print(f"VERIFICATION SUCCESSFUL: 27 files, {total_records} level rows, {total_fields} fields verified with 0 discrepancies!")

if __name__ == "__main__":
    verify()
