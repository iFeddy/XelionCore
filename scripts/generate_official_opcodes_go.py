import os
import sys
import glob

sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from extract_pdb_enums import get_enums

def main():
    all_enums = {}
    for pdb_path in sorted(glob.glob('pdo/*.pdb')):
        enums = get_enums(pdb_path)
        for k, v in enums.items():
            if k not in all_enums:
                all_enums[k] = list(v)
            else:
                existing_names = {x[0] for x in all_enums[k]}
                for fn, fv in v:
                    if fn not in existing_names:
                        all_enums[k].append((fn, fv))
                        existing_names.add(fn)

    dept_definitions = [
        (1, "NC_LOG", "Auditoría y Logs del Sistema", ['PROTOCOL_COMMAND_LOG', 'PROTOCOL_COMMAND_LOG_EASY_LOG_SYS']),
        (2, "NC_MISC", "Utilidades de Red, Tiempos y Semillas XOR", ['PROTOCOL_COMMAND_MISC']),
        (3, "NC_USER", "Autenticación de Cuentas y Selección de Mundo", ['PROTOCOL_COMMAND_USER']),
        (4, "NC_CHAR", "Selección de Personaje y Datos In-Game", ['PROTOCOL_COMMAND_CHAR']),
        (5, "NC_AVATAR", "Creación y Eliminación de Personajes (Avatar)", ['PROTOCOL_COMMAND_AVATAR']),
        (6, "NC_MAP", "Transición y Navegación de Mapas / Zonas", ['PROTOCOL_COMMAND_MAP']),
        (7, "NC_BRIEFINFO", "Spawns, Despawns y Visibilidad en Vista de Campo", ['PROTOCOL_COMMAND_BRIEFINFO']),
        (8, "NC_ACT", "Movimiento, Acciones, Gestos, Tiendas y Menús", ['PROTOCOL_COMMAND_ACT']),
        (9, "NC_BAT", "Combate, Habilidades, Buffs y Debuffs", ['PROTOCOL_COMMAND_BAT']),
        (10, "NC_OPTOOL", "Herramientas de Operación y Administración", ['PROTOCOL_COMMAND_OPTOOL']),
        (12, "NC_ITEM", "Inventario, Ítems, Uso, Comercio y Refinado", ['PROTOCOL_COMMAND_ITEM']),
        (13, "NC_ITEMDB", "Persistencia S2S de Ítems en Base de Datos", ['PROTOCOL_COMMAND_ITEMDB']),
        (14, "NC_PARTY", "Grupos y Partys", ['PROTOCOL_COMMAND_PARTY']),
        (15, "NC_MENU", "Menús e Interfaces de NPCs", ['PROTOCOL_COMMAND_MENU']),
        (16, "NC_CHARSAVE", "Guardado Periódico y Persistencia de Personaje", ['PROTOCOL_COMMAND_CHARSAVE']),
        (17, "NC_QUEST", "Misiones y Progresión", ['PROTOCOL_COMMAND_QUEST']),
        (18, "NC_SKILL", "Habilidades y Aprendizaje", ['PROTOCOL_COMMAND_SKILL']),
        (19, "NC_TRADE", "Intercambio Seguro entre Jugadores", ['PROTOCOL_COMMAND_TRADE']),
        (20, "NC_SOULSTONE", "Piedras de Alma (HP/SP Soul Stones)", ['PROTOCOL_COMMAND_SOULSTONE']),
        (21, "NC_FRIEND", "Lista y Sistema de Amigos", ['PROTOCOL_COMMAND_FRIEND']),
        (22, "NC_KQ", "Kingdom Quests (KQ)", ['PROTOCOL_COMMAND_KINGDOMQUEST']),
        (23, "NC_WT", "Títulos de Armas (Weapon Title)", ['PROTOCOL_COMMAND_WEAPONTITLE']),
        (24, "NC_CT", "Títulos de Personaje (Character Title)", ['PROTOCOL_COMMAND_CHARACTERTITLE']),
        (25, "NC_ANNOUNCE", "Anuncios Globales y Avisos de Servidor", ['PROTOCOL_COMMAND_ANNOUNCE']),
        (26, "NC_BOOTH", "Puestos de Venta Personal (Vendor Booths)", ['PROTOCOL_COMMAND_BOOTH']),
        (27, "NC_SCENARIO", "Eventos y Escenarios Scriptados", ['PROTOCOL_COMMAND_SCENARIO']),
        (28, "NC_CHAR_OPTION", "Configuración de Teclas, QuickBar y Ventanas", ['PROTOCOL_COMMAND_CHAR_OPTION']),
        (29, "NC_GUILD", "Gremios y Torneos de Gremios", ['PROTOCOL_COMMAND_GUILD']),
        (31, "NC_PRISON", "Prisión y Sanciones de Moderación", ['PROTOCOL_COMMAND_PRISON']),
        (32, "NC_REPORT", "Reportes de Jugadores", ['PROTOCOL_COMMAND_REPORT']),
        (35, "NC_MINIHOUSE", "Mini-Casas de Descanso", ['PROTOCOL_COMMAND_MINIHOUSE']),
        (36, "NC_CHARGED", "Buffs y Servicios Pagos / Premium", ['PROTOCOL_COMMAND_CHARGED']),
        (37, "NC_HOLY_PROMISE", "Sistema de Padrinazgo y Alianzas (Holy Promise)", ['PROTOCOL_COMMAND_HOLYPROMISE']),
        (38, "NC_GUILD_ACADEMY", "Academia de Gremios", ['PROTOCOL_COMMAND_GUILD_ACADEMY']),
        (40, "NC_PROMOTION", "Promociones y Recompensas Web", ['PROTOCOL_COMMAND_PROMOTION']),
        (41, "NC_INSTANCE_DUNGEON", "Mazmorras por Instancia", ['PROTOCOL_COMMAND_INSTANCE_DUNGEON']),
        (42, "NC_CHAT_RESTRICT", "Restricciones de Chat y Silenciamientos", ['PROTOCOL_COMMAND_CHAT_RESTRICT']),
        (43, "NC_DICE_TAISAI", "Minijuego de Dados / Taisai", ['PROTOCOL_COMMAND_DICE_TAISAI']),
        (44, "NC_RAID", "Bandas y Raids", ['PROTOCOL_COMMAND_RAID']),
        (45, "NC_USER_CONNECTION", "Control de Sesiones Concurrentes y Traspasos", ['PROTOCOL_USER_CONNECTION']),
        (47, "NC_GAMBLE", "Casa de Apuestas y Lotería", ['PROTOCOL_COMMAND_GAMBLE']),
        (49, "NC_COLLECT", "Álbum de Colección y Cartas", ['PROTOCOL_COMMAND_COLLECT']),
        (50, "NC_SYSLOG", "Métricas y Syslog del Servidor", ['PROTOCOL_COMMAND_SYSLOG']),
        (51, "NC_MOVER", "Monturas y Vehículos (Movers)", ['PROTOCOL_COMMAND_MOVER']),
        (52, "NC_EVENT", "Eventos Globales del Servidor", ['PROTOCOL_COMMAND_EVENT']),
        (53, "NC_PET", "Mascotas y Crecimiento (Pets)", ['PROTOCOL_COMMAND_PET']),
    ]

    lines = []
    lines.append("package protocol")
    lines.append("")
    lines.append("// Catálogo Oficial de Protocolos y Opcodes de Fiesta Online (XelionCore).")
    lines.append("// Extraído directamente desde los flujos TPI/PDB oficiales (Zone.pdb, WorldManager.pdb, Account.pdb, etc.).")
    lines.append("//")
    lines.append("// Fórmula bitfield oficial NETCOMMAND de 16 bits:")
    lines.append("//   protocol = (department << 10) | command")
    lines.append("//")
    lines.append("// Cero Omisiones, Cero Resúmenes, 100% Fidelidad Binaria.")
    lines.append("")
    lines.append("// MakeOpCode packs a department (0..63) and a command (0..1023) into a uint16 opcode.")
    lines.append("func MakeOpCode(dept uint8, cmd uint16) uint16 {")
    lines.append("\treturn (uint16(dept) << 10) | (cmd & 0x3FF)")
    lines.append("}")
    lines.append("")
    lines.append("// ExtractHeader returns the department (category) from an opcode.")
    lines.append("func ExtractHeader(opcode uint16) uint8 {")
    lines.append("\treturn uint8(opcode >> 10)")
    lines.append("}")
    lines.append("")
    lines.append("// ExtractType returns the command ID from an opcode.")
    lines.append("func ExtractType(opcode uint16) uint16 {")
    lines.append("\treturn opcode & 0x3FF")
    lines.append("}")
    lines.append("")
    lines.append("const (")

    seen_names = set()
    total_ops = 0

    for dept_id, dept_code, dept_desc, enum_names in dept_definitions:
        lines.append(f"\t// ========================================================================")
        lines.append(f"\t// Departamento {dept_id:02d}: {dept_code} ({dept_desc})")
        lines.append(f"\t// ========================================================================")
        
        dept_ops = []
        for ename in enum_names:
            fields = all_enums.get(ename, [])
            for fn, fv in fields:
                if fv == 0 or fn.endswith('_NULL') or fn.endswith('_NONE') or fn.endswith('_END'):
                    continue
                go_name = 'PROTO_' + fn
                dept_ops.append((fv, go_name, ename))

        dept_ops.sort(key=lambda x: x[0])
        
        for cmd_id, go_name, ename in dept_ops:
            if go_name in seen_names:
                continue
            seen_names.add(go_name)
            op = (dept_id << 10) | cmd_id
            lines.append(f"\t{go_name:<55} uint16 = ({dept_id:2d} << 10) | {cmd_id:<3d} // 0x{op:04X}")
            total_ops += 1

        lines.append("")

    lines.append(")")
    lines.append("")

    # Preserved types: ServerError, WorldStatus, CreateCharError, Job
    types_block = """// ServerError codes returned in PROTO_NC_USER_LOGINFAIL_ACK packets.
type ServerError uint16

const (
	ErrException          ServerError = 0x42 // Unknown Error / Exception
	ErrDatabase           ServerError = 0x43 // DB Error
	ErrWrongRegion        ServerError = 0x44 // Auth Failed / Wrong Region
	ErrInvalidCredentials ServerError = 0x45 // Please check ID or password
	ErrDatabaseError      ServerError = 0x46 // DB Error
	ErrBlocked            ServerError = 0x47 // The ID has been blocked
	ErrServerMaintenance  ServerError = 0x48 // World server maintenance
	ErrTimeout            ServerError = 0x49 // Timeout
	ErrLoginFailed        ServerError = 0x4A // Login Failed
	ErrAgreementMissing   ServerError = 0x4B // Please accept the agreement
)

// WorldStatus states for server selection display.
type WorldStatus uint8

const (
	WorldStatus_Offline         WorldStatus = 0
	WorldStatus_Maintenance     WorldStatus = 1
	WorldStatus_EmptyServer     WorldStatus = 2
	WorldStatus_Reserved        WorldStatus = 3
	WorldStatus_OfflineUnkError WorldStatus = 4
	WorldStatus_Full            WorldStatus = 5
	WorldStatus_Low             WorldStatus = 6
	WorldStatus_Medium          WorldStatus = 9
	WorldStatus_High            WorldStatus = 10
)

// CreateCharError codes returned in PROTO_NC_AVATAR_CREATEFAIL_ACK packets.
type CreateCharError uint16

const (
	CreateCharError_FailedToCreate CreateCharError = 130 // Failed to create the char
	CreateCharError_WrongClass     CreateCharError = 131 // You have chosen wrong class for your avatar
	CreateCharError_NameTaken      CreateCharError = 132 // The name is already taken by another user
	CreateCharError_ErrorInMaxSlot CreateCharError = 133 // There is an error in the Max slot Number
	CreateCharError_NameInUse      CreateCharError = 385 // The name is already in use
)

// Job identifies character classes in Fiesta Online.
type Job byte

const (
	Job_Fighter   Job = 1
	Job_Cleric    Job = 6
	Job_Archer    Job = 11
	Job_Mage      Job = 16
	Job_Trickster Job = 21
	Job_Crusader  Job = 26
)
"""
    lines.append(types_block)

    raw_aliases = [
        ("PROTO_NC_MISC_SEED_CMD", "PROTO_NC_MISC_SEED_ACK"),
        ("PROTO_NC_MISC_TIME_SYNC_ACK", "PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD"),
        ("PROTO_NC_MISC_UPDATE_TIME", "PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD"),
        ("PROTO_NC_MISC_WORLD_CLIENT_KEY_REQ", "PROTO_NC_USER_LOGINWORLD_REQ"),
        ("PROTO_NC_MISC_WORLD_CLIENT_KEY_ACK", "PROTO_NC_USER_LOGIN_DB"),
        ("PROTO_NC_USER_FILE_CHECK_REQ", "PROTO_NC_USER_XTRAP_REQ"),
        ("PROTO_NC_USER_FILE_CHECK_ACK", "PROTO_NC_USER_XTRAP_ACK"),
        ("PROTO_NC_CHAR_CLIENT_CHARLIST_ACK", "PROTO_NC_USER_LOGINWORLD_ACK"),
        ("PROTO_NC_USER_BACKTOWORLDLIST_REQ", "PROTO_NC_USER_WILL_WORLD_SELECT_REQ"),
        ("PROTO_NC_USER_BACKTOWORLDLIST_ACK", "PROTO_NC_USER_WILL_WORLD_SELECT_ACK"),
        ("PROTO_NC_USER_LOGIN_WITH_OTP_LEGACY", "PROTO_NC_USER_GER_LOGIN_REQ"),
        ("PROTO_NC_USER_WORLDSTATUS_REQ", "PROTO_NC_USER_WORLD_STATUS_REQ"),
        ("PROTO_NC_USER_WORLDSTATUS_ACK", "PROTO_NC_USER_WORLD_STATUS_ACK"),
        ("PROTO_NC_CHAR_SERVER_IP_ACK", "PROTO_NC_CHAR_LOGIN_ACK"),
        ("PROTO_NC_CHAR_POINTS_ACK", "PROTO_NC_CHAR_STAT_REMAINPOINT_CMD"),
        ("PROTO_NC_CHAR_CREATE_REQ", "PROTO_NC_AVATAR_CREATE_REQ"),
        ("PROTO_NC_CHAR_DELETE_REQ", "PROTO_NC_AVATAR_ERASE_REQ"),
        ("PROTO_NC_CHAR_CREATE_FAIL_ACK", "PROTO_NC_AVATAR_CREATEFAIL_ACK"),
        ("PROTO_NC_CHAR_CREATE_ACK", "PROTO_NC_AVATAR_CREATESUCC_ACK"),
        ("PROTO_NC_CHAR_DELETE_ACK", "PROTO_NC_AVATAR_ERASESUCC_ACK"),
        ("PROTO_NC_CHAR_MAPLOGIN_ACK", "PROTO_NC_MAP_LOGIN_ACK"),
        ("PROTO_NC_MAP_LOGINCOMPLETE_CMD", "PROTO_NC_MAP_LOGINFAIL_ACK"),
        ("PROTO_NC_MAP_CHANGE_CMD", "PROTO_NC_MAP_LINKSAME_CMD"),
        ("PROTO_NC_MAP_CHANGE_ZONE_CMD", "PROTO_NC_MAP_LINKOTHER_CMD"),
        ("PROTO_NC_MAP_TELEPORT_REQ", "PROTO_NC_MAP_TOWNPORTAL_REQ"),
        ("PROTO_NC_MAP_TELEPORT_ACK", "PROTO_NC_MAP_TOWNPORTAL_ACK"),
        ("PROTO_NC_BRIEFINFO_UNSHARE_EQUIP_ACK", "PROTO_NC_BRIEFINFO_UNEQUIP_CMD"),
        ("PROTO_NC_BRIEFINFO_SHARE_EQUIP_ACK", "PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD"),
        ("PROTO_NC_BRIEFINFO_CHARACTER_ACK", "PROTO_NC_BRIEFINFO_CHARACTER_CMD"),
        ("PROTO_NC_BRIEFINFO_LEAVE_CMD", "PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD"),
        ("PROTO_NC_ACT_CHAT_ACK", "PROTO_NC_ACT_SOMEONECHAT_CMD"),
        ("PROTO_NC_ACT_GM_NOTICE_ACK", "PROTO_NC_ACT_NOTICE_CMD"),
        ("PROTO_NC_ACT_STOP_ACK", "PROTO_NC_ACT_SOMEONESTOP_CMD"),
        ("PROTO_NC_ACT_WALK_ACK", "PROTO_NC_ACT_SOMEONEMOVEWALK_CMD"),
        ("PROTO_NC_ACT_RUN_ACK", "PROTO_NC_ACT_SOMEONEMOVERUN_CMD"),
        ("PROTO_NC_ACT_TELEPORT_CMD", "PROTO_NC_ACT_MOVEFAIL_CMD"),
        ("PROTO_NC_ACT_EMOTE_REQ", "PROTO_NC_ACT_EMOTICON_CMD"),
        ("PROTO_NC_ACT_JUMP_ACK", "PROTO_NC_ACT_SOMEEONEJUMP_CMD"),
        ("PROTO_NC_ACT_REST_START_REQ", "PROTO_NC_ACT_PITCHTENT_REQ"),
        ("PROTO_NC_ACT_REST_START_ACK", "PROTO_NC_ACT_PITCHTENT_ACK"),
        ("PROTO_NC_ACT_REST_STOP_REQ", "PROTO_NC_ACT_FOLDTENT_REQ"),
        ("PROTO_NC_ACT_REST_STOP_ACK", "PROTO_NC_ACT_FOLDTENT_ACK"),
        ("PROTO_NC_ACT_BLOCK_WALK_ACK", "PROTO_NC_ACT_REINFORCE_STOP_CMD"),
        ("PROTO_NC_ITEM_PREMIUM_LIST_ACK", "PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK"),
        ("PROTO_NC_ITEM_REWARD_LIST_ACK", "PROTO_NC_ITEM_REWARDINVENOPEN_ACK"),
        ("PROTO_NC_FRIEND_LIST_ACK", "PROTO_NC_FRIEND_LIST_CMD"),
        ("PROTO_NC_CHAR_INGAME_REQ", "PROTO_NC_KQ_LIST_REFRESH_REQ"),
        ("PROTO_NC_CHAR_INGAME_ACK", "PROTO_NC_KQ_LIST_TIME_ACK"),
        ("PROTO_NC_KQ_LIST_COUNT_ACK", "PROTO_NC_KQ_LIST_ADD_ACK"),
        ("PROTO_NC_CHAR_GET_QUICKBAR_REQ", "PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_REQ"),
        ("PROTO_NC_CHAR_GET_QUICKBAR_ACK", "PROTO_NC_CHAR_OPTION_GET_SHORTCUTDATA_ACK"),
        ("PROTO_NC_CHAR_GET_QUICKBAR_STATE_REQ", "PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ"),
        ("PROTO_NC_CHAR_GET_QUICKBAR_STATE_ACK", "PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK"),
        ("PROTO_NC_CHAR_GET_GAME_SETTINGS_REQ", "PROTO_NC_CHAR_OPTION_GET_GAME_REQ"),
        ("PROTO_NC_CHAR_GET_GAME_SETTINGS_ACK", "PROTO_NC_CHAR_OPTION_GET_GAME_ACK"),
        ("PROTO_NC_CHAR_GET_CLIENT_SETTINGS_REQ", "PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_REQ"),
        ("PROTO_NC_CHAR_GET_CLIENT_SETTINGS_ACK", "PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK"),
        ("PROTO_NC_CHAR_GET_SHORTCUTS_REQ", "PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_REQ"),
        ("PROTO_NC_CHAR_GET_SHORTCUTS_ACK", "PROTO_NC_CHAR_OPTION_GET_KEYMAPPING_ACK"),
        ("PROTO_NC_CHAR_SAVE_QUICKBAR_REQ", "PROTO_NC_CHAR_OPTION_SET_SHORTCUTDATA_CMD"),
        ("PROTO_NC_CHAR_SAVE_QUICKBAR_STATE_REQ", "PROTO_NC_CHAR_OPTION_SET_SHORTCUTSIZE_CMD"),
        ("PROTO_NC_CHAR_SAVE_GAME_SETTINGS_REQ", "PROTO_NC_CHAR_OPTION_SET_GAME_CMD"),
        ("PROTO_NC_CHAR_SAVE_CLIENT_SETTINGS_REQ", "PROTO_NC_CHAR_OPTION_SET_WINDOWPOS_CMD"),
        ("PROTO_NC_CHAR_SAVE_SHORTCUTS_REQ", "PROTO_NC_CHAR_OPTION_SET_KEYMAPPING_CMD"),
        ("PROTO_NC_SYSTEM_UNKNOWN_REQ", "PROTO_NC_PRISON_GET_REQ"),
        ("PROTO_NC_SYSTEM_UNKNOWN_ACK", "PROTO_NC_PRISON_GET_ACK"),
        ("PROTO_NC_CHAR_BASE_CMD", "PROTO_NC_CHAR_CLIENT_BASE_CMD"),
        ("PROTO_NC_CHAR_LOOK_CMD", "PROTO_NC_CHAR_CLIENT_SHAPE_CMD"),
        ("PROTO_NC_QUEST_BUSY_LIST_CMD", "PROTO_NC_CHAR_CLIENT_QUEST_DOING_CMD"),
        ("PROTO_NC_QUEST_DONE_LIST_CMD", "PROTO_NC_CHAR_CLIENT_QUEST_DONE_CMD"),
        ("PROTO_NC_SKILL_ACTIVE_LIST_CMD", "PROTO_NC_CHAR_CLIENT_SKILL_CMD"),
        ("PROTO_NC_SKILL_PASSIVE_LIST_CMD", "PROTO_NC_CHAR_CLIENT_PASSIVE_CMD"),
        ("PROTO_NC_ITEM_LIST_CMD", "PROTO_NC_CHAR_CLIENT_ITEM_CMD"),
        ("PROTO_NC_CHAR_INFO_END_CMD", "PROTO_NC_CHAR_CLIENT_GAME_CMD"),
        ("PROTO_NC_TITLE_LIST_CMD", "PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD"),
    ]

    lines.append("// ========================================================================")
    lines.append("// Compatibilidad de Nombres Legacy / Transicionales (Deprecated Aliases)")
    lines.append("// ========================================================================")
    lines.append("const (")

    for alias, target in raw_aliases:
        if alias not in seen_names and alias != target:
            lines.append(f"\t{alias:<42} = {target}")

    lines.append(")")
    lines.append("")
    
    output_path = 'pkg/protocol/opcodes.go'
    with open(output_path, 'w', encoding='utf-8') as f:
        f.write("\n".join(lines))
        
    print(f"Successfully generated {total_ops} official opcodes in {output_path}!")

if __name__ == '__main__':
    main()
