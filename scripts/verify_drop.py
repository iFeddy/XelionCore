#!/usr/bin/env python3
"""
Verify 100% numerical and structural parity between Drop txt files and JSON files.
"""
from collections import defaultdict
import json

def verify_item_drop_table():
    with open("data/World/Drop/ItemDropTable.json", "r", encoding="utf-8") as f:
        json_mobs = json.load(f)

    txt_mobs = {}
    with open("data/World/Drop/ItemDropTable.txt", "r", encoding="latin1") as f:
        for line in f:
            if not line.startswith("#RECORD"):
                continue
            parts = line.split()
            mob_id = parts[2]

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

            txt_mobs[mob_id] = {
                "mapArea": parts[1],
                "minLevel": int(parts[3]),
                "maxLevel": int(parts[4]),
                "abStateCnt": int(parts[5]),
                "minCen": int(parts[6]),
                "maxCen": int(parts[7]),
                "cenRate": int(parts[8]),
                "drops": drops
            }

    assert len(json_mobs) == len(txt_mobs) == 1485, f"Count mismatch: json={len(json_mobs)}, txt={len(txt_mobs)}"

    total_drops = 0
    for mob_id, txt_data in txt_mobs.items():
        assert mob_id in json_mobs, f"Missing mob in JSON: {mob_id}"
        j_data = json_mobs[mob_id]

        assert j_data["mapArea"] == txt_data["mapArea"]
        assert j_data["minLevel"] == txt_data["minLevel"]
        assert j_data["maxLevel"] == txt_data["maxLevel"]
        assert j_data["abStateCnt"] == txt_data["abStateCnt"]
        assert j_data["money"]["min"] == txt_data["minCen"]
        assert j_data["money"]["max"] == txt_data["maxCen"]
        assert j_data["money"]["rate"] == txt_data["cenRate"]

        assert len(j_data["drops"]) == len(txt_data["drops"]), f"{mob_id}: drops count mismatch"
        for j_drop, t_drop in zip(j_data["drops"], txt_data["drops"]):
            assert j_drop == t_drop, f"{mob_id} drop mismatch: {j_drop} vs {t_drop}"
            total_drops += 1

    print(f"ItemDropTable VERIFIED: 1485 mobs and {total_drops} active drops verified with 0 discrepancies!")

def verify_item_drop_group():
    with open("data/World/Drop/ItemDropGroup.json", "r", encoding="utf-8") as f:
        json_groups = json.load(f)

    txt_groups = defaultdict(list)
    with open("data/World/Drop/ItemDropGroup.txt", "r", encoding="latin1") as f:
        for line in f:
            if not line.startswith("#RECORD"):
                continue
            parts = line.split()
            grp = parts[1]
            txt_groups[grp].append({
                "itemId": parts[2],
                "minQuantity": int(parts[3]),
                "maxQuantity": int(parts[4]),
                "upgradeRates": [int(p) for p in parts[5:21]]
            })

    assert len(json_groups) == len(txt_groups) == 959, f"Group count mismatch: {len(json_groups)} vs {len(txt_groups)}"

    total_entries = 0
    for grp, entries in txt_groups.items():
        assert grp in json_groups, f"Missing group in JSON: {grp}"
        j_entries = json_groups[grp]
        assert len(j_entries) == len(entries), f"{grp}: entry count mismatch"
        for j_e, t_e in zip(j_entries, entries):
            assert j_e == t_e, f"{grp} entry mismatch: {j_e} vs {t_e}"
            total_entries += 1

    print(f"ItemDropGroup VERIFIED: 959 groups and {total_entries} item mappings verified with 0 discrepancies!")

def main():
    verify_item_drop_table()
    verify_item_drop_group()
    print("ALL DROP FILES VERIFIED WITH 100% PARITY!")

if __name__ == "__main__":
    main()
