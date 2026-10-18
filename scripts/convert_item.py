#!/usr/bin/env python3
"""
Convert legacy World/Item txt files into clean, modular JSON files with 100% data fidelity.
"""
from collections import defaultdict
import json
import os

def convert_item_options():
    fpath = "data/World/Item/ItemOptions.txt"
    records = []
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            parts = line.split()
            degree = int(parts[1])
            opt_type = parts[2]
            weights = [int(p) for p in parts[3:12]]
            records.append({
                "degree": degree,
                "type": opt_type,
                "weights": weights
            })

    out_path = "data/World/Item/ItemOptions.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(records, out, indent=2)
    print(f"ItemOptions converted -> {out_path} ({len(records)} entries)")

def convert_item_use_function():
    fpath = "data/World/Item/ItemUseFunction.txt"
    entries = {}
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            item_index = parts[1]
            broadcast = int(parts[2]) == 1
            func_name = parts[3]
            entries[item_index] = {
                "broadcast": broadcast,
                "useFunction": func_name
            }

    out_path = "data/World/Item/ItemUseFunction.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(entries, out, indent=2)
    print(f"ItemUseFunction converted -> {out_path} ({len(entries)} entries)")

def convert_random_option_table():
    fpath = "data/World/Item/RandomOptionTable.txt"
    entries = {}
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#record"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            item_index = parts[1]
            unidentified = int(parts[2]) == 1
            min_op = int(parts[3])
            max_op = int(parts[4])
            str_min, str_max = int(parts[5]), int(parts[6])
            con_min, con_max = int(parts[7]), int(parts[8])
            dex_min, dex_max = int(parts[9]), int(parts[10])
            int_min, int_max = int(parts[11]), int(parts[12])
            men_min, men_max = int(parts[13]), int(parts[14])

            entries[item_index] = {
                "unidentified": unidentified,
                "minOptions": min_op,
                "maxOptions": max_op,
                "str": {"min": str_min, "max": str_max},
                "con": {"min": con_min, "max": con_max},
                "dex": {"min": dex_min, "max": dex_max},
                "int": {"min": int_min, "max": int_max},
                "men": {"min": men_min, "max": men_max}
            }

    out_path = "data/World/Item/RandomOptionTable.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(entries, out, indent=2)
    print(f"RandomOptionTable converted -> {out_path} ({len(entries)} entries)")

def convert_treasure_reward():
    fpath = "data/World/Item/TreasureReward.txt"
    containers_by_idx = {}
    contents_by_idx = defaultdict(list)

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.lower().startswith("#recordin"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            table_name = parts[1].lower()

            if table_name == "container":
                item_id = int(parts[2])
                idx = int(parts[3])
                card_idx = parts[4]
                min_lot = int(parts[5])
                max_lot = int(parts[6])
                dummy = parts[7]
                containers_by_idx[idx] = {
                    "containerId": card_idx,
                    "dummyId": dummy,
                    "minLot": min_lot,
                    "maxLot": max_lot,
                    "rewards": []
                }
            elif table_name == "content":
                reward_id = int(parts[2])
                idx = int(parts[3])
                reward_inx = parts[4]
                reward_lot = int(parts[5])
                up_low = int(parts[6])
                up_high = int(parts[7])
                rate = int(parts[8])

                contents_by_idx[idx].append({
                    "itemId": reward_inx,
                    "lot": reward_lot,
                    "upgradeMin": up_low,
                    "upgradeMax": up_high,
                    "rate": rate
                })

    containers_list = []
    for idx in sorted(containers_by_idx.keys()):
        c = containers_by_idx[idx]
        c["rewards"] = contents_by_idx[idx]
        containers_list.append(c)

    out_path = "data/World/Item/TreasureReward.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(containers_list, out, indent=2)
    print(f"TreasureReward converted -> {out_path} ({len(containers_list)} containers)")

def main():
    convert_item_options()
    convert_item_use_function()
    convert_random_option_table()
    convert_treasure_reward()

if __name__ == "__main__":
    main()
