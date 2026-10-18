#!/usr/bin/env python3
"""
Convert legacy World/Map/Field.txt into 4 modular JSON files with 100% data fidelity.
"""
import json
import os

def parse_int(v, default=0):
    try:
        return int(v)
    except (ValueError, TypeError):
        return default

def convert():
    fpath = "data/World/Map/Field.txt"
    field_list = []
    guild_tournaments = []
    instance_dungeons = []
    tutorial_map = []

    current_table = None

    with open(fpath, "r", encoding="latin1") as fp:
        for line in fp:
            line_s = line.strip()
            if not line_s or line_s.startswith(";"):
                continue

            if line_s.startswith("#Table"):
                current_table = line_s.split()[1]
                continue

            if line_s.startswith("#Record") or line_s.startswith("#record"):
                clean = line.split(";")[0].strip()
                # tab split
                parts = [p.strip() for p in clean.split("\t") if p.strip()][1:]

                if current_table == "FieldList":
                    map_id = parts[0]
                    sub_from = parts[1]
                    sub_to = parts[2]
                    serial = parse_int(parts[3])
                    name = parts[4].replace("#", " ")
                    kq = parse_int(parts[5])
                    xsize = parse_int(parts[6])
                    ysize = parse_int(parts[7])
                    immortal_sec = parse_int(parts[8])
                    script_name = parts[9]

                    # Combat flags
                    item_drop = parse_int(parts[10])
                    pker = parse_int(parts[11])
                    enemy_guild = parse_int(parts[12])
                    party = parse_int(parts[13])
                    monster = parse_int(parts[14])
                    summoned = parse_int(parts[15])
                    is_pk_kq = parse_int(parts[16]) == 1
                    is_free_pk = parse_int(parts[17]) == 1
                    is_party_battle = parse_int(parts[18]) == 1
                    name_hide = parse_int(parts[19]) == 1

                    # Teleport & system
                    link_in = parse_int(parts[20]) == 1
                    link_out = parse_int(parts[21]) == 1
                    system_map = parse_int(parts[22]) == 1

                    # Respawn
                    regen_city = parts[23]
                    r_xa, r_ya = parse_int(parts[24]), parse_int(parts[25])
                    r_xb, r_yb = parse_int(parts[26]), parse_int(parts[27])
                    r_xc, r_yc = parse_int(parts[28]), parse_int(parts[29])
                    r_xd, r_yd = parse_int(parts[30]), parse_int(parts[31])
                    regen_spot = parse_int(parts[32])

                    # Permissions
                    can_restart = parse_int(parts[33]) == 1
                    can_trade = parse_int(parts[34]) == 1
                    can_mini_house = parse_int(parts[35]) == 1
                    can_item = parse_int(parts[36]) == 1
                    can_skill = parse_int(parts[37]) == 1
                    chat = parse_int(parts[38]) == 1
                    shout = parse_int(parts[39]) == 1
                    can_booth = parse_int(parts[40]) == 1
                    can_produce = parse_int(parts[41]) == 1
                    can_ride = parse_int(parts[42]) == 1
                    can_stone = parse_int(parts[43]) == 1
                    can_party = parse_int(parts[44]) == 1

                    # Exp Loss & sublayer
                    exp_mob = parse_int(parts[45])
                    exp_ply = parse_int(parts[46])
                    usr_sub_layer = parse_int(parts[47])
                    checksum = parse_int(parts[48])
                    fiesta = parse_int(parts[49])

                    field_list.append({
                        "serial": serial,
                        "mapId": map_id,
                        "subFrom": sub_from,
                        "subTo": sub_to,
                        "name": name,
                        "isKq": kq != 0,
                        "kq": kq,
                        "width": xsize,
                        "height": ysize,
                        "immortalSec": immortal_sec,
                        "scriptName": script_name,
                        "combat": {
                            "itemDropByKilled": item_drop,
                            "pker": pker,
                            "enemyGuild": enemy_guild,
                            "party": party,
                            "monster": monster,
                            "summoned": summoned,
                            "isPkKq": is_pk_kq,
                            "isFreePk": is_free_pk,
                            "isPartyBattle": is_party_battle,
                            "nameHide": name_hide
                        },
                        "teleport": {
                            "linkIn": link_in,
                            "linkOut": link_out,
                            "systemMap": system_map
                        },
                        "respawn": {
                            "city": regen_city,
                            "spotCount": regen_spot,
                            "spots": [
                                {"x": r_xa, "y": r_ya},
                                {"x": r_xb, "y": r_yb},
                                {"x": r_xc, "y": r_yc},
                                {"x": r_xd, "y": r_yd}
                            ]
                        },
                        "permissions": {
                            "canRestart": can_restart,
                            "canTrade": can_trade,
                            "canMiniHouse": can_mini_house,
                            "canItem": can_item,
                            "canSkill": can_skill,
                            "chat": chat,
                            "shout": shout,
                            "canBooth": can_booth,
                            "canProduce": can_produce,
                            "canRide": can_ride,
                            "canStone": can_stone,
                            "canParty": can_party
                        },
                        "expLoss": {
                            "byMonster": exp_mob,
                            "byPlayer": exp_ply
                        },
                        "usrSubLayer": usr_sub_layer,
                        "checkSum": checksum,
                        "fiesta": fiesta
                    })

                elif current_table == "GuildTournamentMap":
                    map_id = parts[0]
                    gt_no = parse_int(parts[1])
                    stone_x = parse_int(parts[2])
                    stone_y = parse_int(parts[3])
                    stone_dir = parse_int(parts[4])
                    gld_a_x, gld_a_y = parse_int(parts[5]), parse_int(parts[6])
                    gld_b_x, gld_b_y = parse_int(parts[7]), parse_int(parts[8])

                    doors_a = []
                    doors_b = []
                    # 3 doors for A and B
                    # door 0: 9..12 (A), 13..16 (B)
                    doors_a.append({"x": parse_int(parts[9]), "y": parse_int(parts[10]), "dir": parse_int(parts[11]), "block": parts[12]})
                    doors_b.append({"x": parse_int(parts[13]), "y": parse_int(parts[14]), "dir": parse_int(parts[15]), "block": parts[16]})
                    # door 1: 17..20 (A), 21..24 (B)
                    doors_a.append({"x": parse_int(parts[17]), "y": parse_int(parts[18]), "dir": parse_int(parts[19]), "block": parts[20]})
                    doors_b.append({"x": parse_int(parts[21]), "y": parse_int(parts[22]), "dir": parse_int(parts[23]), "block": parts[24]})
                    # door 2: 25..28 (A), 29..32 (B)
                    doors_a.append({"x": parse_int(parts[25]), "y": parse_int(parts[26]), "dir": parse_int(parts[27]), "block": parts[28]})
                    doors_b.append({"x": parse_int(parts[29]), "y": parse_int(parts[30]), "dir": parse_int(parts[31]), "block": parts[32]})

                    guild_tournaments.append({
                        "mapId": map_id,
                        "gtMapNo": gt_no,
                        "stone": {"x": stone_x, "y": stone_y, "dir": stone_dir},
                        "guildARegion": {"x": gld_a_x, "y": gld_a_y},
                        "guildBRegion": {"x": gld_b_x, "y": gld_b_y},
                        "doorsA": doors_a,
                        "doorsB": doors_b
                    })

                elif current_table == "InstanceDungeon":
                    argument = parts[0]
                    id_no = parse_int(parts[1])
                    zone_number = parse_int(parts[2])
                    map_id = parts[3]
                    script_name = parts[4]
                    mode_id_lv = parse_int(parts[5])
                    entrance_type = parse_int(parts[6])
                    guild = parse_int(parts[7])
                    individual = parse_int(parts[8])
                    need_quest = parse_int(parts[9])
                    need_item = parts[10]
                    consume = parse_int(parts[11])
                    limit_time = parse_int(parts[12])
                    max_use_chr = parse_int(parts[13])

                    instance_dungeons.append({
                        "argument": argument,
                        "idNo": id_no,
                        "zoneNumber": zone_number,
                        "mapId": map_id,
                        "scriptName": script_name,
                        "modeIdLv": mode_id_lv,
                        "entranceType": entrance_type,
                        "guild": guild,
                        "individual": individual,
                        "needQuest": need_quest,
                        "needItem": need_item,
                        "consume": consume,
                        "limitTime": limit_time,
                        "maxUseChr": max_use_chr
                    })

                elif current_table == "TutorialMap":
                    tutorial_map.append({
                        "mapId": parts[0],
                        "maxChr": parse_int(parts[1]),
                        "scriptName": parts[2]
                    })

    # Save 4 files
    with open("data/World/Map/FieldList.json", "w", encoding="utf-8") as out:
        json.dump(field_list, out, indent=2)
    print(f"FieldList.json written ({len(field_list)} fields)")

    with open("data/World/Map/GuildTournament.json", "w", encoding="utf-8") as out:
        json.dump(guild_tournaments, out, indent=2)
    print(f"GuildTournament.json written ({len(guild_tournaments)} tournaments)")

    with open("data/World/Map/InstanceDungeon.json", "w", encoding="utf-8") as out:
        json.dump(instance_dungeons, out, indent=2)
    print(f"InstanceDungeon.json written ({len(instance_dungeons)} instance dungeons)")

    with open("data/World/Map/TutorialMap.json", "w", encoding="utf-8") as out:
        json.dump(tutorial_map, out, indent=2)
    print(f"TutorialMap.json written ({len(tutorial_map)} tutorials)")

if __name__ == "__main__":
    convert()
