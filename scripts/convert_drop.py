#!/usr/bin/env python3
"""
Convert legacy World/Drop txt files into clean, modular JSON files with 100% data fidelity.
"""
from collections import defaultdict
import json
import os

def convert_item_drop_table():
    fpath = "data/World/Drop/ItemDropTable.txt"
    mobs_map = {}

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.startswith("#RECORD"):
                continue

            parts = line.split()
            map_area = parts[1]
            mob_id = parts[2]
            min_lvl = int(parts[3])
            max_lvl = int(parts[4])
            ab_cnt = int(parts[5])
            min_cen = int(parts[6])
            max_cen = int(parts[7])
            cen_rate = int(parts[8])

            trade_boxes = []
            for b_idx in [(9, 10), (11, 12), (13, 14)]:
                box_name = parts[b_idx[0]]
                box_rate = int(parts[b_idx[1]])
                if box_name != "-" and box_rate > 0:
                    trade_boxes.append({"box": box_name, "rate": box_rate})

            drops = []
            for slot in range(45):
                b = 15 + slot * 6
                item = parts[b]
                rate = int(parts[b + 1])
                up_min = int(parts[b + 2])
                up_max = int(parts[b + 3])
                rule = parts[b + 4]
                num = int(parts[b + 5])
                if item != "-" and rate > 0:
                    drops.append({
                        "item": item,
                        "rate": rate,
                        "minUpgrade": up_min,
                        "maxUpgrade": up_max,
                        "rule": rule,
                        "num": num
                    })

            exc = []
            for e_idx in range(285, 290):
                if e_idx < len(parts) and parts[e_idx] != "-":
                    exc.append(parts[e_idx])

            mobs_map[mob_id] = {
                "mapArea": map_area,
                "minLevel": min_lvl,
                "maxLevel": max_lvl,
                "abStateCnt": ab_cnt,
                "money": {"min": min_cen, "max": max_cen, "rate": cen_rate},
                "tradeBoxes": trade_boxes,
                "drops": drops,
                "excludedItems": exc
            }

    out_path = "data/World/Drop/ItemDropTable.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(mobs_map, out, indent=2)

    print(f"ItemDropTable converted -> {out_path} ({len(mobs_map)} mobs)")

def convert_item_drop_group():
    fpath = "data/World/Drop/ItemDropGroup.txt"
    groups = defaultdict(list)

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line.startswith("#RECORD"):
                continue

            parts = line.split()
            grp = parts[1]
            item = parts[2]
            min_q = int(parts[3])
            max_q = int(parts[4])
            rates = [int(p) for p in parts[5:21]]

            groups[grp].append({
                "itemId": item,
                "minQuantity": min_q,
                "maxQuantity": max_q,
                "upgradeRates": rates
            })

    out_path = "data/World/Drop/ItemDropGroup.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(groups, out, indent=2)

    print(f"ItemDropGroup converted -> {out_path} ({len(groups)} groups)")

def main():
    convert_item_drop_table()
    convert_item_drop_group()

if __name__ == "__main__":
    main()
