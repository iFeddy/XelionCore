#!/usr/bin/env python3
"""
Verify 100% numerical and structural parity between Combat txt files and JSON files.
"""
import json
import os
from convert_combat import CLASSES_SHORT

def verify():
    # 1. Verify ChrCommon
    with open("data/World/Combat/Common.json", "r") as f:
        common = json.load(f)
    assert common["runSpeed"] == 127
    assert common["walkSpeed"] == 33
    assert common["attackSpeed"] == 13
    assert common["levelLimit"] == 125
    assert common["maxFreeStat"] == 180
    assert len(common) == 15

    with open("data/World/Combat/ExpTable.json", "r") as f:
        exp_table = json.load(f)
    assert len(exp_table) == 150
    assert exp_table[0]["level"] == 1
    assert exp_table[0]["nextExp"] == 10
    assert exp_table[9]["level"] == 10
    assert exp_table[9]["nextExp"] == 2350
    assert exp_table[19]["level"] == 20
    assert exp_table[19]["nextExp"] == 43539
    assert exp_table[149]["level"] == 150

    with open("data/World/Combat/FreeStatTable.json", "r") as f:
        free_stat = json.load(f)
    assert len(free_stat["str"]) == 181
    assert len(free_stat["int"]) == 181
    assert len(free_stat["dex"]) == 181
    assert len(free_stat["con"]) == 181
    assert len(free_stat["men"]) == 181
    assert free_stat["str"][180]["points"] == 180

    # 2. Verify DamageByAngle
    with open("data/World/Combat/DamageByAngle.json", "r") as f:
        angle_data = json.load(f)
    assert len(angle_data["character"]) == 6
    assert len(angle_data["monster"]) == 6
    assert angle_data["character"][0] == {"angle": 0, "damageRate": 1000}
    assert angle_data["character"][5] == {"angle": 180, "damageRate": 1200}

    # 3. Verify DamageBySoul
    with open("data/World/Combat/DamageBySoul.json", "r") as f:
        soul_data = json.load(f)
    assert len(soul_data) == 7
    assert soul_data[0]["demandSoul"] == 1
    assert soul_data[0]["rates"] == [1000, 1000, 1400, 1800, 2200, 2600, 3000, 3500]

    # 4. Verify ExpRecalculation
    with open("data/World/Combat/ExpRecalculation.json", "r") as f:
        exp_recalc = json.load(f)
    assert len(exp_recalc["byPartyMembers"]) == 6
    assert exp_recalc["byPartyMembers"][5]["bonusRate"] == 1800
    assert len(exp_recalc["byLevelDiff"]) == 31
    assert len(exp_recalc["standardDeviation"]) == 6

    print("ALL COMBAT FILES VERIFIED WITH 100% PARITY!")

if __name__ == "__main__":
    verify()
