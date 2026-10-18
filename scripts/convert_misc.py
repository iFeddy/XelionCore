#!/usr/bin/env python3
"""
Convert legacy World/Misc txt files into clean, modular JSON files with 100% data fidelity.
"""
from collections import defaultdict
import json
import os

def parse_int(v, default=0):
    try:
        return int(v)
    except (ValueError, TypeError):
        return default

def convert_recall_coord():
    fpath = "data/World/Misc/RecallCoord.txt"
    records = []
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line_s = line.strip()
            if not line_s.startswith("#Record"):
                continue
            clean = line.split(";")[0].strip()
            parts = clean.split()
            item_index = parts[1]
            item_ident = parse_int(parts[2])
            map_name = parts[3]
            x = parse_int(parts[4])
            y = parse_int(parts[5])
            records.append({
                "itemIndex": item_index,
                "itemIdent": item_ident,
                "mapName": map_name,
                "x": x,
                "y": y
            })

    out_path = "data/World/Misc/RecallCoord.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(records, out, indent=2)
    print(f"RecallCoord converted -> {out_path} ({len(records)} entries)")

def convert_stealth_matrix():
    fpath = "data/World/Misc/SubLayerInteract.txt"
    visibility = {}
    can_attack = {}
    current_table = None

    layers = ["Base", "RangerStealth", "Observer", "AdminHide", "GMDoor", "GMPlayer"]

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line_s = line.strip()
            if not line_s or line_s.startswith(";"):
                continue
            if line_s.startswith("#Table"):
                current_table = line_s.split()[1]
                continue
            if line_s.startswith("#Record"):
                clean = line.split(";")[0].strip()
                parts = clean.split()[1:]
                # parts[0] is CheckSumA (0..5), parts[1] is Layer Name
                layer_name = parts[1]
                # parts[2..7] correspond to layers Base..GMPlayer
                row_map = {}
                for idx, target_layer in enumerate(layers):
                    row_map[target_layer] = parse_int(parts[2 + idx]) == 1

                if current_table == "SendMyBrief":
                    visibility[layer_name] = row_map
                elif current_table == "CanAttack":
                    can_attack[layer_name] = row_map

    out_data = {
        "visibility": visibility,
        "canAttack": can_attack
    }
    out_path = "data/World/Misc/StealthMatrix.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(out_data, out, indent=2)
    print(f"StealthMatrix converted -> {out_path}")

def convert_skill_mechanics():
    fpath = "data/World/Misc/MiscDataTable.txt"
    skill_breed_mob = []
    polymorph = []
    neutralization = []
    abstate_reset = []
    expand_skill = []

    current_table = None
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line_s = line.strip()
            if not line_s or line_s.startswith(";"):
                continue
            if line_s.startswith("#Table"):
                current_table = line_s.split()[1]
                continue
            if line_s.startswith("#Record"):
                clean = line.split(";")[0].strip()
                parts = clean.split()[1:]

                if current_table == "SkillBreedMob":
                    skill_breed_mob.append({
                        "serial": parse_int(parts[0]),
                        "skillIndex": parts[1],
                        "mobIndex": parts[2],
                        "num": parse_int(parts[3]),
                        "moveType": parts[4],
                        "target": parts[5],
                        "range": parse_int(parts[6]),
                        "shape": parts[7],
                        "lifeTime": parse_int(parts[8]),
                        "abState": parts[9],
                        "filter": parts[10],
                        "aggro": parse_int(parts[11]),
                        "interval": parse_int(parts[12]),
                        "activeTime": parse_int(parts[13]),
                        "subSkill": parts[14],
                        "applyTarget": parts[15],
                        "disappear": parts[16],
                        "option": parts[17] if len(parts) > 17 else "-"
                    })
                elif current_table == "Polymorph":
                    polymorph.append({
                        "serial": parse_int(parts[0]),
                        "skill": parts[1],
                        "mobIndex": parts[2],
                        "canMove": parse_int(parts[3]) == 1,
                        "canAttack": parse_int(parts[4]) == 1
                    })
                elif current_table == "Neutralization":
                    neutralization.append({
                        "serial": parse_int(parts[0]),
                        "skill": parts[1],
                        "weapon": parse_int(parts[2]),
                        "shield": parse_int(parts[3]),
                        "body": parse_int(parts[4]),
                        "leg": parse_int(parts[5])
                    })
                elif current_table == "AbStateReset":
                    abstate_reset.append({
                        "serial": parse_int(parts[0]),
                        "abState": parts[1],
                        "run": parse_int(parts[2]) == 1,
                        "walk": parse_int(parts[3]) == 1,
                        "attack": parse_int(parts[4]) == 1,
                        "attacked": parse_int(parts[5]) == 1,
                        "equipWhere": parts[6],
                        "equipClass": parts[7]
                    })
                elif current_table == "ExpandSkill":
                    expand_skill.append({
                        "serial": parse_int(parts[0]),
                        "skill": parts[1],
                        "condition": parts[2],
                        "dmg": parse_int(parts[3]),
                        "abState": parts[4],
                        "critical": parse_int(parts[5])
                    })

    out_data = {
        "skillBreedMobs": skill_breed_mob,
        "polymorph": polymorph,
        "neutralization": neutralization,
        "abStateReset": abstate_reset,
        "expandSkill": expand_skill
    }
    out_path = "data/World/Misc/SkillMechanics.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(out_data, out, indent=2)
    print(f"SkillMechanics converted -> {out_path} ({len(skill_breed_mob)} breed mobs, {len(polymorph)} poly, {len(neutralization)} neut, {len(abstate_reset)} resets, {len(expand_skill)} expands)")

def convert_mob_chat():
    fpath = "data/World/Misc/MobChat.txt"
    table_keys = {
        "PIECE": "peace",
        "ATTACK": "attack",
        "DAMAGED": "damaged",
        "DEAD": "dead",
        "HELPMAIN": "helpMain",
        "HELPSUB": "helpSub",
        "RACERESCUE": "raceRescue"
    }
    chats = defaultdict(list)

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line_s = line.strip()
            if not line_s.startswith("#RecordIn"):
                continue
            clean = line.split(";")[0].strip()
            parts = [p.strip() for p in clean.split("\t") if p.strip()]
            # parts: [#RecordIn, Table, MobIndex, Rate0, Rate1, Script0, Script1, Script2, Script3]
            table_name = parts[1]
            key = table_keys.get(table_name, table_name.lower())
            mob_index = parts[2]
            rate0 = parse_int(parts[3])
            rate1 = parse_int(parts[4]) if len(parts) > 4 else 0
            scripts = []
            for s_idx in range(5, len(parts)):
                text = parts[s_idx]
                if text and text != "-":
                    # unquote if quoted
                    if text.startswith('"') and text.endswith('"'):
                        text = text[1:-1]
                    scripts.append(text)

            chats[key].append({
                "mob": mob_index,
                "rate0": rate0,
                "rate1": rate1,
                "messages": scripts
            })

    out_path = "data/World/Misc/MobChat.json"
    with open(out_path, "w", encoding="utf-8") as out:
        json.dump(chats, out, indent=2)
    total_chats = sum(len(v) for v in chats.values())
    print(f"MobChat converted -> {out_path} ({total_chats} entries across {len(chats)} categories)")

def main():
    convert_recall_coord()
    convert_stealth_matrix()
    convert_skill_mechanics()
    convert_mob_chat()

if __name__ == "__main__":
    main()
