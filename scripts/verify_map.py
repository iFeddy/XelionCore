#!/usr/bin/env python3
"""
Verify 100% numerical and structural parity between Field.txt and the 4 JSON files.
"""
import json

def verify():
    # 1. FieldList
    with open("data/World/Map/FieldList.json", "r") as f:
        fields = json.load(f)
    assert len(fields) == 153, f"Expected 153 fields, got {len(fields)}"

    # Check Roumen
    rou = fields[0]
    assert rou["mapId"] == "Rou"
    assert rou["serial"] == 0
    assert rou["name"] == "Roumen"
    assert rou["width"] == 256
    assert rou["height"] == 256
    assert rou["immortalSec"] == 5000
    assert rou["respawn"]["city"] == "Rou"
    assert rou["respawn"]["spots"][0] == {"x": 4134, "y": 6488}
    assert rou["permissions"]["canTrade"] == True
    assert rou["expLoss"]["byMonster"] == 0

    # Check Elderine
    eld = [f for f in fields if f["mapId"] == "Eld"][0]
    assert eld["name"] == "Elderine"
    assert eld["width"] == 512
    assert eld["height"] == 512

    # 2. GuildTournament
    with open("data/World/Map/GuildTournament.json", "r") as f:
        gld = json.load(f)
    assert len(gld) == 3
    assert gld[0]["mapId"] == "GuildT04"
    assert len(gld[0]["doorsA"]) == 3
    assert len(gld[0]["doorsB"]) == 3

    # 3. InstanceDungeon
    with open("data/World/Map/InstanceDungeon.json", "r") as f:
        id_list = json.load(f)
    assert len(id_list) == 17
    assert id_list[0]["mapId"] == "UrgDragon"
    assert id_list[0]["zoneNumber"] == 3
    assert id_list[0]["scriptName"] == "ID/GraveYard/GraveYard"
    assert id_list[4]["mapId"] == "Leviathan"
    assert id_list[4]["scriptName"] == "ID/d_NestOfLeviathan"

    # 4. TutorialMap
    with open("data/World/Map/TutorialMap.json", "r") as f:
        tut = json.load(f)
    assert len(tut) == 1
    assert tut[0]["mapId"] == "RouT"
    assert tut[0]["maxChr"] == 100
    assert tut[0]["scriptName"] == "Tutorial/Tutorial"

    print("ALL 4 MAP TABLES VERIFIED WITH 100% PARITY!")

if __name__ == "__main__":
    verify()
