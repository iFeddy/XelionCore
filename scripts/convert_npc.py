#!/usr/bin/env python3
"""
Convert data/World/NPC/ legacy files (NPC.txt and NPCAction.txt) into modular JSON files:
- data/World/NPC/Maps/<Map>.json (95 map files with ShineNPC records)
- data/World/NPC/LinkTable.json (portal / warp transitions)
- data/World/NPC/NPCAction.json (conditions and actions)

Preserves all commented entries with "enabled": false.
"""

import json
import os
from collections import defaultdict

BASE_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
NPC_DIR = os.path.join(BASE_DIR, "data", "World", "NPC")
MAPS_DIR = os.path.join(NPC_DIR, "Maps")

def parse_npc_txt():
    npc_txt_path = os.path.join(NPC_DIR, "NPC.txt")
    with open(npc_txt_path, "r", encoding="latin1") as f:
        lines = f.readlines()

    shine_by_map = defaultdict(list)
    link_records = []

    for line_num, line in enumerate(lines, 1):
        raw = line.strip()
        if not raw:
            continue

        is_comment = line.lstrip().startswith(";")
        uncommented = line.lstrip(";\t ")
        if not (uncommented.startswith("#recordin") or uncommented.startswith("#Record") or uncommented.startswith("#record")):
            continue

        tokens_all = uncommented.split()
        table_name = tokens_all[1]

        first_tok = tokens_all[0]
        after_first = uncommented[len(first_tok):]
        is_tab_delimited = after_first.startswith("\t")

        if table_name == "ShineNPC":
            if is_tab_delimited:
                tabs = [t.strip() for t in line.split("\t")]
                rec_idx = next(i for i, t in enumerate(tabs) if "#recordin" in t or "#Record" in t)
                fields = tabs[rec_idx+2 : rec_idx+10]
                while len(fields) < 8:
                    fields.append("")
            else:
                fields = tokens_all[2:10]
                while len(fields) < 8:
                    fields.append("")

            map_name = fields[1]
            rec = {
                "mobName": fields[0],
                "map": map_name,
                "x": int(fields[2]) if fields[2] and fields[2] != "-" else 0,
                "y": int(fields[3]) if fields[3] and fields[3] != "-" else 0,
                "direct": int(fields[4]) if fields[4] and fields[4] != "-" else 0,
                "menu": int(fields[5]) if fields[5] and fields[5] != "-" else 0,
                "role": fields[6] if fields[6] else "",
                "roleArg": fields[7] if fields[7] else "",
                "enabled": not is_comment
            }
            shine_by_map[map_name].append(rec)

        elif table_name == "LinkTable":
            tabs = [t.strip() for t in line.split("\t")]
            rec_idx = next(i for i, t in enumerate(tabs) if "#recordin" in t or "#Record" in t)
            fields = tabs[rec_idx+9 : rec_idx+16]
            while len(fields) < 7:
                fields.append("")

            rec = {
                "argument": fields[0],
                "mapServer": fields[1],
                "mapClient": fields[2],
                "x": int(fields[3]) if fields[3] and fields[3] != "-" else 0,
                "y": int(fields[4]) if fields[4] and fields[4] != "-" else 0,
                "direct": int(fields[5]) if fields[5] and fields[5] != "-" else 0,
                "party": int(fields[6]) if fields[6] and fields[6] != "-" else 0,
                "enabled": not is_comment
            }
            link_records.append(rec)

    return shine_by_map, link_records

def parse_npc_action_txt():
    action_txt_path = os.path.join(NPC_DIR, "NPCAction.txt")
    with open(action_txt_path, "r", encoding="latin1") as f:
        lines = f.readlines()

    table = None
    conditions = []
    actions = []

    for line in lines:
        raw = line.split(";")[0].strip()
        if not raw:
            continue
        if raw.startswith("#Table"):
            table = raw.split()[1]
            continue
        if raw.startswith("#Record"):
            tokens = raw.split()
            if table == "NPCCondition":
                rec = {
                    "id": int(tokens[1]),
                    "conditionA": tokens[2],
                    "typeA": tokens[3],
                    "ax": int(tokens[4]),
                    "ay": int(tokens[5]),
                    "conditionB": tokens[6],
                    "typeB": tokens[7],
                    "bx": int(tokens[8]),
                    "by": int(tokens[9])
                }
                conditions.append(rec)
            elif table == "NPCAction":
                rec = {
                    "id": int(tokens[1]),
                    "actionIndex": tokens[2],
                    "type": tokens[3],
                    "x": int(tokens[4]),
                    "y": int(tokens[5])
                }
                actions.append(rec)

    return conditions, actions

def main():
    os.makedirs(MAPS_DIR, exist_ok=True)

    print("Parsing NPC.txt...")
    shine_by_map, link_records = parse_npc_txt()

    total_npcs = 0
    total_enabled_npcs = 0
    total_disabled_npcs = 0

    for map_name in sorted(shine_by_map.keys()):
        npcs = shine_by_map[map_name]
        total_npcs += len(npcs)
        total_enabled_npcs += sum(1 for n in npcs if n["enabled"])
        total_disabled_npcs += sum(1 for n in npcs if not n["enabled"])

        out_path = os.path.join(MAPS_DIR, f"{map_name}.json")
        with open(out_path, "w", encoding="utf-8") as f:
            json.dump(npcs, f, indent=2, ensure_ascii=False)

    print(f"Generated {len(shine_by_map)} map JSON files in {MAPS_DIR}")
    print(f"  Total NPCs: {total_npcs} (Enabled: {total_enabled_npcs}, Disabled: {total_disabled_npcs})")

    link_path = os.path.join(NPC_DIR, "LinkTable.json")
    with open(link_path, "w", encoding="utf-8") as f:
        json.dump(link_records, f, indent=2, ensure_ascii=False)

    total_enabled_links = sum(1 for l in link_records if l["enabled"])
    total_disabled_links = sum(1 for l in link_records if not l["enabled"])
    print(f"Generated LinkTable.json with {len(link_records)} records (Enabled: {total_enabled_links}, Disabled: {total_disabled_links})")

    print("Parsing NPCAction.txt...")
    conditions, actions = parse_npc_action_txt()
    action_data = {
        "conditions": conditions,
        "actions": actions
    }
    action_path = os.path.join(NPC_DIR, "NPCAction.json")
    with open(action_path, "w", encoding="utf-8") as f:
        json.dump(action_data, f, indent=2, ensure_ascii=False)

    print(f"Generated NPCAction.json with {len(conditions)} conditions and {len(actions)} actions")
    print("Done!")

if __name__ == "__main__":
    main()
