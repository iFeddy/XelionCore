#!/usr/bin/env python3
"""
Modernize data/World/Quest:
- Remove obsolete dead-code Quest.txt (never read by Zone.exe)
- Convert QuestParser.txt (read by Zone.exe into gQuestParserData) into clean QuestParser.json
"""
import json
import os

def convert():
    txt_path = "data/World/Quest/QuestParser.txt"
    json_path = "data/World/Quest/QuestParser.json"

    class_map = {}
    with open(txt_path, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line or line.startswith(";"):
                continue
            parts = line.split()
            class_name = parts[0]
            class_id = int(parts[1])
            class_map[class_name] = class_id

    assert len(class_map) == 25, f"Expected 25 classes, got {len(class_map)}"

    with open(json_path, "w", encoding="utf-8") as out:
        json.dump(class_map, out, indent=2)

    print(f"QuestParser converted -> {json_path} (25 classes mapped)")

    # Remove dead code Quest.txt if exists
    quest_txt = "data/World/Quest/Quest.txt"
    if os.path.exists(quest_txt):
        os.remove(quest_txt)
        print("Removed dead legacy code Quest.txt")

if __name__ == "__main__":
    convert()
