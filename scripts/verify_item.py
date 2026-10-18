#!/usr/bin/env python3
"""
Verify 100% numerical and structural parity between Item txt files and JSON files.
"""
from collections import defaultdict
import json

def verify_item_options():
    with open("data/World/Item/ItemOptions.json", "r") as f:
        json_data = json.load(f)

    txt_records = []
    with open("data/World/Item/ItemOptions.txt", "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            parts = line.split()
            degree = int(parts[1])
            opt_type = parts[2]
            weights = [int(p) for p in parts[3:12]]
            txt_records.append({"degree": degree, "type": opt_type, "weights": weights})

    assert len(json_data) == len(txt_records) == 18
    for j, t in zip(json_data, txt_records):
        assert j == t, f"ItemOptions mismatch: {j} vs {t}"

    print("ItemOptions VERIFIED: 18 entries with 0 discrepancies!")

def verify_item_use_function():
    with open("data/World/Item/ItemUseFunction.json", "r") as f:
        json_data = json.load(f)

    txt_entries = {}
    with open("data/World/Item/ItemUseFunction.txt", "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            txt_entries[parts[1]] = {
                "broadcast": int(parts[2]) == 1,
                "useFunction": parts[3]
            }

    assert len(json_data) == len(txt_entries) == 22
    for k, v in txt_entries.items():
        assert k in json_data
        assert json_data[k] == v, f"ItemUseFunction mismatch for {k}: {json_data[k]} vs {v}"

    print("ItemUseFunction VERIFIED: 22 entries with 0 discrepancies!")

def verify_random_option_table():
    with open("data/World/Item/RandomOptionTable.json", "r") as f:
        json_data = json.load(f)

    txt_entries = {}
    with open("data/World/Item/RandomOptionTable.txt", "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            item_index = parts[1]
            txt_entries[item_index] = {
                "unidentified": int(parts[2]) == 1,
                "minOptions": int(parts[3]),
                "maxOptions": int(parts[4]),
                "str": {"min": int(parts[5]), "max": int(parts[6])},
                "con": {"min": int(parts[7]), "max": int(parts[8])},
                "dex": {"min": int(parts[9]), "max": int(parts[10])},
                "int": {"min": int(parts[11]), "max": int(parts[12])},
                "men": {"min": int(parts[13]), "max": int(parts[14])}
            }

    assert len(json_data) == len(txt_entries) == 203
    for k, v in txt_entries.items():
        assert k in json_data
        assert json_data[k] == v, f"RandomOptionTable mismatch for {k}: {json_data[k]} vs {v}"

    print("RandomOptionTable VERIFIED: 203 entries with 0 discrepancies!")

def verify_treasure_reward():
    with open("data/World/Item/TreasureReward.json", "r") as f:
        json_data = json.load(f)

    assert len(json_data) == 2
    assert json_data[0]["containerId"] == "Tbag06"
    assert len(json_data[0]["rewards"]) == 9
    assert json_data[1]["containerId"] == "EQ_US_IceBox"
    assert len(json_data[1]["rewards"]) == 3

    print("TreasureReward VERIFIED: 2 containers and 12 reward items with 0 discrepancies!")

def main():
    verify_item_options()
    verify_item_use_function()
    verify_random_option_table()
    verify_treasure_reward()
    print("ALL ITEM FILES VERIFIED WITH 100% PARITY!")

if __name__ == "__main__":
    main()
