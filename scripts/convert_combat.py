#!/usr/bin/env python3
"""
Convert legacy World/Combat txt files into clean, modular JSON files with 100% data fidelity.
"""
import json
import os

CLASSES_SHORT = [
    "Fig", "Cfig", "War", "Gla", "Kni", "Cle", "Hcle", "Pal", "Hol", "Gua",
    "Arc", "Harc", "Sco", "Sha", "Ran", "Mag", "Wmag", "Enc", "Warl", "Wiz",
    "Jok", "Chs", "Cru", "Cls", "Ass", "Sen", "Sav"
]

def convert_chr_common():
    fpath = "data/World/Combat/ChrCommon.txt"
    common_data = {}
    stat_table = []
    free_stat_str = []
    free_stat_int = []
    free_stat_dex = []
    free_stat_con = []
    free_stat_men = []

    current_table = None

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line or line.startswith(";"):
                continue

            if line.lower().startswith("#table"):
                current_table = line.split()[1]
                continue

            if not line.startswith("#record"):
                continue

            # strip trailing comments from record
            clean_line = line.split(";")[0].strip()
            parts = clean_line.split()[1:]

            if current_table == "Common":
                key = parts[0]
                val = int(parts[1])
                # camelCase key
                key_camel = key[0].lower() + key[1:]
                common_data[key_camel] = val

            elif current_table == "StatTable":
                # level, NextExp, dummy, 27 classes, ExpLostAtPvP, CheckSum
                level = int(parts[0])
                next_exp = int(parts[1])
                # parts[2] is dummy (0)
                class_points = {}
                for idx, cname in enumerate(CLASSES_SHORT):
                    class_points[cname.lower()] = int(parts[3 + idx])
                exp_lost_pvp = int(parts[30])

                stat_table.append({
                    "level": level,
                    "nextExp": next_exp,
                    "expLostAtPvP": exp_lost_pvp,
                    "classStatPoints": class_points
                })

            elif current_table == "FreeStatStrTable":
                stat = int(parts[0])
                wc = int(parts[1])
                free_stat_str.append({"points": stat, "physicalAttack": wc})

            elif current_table == "FreeStatIntTable":
                stat = int(parts[0])
                ma = int(parts[1])
                free_stat_int.append({"points": stat, "magicAttack": ma})

            elif current_table == "FreeStatDexTable":
                stat = int(parts[0])
                th = int(parts[1])
                tb = int(parts[2])
                free_stat_dex.append({"points": stat, "aimRate": th, "evasionRate": tb})

            elif current_table == "FreeStatConTable":
                stat = int(parts[0])
                ac = int(parts[1])
                block = int(parts[2])
                max_hp = int(parts[3])
                free_stat_con.append({
                    "points": stat,
                    "physicalDef": ac,
                    "blockRate": block,
                    "maxHp": max_hp
                })

            elif current_table == "FreeStatMenTable":
                stat = int(parts[0])
                mr = int(parts[1])
                cri = int(parts[2])
                max_sp = int(parts[3])
                free_stat_men.append({
                    "points": stat,
                    "magicDef": mr,
                    "critRate": cri,
                    "maxSp": max_sp
                })

    with open("data/World/Combat/Common.json", "w", encoding="utf-8") as out:
        json.dump(common_data, out, indent=2)

    with open("data/World/Combat/ExpTable.json", "w", encoding="utf-8") as out:
        json.dump(stat_table, out, indent=2)

    free_stat_all = {
        "str": free_stat_str,
        "int": free_stat_int,
        "dex": free_stat_dex,
        "con": free_stat_con,
        "men": free_stat_men
    }
    with open("data/World/Combat/FreeStatTable.json", "w", encoding="utf-8") as out:
        json.dump(free_stat_all, out, indent=2)

    print("ChrCommon converted -> Common.json, ExpTable.json, FreeStatTable.json")

def convert_damage_by_angle():
    fpath = "data/World/Combat/DamageByAngle.txt"
    chr_records = []
    mob_records = []
    current_table = None

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line or line.startswith(";"):
                continue
            if line.lower().startswith("#table"):
                current_table = line.split()[1]
                continue
            if line.lower().startswith("#record"):
                parts = line.split()[1:]
                angle = int(parts[0])
                rate = int(parts[1])
                entry = {"angle": angle, "damageRate": rate}
                if current_table == "DamageByAngle_Chr":
                    chr_records.append(entry)
                elif current_table == "DamageByAngle_Mob":
                    mob_records.append(entry)

    out_data = {
        "character": chr_records,
        "monster": mob_records
    }
    with open("data/World/Combat/DamageByAngle.json", "w", encoding="utf-8") as out:
        json.dump(out_data, out, indent=2)
    print("DamageByAngle converted -> DamageByAngle.json")

def convert_damage_by_soul():
    fpath = "data/World/Combat/DamageBySoul.txt"
    records = []
    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line or line.startswith(";"):
                continue
            if line.lower().startswith("#record"):
                parts = line.split()[1:]
                demand_soul = int(parts[0])
                rates = [int(p) for p in parts[1:9]]
                records.append({
                    "demandSoul": demand_soul,
                    "rates": rates
                })

    with open("data/World/Combat/DamageBySoul.json", "w", encoding="utf-8") as out:
        json.dump(records, out, indent=2)
    print("DamageBySoul converted -> DamageBySoul.json")

def convert_exp_recalculation():
    fpath = "data/World/Combat/ExpRecalculation.txt"
    by_party = []
    by_diff = []
    std_dev = []
    current_table = None

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line = line.strip()
            if not line or line.startswith(";"):
                continue
            if line.lower().startswith("#table"):
                current_table = line.split()[1]
                continue
            if line.lower().startswith("#record"):
                clean = line.split(";")[0].strip()
                parts = clean.split()[1:]
                if current_table == "ByPartyMem":
                    by_party.append({
                        "partyMembers": int(parts[0]),
                        "bonusRate": int(parts[1])
                    })
                elif current_table == "ByLevelDiff":
                    by_diff.append({
                        "levelDiff": int(parts[0]),
                        "bonusRate": int(parts[1])
                    })
                elif current_table == "StandardDeviation":
                    std_dev.append({
                        "lowerBound": int(parts[0]),
                        "handicapRate": int(parts[1])
                    })

    out_data = {
        "byPartyMembers": by_party,
        "byLevelDiff": by_diff,
        "standardDeviation": std_dev
    }
    with open("data/World/Combat/ExpRecalculation.json", "w", encoding="utf-8") as out:
        json.dump(out_data, out, indent=2)
    print("ExpRecalculation converted -> ExpRecalculation.json")

def main():
    convert_chr_common()
    convert_damage_by_angle()
    convert_damage_by_soul()
    convert_exp_recalculation()

if __name__ == "__main__":
    main()
