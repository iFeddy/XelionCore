import glob
import os
import sys
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from extract_pdb_enums import get_enums

def build_full_opcode_catalog():
    # Load enums from all PDBs, merging values
    all_enums = {}
    for pdb_path in sorted(glob.glob('pdo/*.pdb')):
        enums = get_enums(pdb_path)
        for k, v in enums.items():
            if k not in all_enums:
                all_enums[k] = v
            else:
                existing_names = {x[0] for x in all_enums[k]}
                for fn, fv in v:
                    if fn not in existing_names:
                        all_enums[k].append((fn, fv))
                        existing_names.add(fn)

    # PROTOCOL_COMMAND: department list
    dept_enum = all_enums.get('PROTOCOL_COMMAND', [])
    print(f"Departments in PROTOCOL_COMMAND: {len(dept_enum)}")
    
    # Mapping department name/ID to enum name
    dept_to_enum = {
        1: ('NC_LOG', ['PROTOCOL_COMMAND_LOG', 'PROTOCOL_COMMAND_LOG_EASY_LOG_SYS']),
        2: ('NC_MISC', ['PROTOCOL_COMMAND_MISC']),
        3: ('NC_USER', ['PROTOCOL_COMMAND_USER']),
        4: ('NC_CHAR', ['PROTOCOL_COMMAND_CHAR']),
        5: ('NC_AVATAR', ['PROTOCOL_COMMAND_AVATAR']),
        6: ('NC_MAP', ['PROTOCOL_COMMAND_MAP']),
        7: ('NC_BRIEFINFO', ['PROTOCOL_COMMAND_BRIEFINFO']),
        8: ('NC_ACT', ['PROTOCOL_COMMAND_ACT']),
        9: ('NC_BAT', ['PROTOCOL_COMMAND_BAT']),
        10: ('NC_OPTOOL', ['PROTOCOL_COMMAND_OPTOOL']),
        12: ('NC_ITEM', ['PROTOCOL_COMMAND_ITEM']),
        13: ('NC_ITEMDB', ['PROTOCOL_COMMAND_ITEMDB']),
        14: ('NC_PARTY', ['PROTOCOL_COMMAND_PARTY']),
        15: ('NC_MENU', ['PROTOCOL_COMMAND_MENU']),
        16: ('NC_CHARSAVE', ['PROTOCOL_COMMAND_CHARSAVE']),
        17: ('NC_QUEST', ['PROTOCOL_COMMAND_QUEST']),
        18: ('NC_SKILL', ['PROTOCOL_COMMAND_SKILL']),
        19: ('NC_TRADE', ['PROTOCOL_COMMAND_TRADE']),
        20: ('NC_SOULSTONE', ['PROTOCOL_COMMAND_SOULSTONE']),
        21: ('NC_FRIEND', ['PROTOCOL_COMMAND_FRIEND']),
        22: ('NC_KQ', ['PROTOCOL_COMMAND_KINGDOMQUEST']),
        23: ('NC_WT', ['PROTOCOL_COMMAND_WEAPONTITLE']),
        24: ('NC_CT', ['PROTOCOL_COMMAND_CHARACTERTITLE']),
        25: ('NC_ANNOUNCE', ['PROTOCOL_COMMAND_ANNOUNCE']),
        26: ('NC_BOOTH', ['PROTOCOL_COMMAND_BOOTH']),
        27: ('NC_SCENARIO', ['PROTOCOL_COMMAND_SCENARIO']),
        28: ('NC_CHAR_OPTION', ['PROTOCOL_COMMAND_CHAR_OPTION']),
        29: ('NC_GUILD', ['PROTOCOL_COMMAND_GUILD']),
        31: ('NC_PRISON', ['PROTOCOL_COMMAND_PRISON']),
        32: ('NC_REPORT', ['PROTOCOL_COMMAND_REPORT']),
        35: ('NC_MINIHOUSE', ['PROTOCOL_COMMAND_MINIHOUSE']),
        36: ('NC_CHARGED', ['PROTOCOL_COMMAND_CHARGED']),
        37: ('NC_HOLY_PROMISE', ['PROTOCOL_COMMAND_HOLYPROMISE']),
        38: ('NC_GUILD_ACADEMY', ['PROTOCOL_COMMAND_GUILD_ACADEMY']),
        40: ('NC_PROMOTION', ['PROTOCOL_COMMAND_PROMOTION']),
        41: ('NC_INSTANCE_DUNGEON', ['PROTOCOL_COMMAND_INSTANCE_DUNGEON']),
        42: ('NC_CHAT_RESTRICT', ['PROTOCOL_COMMAND_CHAT_RESTRICT']),
        43: ('NC_DICE_TAISAI', ['PROTOCOL_COMMAND_DICE_TAISAI']),
        44: ('NC_RAID', ['PROTOCOL_COMMAND_RAID']),
        45: ('NC_USER_CONNECTION', ['PROTOCOL_USER_CONNECTION']),
        47: ('NC_GAMBLE', ['PROTOCOL_COMMAND_GAMBLE']),
        49: ('NC_COLLECT', ['PROTOCOL_COMMAND_COLLECT']),
        50: ('NC_SYSLOG', ['PROTOCOL_COMMAND_SYSLOG']),
        51: ('NC_MOVER', ['PROTOCOL_COMMAND_MOVER']),
        52: ('NC_EVENT', ['PROTOCOL_COMMAND_EVENT']),
        53: ('NC_PET', ['PROTOCOL_COMMAND_PET']),
    }
    
    opcodes = {}
    for dept_id, (dept_name, enum_names) in dept_to_enum.items():
        for ename in enum_names:
            fields = all_enums.get(ename, [])
            for fn, fv in fields:
                if fv == 0 or fn.endswith('_NULL') or fn.endswith('_NONE') or fn.endswith('_END'):
                    continue
                # op formula
                op = (dept_id << 10) | fv
                go_name = 'PROTO_' + fn
                opcodes[op] = (dept_id, fv, go_name, dept_name, ename)
                
    print(f"Total official opcodes generated: {len(opcodes)}")
    return opcodes

if __name__ == '__main__':
    ops = build_full_opcode_catalog()
