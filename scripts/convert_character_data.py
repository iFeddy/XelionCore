#!/usr/bin/env python3
import json
import os
import re

CLASS_NAMES = {
    1: "Fighter",
    6: "Cleric",
    11: "Archer",
    16: "Mage",
    21: "Joker",
    26: "Sentinel",
}


def read_shine_text(filepath):
    with open(filepath, "rb") as bf:
        raw = bf.read()
    if raw.startswith(b"\xff\xfe") or raw.startswith(b"\xfe\xff"):
        return raw.decode("utf-16")
    for enc in ("utf-8", "cp949", "euc-kr", "latin-1"):
        try:
            return raw.decode(enc)
        except UnicodeDecodeError:
            continue
    return raw.decode("latin-1", errors="replace")


def parse_default_character_data(filepath):
    classes = {}
    class_skills = {}
    all_starter_quests = []

    for cid, cname in CLASS_NAMES.items():
        classes[str(cid)] = {
            "class_id": cid,
            "class_name": cname,
            "start_map": "RouN",
            "start_x": 0,
            "start_y": 0,
            "hp": 0,
            "sp": 0,
            "hp_soul": 0,
            "sp_soul": 0,
            "money": 0,
            "init_level": 1,
            "init_exp": 0,
            "max_shape": {"max_face": 3, "max_hair": 3, "max_hair_color": 3},
            "starter_items": [],
            "shortcuts": [],
            "starter_skills": [],
            "starter_quests": [],
        }

    text = read_shine_text(filepath)
    in_define = False
    for raw_line in text.splitlines():
        line = raw_line.strip()
        if not line or line.startswith(";"):
            continue

        upper_line = line.upper()
        if upper_line.startswith("#DEFINE"):
            in_define = True
            continue
        if upper_line.startswith("#ENDDEFINE"):
            in_define = False
            continue
        if in_define or upper_line.startswith("#END"):
            continue

        comment = ""
        if ";" in line:
            parts = line.split(";", 1)
            line = parts[0].strip()
            comment = parts[1].strip()
        if not line:
            continue

        tokens = [t.strip().strip('"') for t in re.split(r"[\s,]+", line) if t.strip()]
        if not tokens:
            continue

        directive = tokens[0].upper()
        args = tokens[1:]

        if directive == "CHARACTER" and len(args) >= 10:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["start_map"] = args[1]
                classes[cid]["start_x"] = int(args[2])
                classes[cid]["start_y"] = int(args[3])
                classes[cid]["hp"] = int(args[4])
                classes[cid]["sp"] = int(args[5])
                classes[cid]["hp_soul"] = int(args[6])
                classes[cid]["sp_soul"] = int(args[7])
                classes[cid]["money"] = int(args[8])
                classes[cid]["init_level"] = int(args[9])
                if len(args) >= 11:
                    classes[cid]["init_exp"] = int(args[10])

        elif directive == "MAXSHAPE" and len(args) >= 4:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["max_shape"] = {
                    "max_face": int(args[1]),
                    "max_hair": int(args[2]),
                    "max_hair_color": int(args[3]),
                }

        elif directive == "ITEM" and len(args) >= 3:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["starter_items"].append({
                    "item_id": int(args[1]),
                    "amount": int(args[2]),
                    "name": comment or f"Item {args[1]}",
                })

        elif directive == "SHORTCUT" and len(args) >= 4:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["shortcuts"].append({
                    "slot": int(args[1]),
                    "code_no": int(args[2]),
                    "value": int(args[3]),
                })

        elif directive == "SKILL" and len(args) >= 2:
            cid = str(int(args[0]))
            skill_id = int(args[1])
            class_skills.setdefault(cid, [])
            if skill_id not in class_skills[cid]:
                class_skills[cid].append(skill_id)
            if cid in classes and skill_id not in classes[cid]["starter_skills"]:
                classes[cid]["starter_skills"].append(skill_id)

        elif directive == "QUEST" and len(args) >= 2:
            cid_int = int(args[0])
            cid = str(cid_int)
            quest_id = int(args[1])
            all_starter_quests.append({"class_id": cid_int, "quest_id": quest_id})
            if cid in classes and quest_id not in classes[cid]["starter_quests"]:
                classes[cid]["starter_quests"].append(quest_id)

    return {
        "classes": classes,
        "class_skills": class_skills,
        "starter_quests": all_starter_quests,
    }


def parse_tutorial_character_data(filepath):
    classes = {}
    text = read_shine_text(filepath)
    in_define = False

    for raw_line in text.splitlines():
        line = raw_line.strip()
        if not line or line.startswith(";"):
            continue

        upper_line = line.upper()
        if upper_line.startswith("#DEFINE"):
            in_define = True
            continue
        if upper_line.startswith("#ENDDEFINE"):
            in_define = False
            continue
        if in_define or upper_line.startswith("#END"):
            continue

        comment = ""
        if ";" in line:
            parts = line.split(";", 1)
            line = parts[0].strip()
            comment = parts[1].strip()
        if not line:
            continue

        tokens = [t.strip().strip('"') for t in re.split(r"[\s,]+", line) if t.strip()]
        if not tokens:
            continue

        directive = tokens[0].upper()
        args = tokens[1:]

        if directive == "TUTORIAL_CHARACTER" and len(args) >= 9:
            cid_int = int(args[0])
            cid = str(cid_int)
            classes[cid] = {
                "class_id": cid_int,
                "class_name": CLASS_NAMES.get(cid_int, f"Class {cid_int}"),
                "start_map": args[1],
                "start_x": int(args[2]),
                "start_y": int(args[3]),
                "hp": int(args[4]),
                "sp": int(args[5]),
                "hp_soul": int(args[6]),
                "sp_soul": int(args[7]),
                "money": int(args[8]),
                "init_level": 1,
                "init_exp": 0,
                "max_shape": {"max_face": 3, "max_hair": 3, "max_hair_color": 3},
                "starter_items": [],
                "tutorial_items": [],
                "shortcuts": [],
                "starter_skills": [],
                "starter_quests": [],
            }
        elif directive == "DEFAULT_ITEM" and len(args) >= 3:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["starter_items"].append({
                    "item_id": int(args[1]),
                    "amount": int(args[2]),
                    "name": comment or f"Item {args[1]}",
                })
        elif directive == "TUTORIAL_ITEM" and len(args) >= 3:
            cid = str(int(args[0]))
            if cid in classes:
                classes[cid]["tutorial_items"].append({
                    "item_id": int(args[1]),
                    "amount": int(args[2]),
                    "name": comment or f"Item {args[1]}",
                })

    return {
        "classes": classes,
        "tutorials": classes,
    }


if __name__ == "__main__":
    def_path = (
        "data/DefaultCharacterData.txt"
        if os.path.exists("data/DefaultCharacterData.txt")
        else "../AstraServer/Server/9Data/Shine/DefaultCharacterData.txt"
    )
    def_data = parse_default_character_data(def_path)
    with open("data/default_character_data.json", "w", encoding="utf-8") as out:
        json.dump(def_data, out, indent=2, ensure_ascii=False)
    print(f"Generated data/default_character_data.json ({len(def_data['classes'])} classes).")

    tut_path = (
        "data/TutorialCharacterData.txt"
        if os.path.exists("data/TutorialCharacterData.txt")
        else "../AstraServer/Server/9Data/Shine/TutorialCharacterData.txt"
    )
    tut_data = parse_tutorial_character_data(tut_path)
    with open("data/tutorial_character_data.json", "w", encoding="utf-8") as out:
        json.dump(tut_data, out, indent=2, ensure_ascii=False)
    print(f"Generated data/tutorial_character_data.json ({len(tut_data['classes'])} tutorial classes).")
