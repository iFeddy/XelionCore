import re

type_map = {
    'int': 'INT NOT NULL DEFAULT 0',
    'smallint': 'SMALLINT NOT NULL DEFAULT 0',
    'tinyint': 'TINYINT UNSIGNED NOT NULL DEFAULT 0',
    'bigint': 'BIGINT NOT NULL DEFAULT 0',
    'bit': 'TINYINT(1) NOT NULL DEFAULT 0',
    'datetime': 'DATETIME DEFAULT CURRENT_TIMESTAMP',
}

AUTO_INC = {
    "tCharacterChangeID": "nRowNo",
    "tCharacterChargedBuffer": "nNo",
    "tCharacterTitle": "nNo",
    "tGMItemRewardLog": "nLogNo",
    "tGambleCoinExchangeLog": "nLogNo",
    "tGambleJackpotLog": "nLogNo",
    "tGambleSlotMachineLog": "nGameNo",
    "tGuild": "nNo",
    "tGuildAcademyRanking": "nRowNo",
    "tLog_GuildMasterChange": "nRowNo",
    "tMinihouseGuest": "nArticleNo",
    "tParty": "nPartyNo",
    "tPetNew": "nPetNo",
    "tPrison": "nPrisonNum",
    "tStatisticsGameMoney": "nRowNo",
}

TABLE_KEYS = {
    "tArenaChallenge": ("PRIMARY KEY (`nCharNo`, `nMonsterLevel`)", []),
    "tArenaFBZ": ("PRIMARY KEY (`nCharNo`)", []),
    "tArenaGuild": ("PRIMARY KEY (`nGuildNo`)", []),
    "tArenaPvP": ("PRIMARY KEY (`nCharNo`)", []),
    "tAttendance": ("PRIMARY KEY (`nCharNo`, `dAttendanceDate`)", []),
    "tAttendanceReward": ("PRIMARY KEY (`nCharNo`, `dEventDate`, `nRewardID`)", []),
    "tAutoMon": ("PRIMARY KEY (`nCharNo`, `nType`, `nSlot`)", []),
    "tAutoPick": ("PRIMARY KEY (`nCharNo`)", []),
    "tCharacterChangeID": ("PRIMARY KEY (`nRowNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tCharacterChargedBuffer": ("PRIMARY KEY (`nNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tCharacterGameOption": ("PRIMARY KEY (`nCharNo`, `nOptionNo`)", []),
    "tCharacterKeyMap": ("PRIMARY KEY (`nCharNo`, `nFunctionNo`)", []),
    "tCharacterLoading": (None, ["KEY `idx_char` (`nCharNo`)"]),
    "tCharacterSingleOption": ("PRIMARY KEY (`nCharNo`)", []),
    "tCharacterState": ("PRIMARY KEY (`nCharNo`, `nStateNo`)", []),
    "tCharacterTitle": ("PRIMARY KEY (`nNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tChatBlock": ("PRIMARY KEY (`nUserNo`, `nCharNo`)", []),
    "tChatColor": ("PRIMARY KEY (`nCharNo`)", []),
    "tCollectionCard": ("PRIMARY KEY (`nCharNo`, `nCardID`, `nSerialNo`)", ["KEY `idx_card` (`nCharNo`, `nCardID`)"]),
    "tCollectionFavorites": ("PRIMARY KEY (`nCharNo`, `nSlotNo`)", []),
    "tCollectionReward": ("PRIMARY KEY (`nCharNo`, `nReward`)", []),
    "tCollectionSerialCount": ("PRIMARY KEY (`nCardID`)", []),
    "tDice": ("PRIMARY KEY (`nGameNo`, `nCharNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tDiceLog": (None, ["KEY `idx_game` (`nGameNo`)", "KEY `idx_char` (`nCustomCharNo`)"]),
    "tDicePayLog": (None, ["KEY `idx_char` (`nCharNo`)", "KEY `idx_item` (`nItemKey`)"]),
    "tDiceResult": (None, ["KEY `idx_game` (`nGameNo`)", "KEY `idx_char` (`nCustomCharNo`)"]),
    "tExp": ("PRIMARY KEY (`level`)", []),
    "tFriend": ("PRIMARY KEY (`nCharNo`, `nFriendCharNo`)", ["KEY `idx_friend` (`nFriendCharNo`)"]),
    "tGMItemReward": ("PRIMARY KEY (`nGroupNo`)", []),
    "tGMItemRewardLog": ("PRIMARY KEY (`nLogNo`)", ["KEY `idx_group` (`nGroupNo`)", "KEY `idx_char` (`nCharNo`)"]),
    "tGambleCoin": ("PRIMARY KEY (`nCharNo`)", []),
    "tGambleCoinExchangeLog": ("PRIMARY KEY (`nLogNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tGambleDice": ("PRIMARY KEY (`nGameNo`)", []),
    "tGambleDiceLog": (None, ["KEY `idx_game` (`nGameNo`)", "KEY `idx_char` (`nCharNo`)"]),
    "tGambleDiceResult": (None, ["KEY `idx_game` (`nGameNo`)", "KEY `idx_char` (`nCharNo`)"]),
    "tGambleJackpot": ("PRIMARY KEY (`sGameName`, `nBattingUnit`)", []),
    "tGambleJackpotLog": ("PRIMARY KEY (`nLogNo`)", []),
    "tGambleSlotMachineLog": ("PRIMARY KEY (`nGameNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tGameOptionInit": ("PRIMARY KEY (`nOptionNo`)", []),
    "tGuild": ("PRIMARY KEY (`nNo`)", ["UNIQUE KEY `idx_name` (`sName`)"]),
    "tGuildAcademy": ("PRIMARY KEY (`nNo`)", ["KEY `idx_master` (`nMasterCharNo`)"]),
    "tGuildAcademyGraduate": ("PRIMARY KEY (`nNo`, `nCharNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tGuildAcademyHistory": (None, ["KEY `idx_academy` (`nAcademyID`)", "KEY `idx_char` (`nCharNo`)"]),
    "tGuildAcademyMember": ("PRIMARY KEY (`nNo`, `nCharNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tGuildAcademyRanking": ("PRIMARY KEY (`nRowNo`)", ["KEY `idx_ranking` (`dDate`, `nRanking`)"]),
    "tGuildAcademyReward": ("PRIMARY KEY (`nNo`, `nLevel`, `nTypeID`)", []),
    "tGuildAcademyRewardType": ("PRIMARY KEY (`nTypeID`)", []),
    "tGuildEmblem": ("PRIMARY KEY (`nGuildNo`)", []),
    "tGuildMember": ("PRIMARY KEY (`nNo`, `nCharNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tGuildToken": ("PRIMARY KEY (`nNo`)", []),
    "tGuildTournament": ("PRIMARY KEY (`nGTNo`, `nGuildNo`)", ["KEY `idx_guild` (`nGuildNo`)"]),
    "tGuildTournamentRegist": ("PRIMARY KEY (`nGTNo`, `nNo`)", ["KEY `idx_guild` (`nNo`)"]),
    "tHolyPromise": ("PRIMARY KEY (`nCharNo`, `nUpCharNo`)", ["KEY `idx_upchar` (`nUpCharNo`)"]),
    "tHolyPromiseReward": ("PRIMARY KEY (`nCharNo`, `nSeqNo`)", []),
    "tInGameGuide": ("PRIMARY KEY (`nCharNo`, `nGuideNo`)", []),
    "tItemActionCoolTime": ("PRIMARY KEY (`nCharNo`, `nActionGroup`)", []),
    "tItemEnchant": ("PRIMARY KEY (`nItemKey`, `nSlot`)", ["KEY `idx_item` (`nItemKey`)"]),
    "tItemMobList": ("PRIMARY KEY (`nItemKey`, `nMobID`)", ["KEY `idx_item` (`nItemKey`)"]),
    "tItemOptions": ("PRIMARY KEY (`nItemKey`, `nOptionType`)", ["KEY `idx_item` (`nItemKey`)"]),
    "tItemStorageFindEmpty": (None, ["KEY `idx_storage` (`nStorage`)"]),
    "tItemTitle": ("PRIMARY KEY (`nItemKey`, `nTitleType`)", ["KEY `idx_item` (`nItemKey`)"]),
    "tKeyMapInit": ("PRIMARY KEY (`nFunctionNo`)", []),
    "tKingdomQuest": ("PRIMARY KEY (`nQuestNo`)", []),
    "tLog_GuildMasterChange": ("PRIMARY KEY (`nRowNo`)", ["KEY `idx_guild` (`nGuildNo`)"]),
    "tMarriage": ("PRIMARY KEY (`nCharNo`)", ["KEY `idx_partner` (`nPartnerCharNo`)"]),
    "tMasterPupils": ("PRIMARY KEY (`nPupils`, `nMaster`)", ["KEY `idx_master` (`nMaster`)"]),
    "tMinihouse": ("PRIMARY KEY (`nCharNo`)", []),
    "tMinihouseGate": ("PRIMARY KEY (`nCharNo`, `nItemID`)", []),
    "tMinihouseGuest": ("PRIMARY KEY (`nArticleNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tParty": ("PRIMARY KEY (`nPartyNo`)", ["KEY `idx_master` (`nMasterCharNo`)"]),
    "tPartyMember": ("PRIMARY KEY (`nPartyNo`, `nCharNo`)", ["KEY `idx_char` (`nCharNo`)"]),
    "tPet": ("PRIMARY KEY (`nStorage`, `nPetNumber`, `nOwner`)", ["KEY `idx_owner` (`nOwner`)", "KEY `idx_petkey` (`nPetKey`)"]),
    "tPetNew": ("PRIMARY KEY (`nPetNo`)", ["KEY `idx_owner` (`nOwnerNo`)"]),
    "tPetSkill": ("PRIMARY KEY (`nPetKey`, `nSkillNo`, `nSkillLevel`)", []),
    "tPrison": ("PRIMARY KEY (`nPrisonNum`)", ["KEY `idx_char` (`nCharID`)"]),
    "tQuest": ("PRIMARY KEY (`nCharNo`, `nQuestNo`)", []),
    "tQuestTimes": ("PRIMARY KEY (`nCharNo`, `nQuestNo`)", []),
    "tRestExp": ("PRIMARY KEY (`nCharNo`)", []),
    "tSkill": ("PRIMARY KEY (`nCharNo`, `nSkillNo`)", []),
    "tSpamBlockTime": ("PRIMARY KEY (`nCharNo`)", []),
    "tSpamFilter": ("PRIMARY KEY (`nUserNo`)", []),
    "tStatCharLevel": (None, ["KEY `idx_stat` (`dDate`, `nLevel`)"]),
    "tStatGameMoney": (None, ["KEY `idx_date` (`dDate`)"]),
    "tStatisticsGameMoney": ("PRIMARY KEY (`nRowNo`)", ["KEY `idx_date` (`dDate`)"]),
    "tTutorial": ("PRIMARY KEY (`nUserNo`, `nCharNo`)", []),
    "tUserMoney": ("PRIMARY KEY (`nUserNo`)", []),
}

def convert_type(t, is_auto_inc=False):
    t = t.strip()
    if is_auto_inc:
        if t == 'bigint':
            return 'BIGINT NOT NULL AUTO_INCREMENT'
        return 'INT NOT NULL AUTO_INCREMENT'
    if t.startswith('nvarchar(max)') or t.startswith('varchar(max)'):
        return 'TEXT DEFAULT NULL'
    m = re.match(r'(?:nvarchar|varchar|char)\((\d+)\)', t)
    if m:
        return f"VARCHAR({m.group(1)}) NOT NULL DEFAULT ''"
    m = re.match(r'varbinary\((\d+)\)', t)
    if m:
        return f"VARBINARY({m.group(1)}) DEFAULT NULL"
    return type_map.get(t, "VARCHAR(255) NOT NULL DEFAULT ''")

with open('docs/TablasLegacy.md', 'r') as f:
    content = f.read()

matches = re.findall(r'\|\s*`World00_Character\.bak`\s*\|\s*\*\*`?([^`*]+)`?\*\*\s*\|\s*([^|]+)\|', content)

skip_prefixes = ('chr_', 'TEMP_', 'Temp_', 'temp_', 'tTEMP_', 'deleted_', 'audit', 'Waits_', 'removeskill', 'skilldata', 'QuestItemFix', 'num')
skip_suffixes = ('_Backup', '_BackUp', '_deleted', '_20091204', '_100610')

existing_tables = {
    'tCharacter', 'tCharacterShape', 'tCharacterOptions', 'tCharacterShortCut', 'tItem'
}

ddl_statements = []

for tname, col_str in sorted(matches, key=lambda x: x[0]):
    if tname in existing_tables:
        continue
    if any(tname.startswith(p) for p in skip_prefixes):
        continue
    if any(tname.endswith(s) for s in skip_suffixes):
        continue
    
    col_matches = re.findall(r'`?([a-zA-Z0-9_]+)`?\s*\(([^)]+)\)', col_str)
    if not col_matches:
        continue
    
    auto_col = AUTO_INC.get(tname)
    
    col_defs = []
    seen_cols = set()
    for cname, ctype in col_matches:
        if cname in seen_cols:
            continue
        seen_cols.add(cname)
        is_ai = (cname == auto_col)
        col_defs.append(f"    `{cname}` {convert_type(ctype, is_ai)}")
    
    pk_info, extra_keys = TABLE_KEYS.get(tname, (None, []))
    
    constraints = []
    if pk_info:
        constraints.append(f"    {pk_info}")
    for ek in extra_keys:
        constraints.append(f"    {ek}")
        
    all_lines = col_defs + constraints
    
    sql = f"-- ----------------------------------------------------------------------------\n"
    sql += f"-- Table: {tname}\n"
    sql += f"-- ----------------------------------------------------------------------------\n"
    sql += f"CREATE TABLE IF NOT EXISTS `{tname}` (\n"
    sql += ",\n".join(all_lines)
    sql += "\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;\n"
    ddl_statements.append((tname, sql))

print(f"Generated DDL for {len(ddl_statements)} official tables.")

with open('database/zone_tables.sql', 'w') as out:
    out.write("-- ============================================================================\n")
    out.write("-- XELION ONLINE - TABLAS OFICIALES DE WORLD00_CHARACTER (ZONE & CHARACTER)\n")
    out.write("-- Generado con 100% de fidelidad respecto a docs/TablasLegacy.md\n")
    out.write("-- ============================================================================\n\n")
    out.write("USE `xelion_online`;\n\n")
    for tname, sql in ddl_statements:
        out.write(sql + "\n")
