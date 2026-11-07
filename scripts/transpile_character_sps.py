#!/usr/bin/env python3
"""
transpile_character_sps.py - Transpilador 100% exhaustivo y pool-safe de los 425 Stored Procedures
de World00_Character (311 ODBC de Character.exe U 405 no-DBA de World00_Character.bak = 419 oficiales + 6 S2S)
a MySQL 8.4 usando encapsulación proc_logic + SELECT final incondicional y guarda de llamadas anidadas.
"""

import os
import re
import shutil
import mysql.connector

TARGET_DIR = "docs/sql/procedures/World00_Character"
DBA_EXCLUDED = {"adam_listproc", "commandexecute", "indexoptimize"}

# 1. Extraer los 311 Stored Procedures ODBC exactos invocados por Character.exe (`{CALL ...}`)
CHAR_EXE_PATH = "../AstraServer/Server/Character/Character.exe"
odbc_311 = []
if os.path.exists(CHAR_EXE_PATH):
    with open(CHAR_EXE_PATH, "rb") as f:
        exe_data = f.read()
    call_matches = re.findall(rb"\{\s*CALL\s+([a-zA-Z0-9_]+)", exe_data, flags=re.IGNORECASE)
    seen_odbc = {}
    for cm in call_matches:
        nm = cm.decode("ascii", errors="ignore")
        seen_odbc.setdefault(nm.lower(), nm)
    odbc_311 = sorted(seen_odbc.values(), key=lambda s: s.lower())

# 2. Mapear los 405 archivos .sql no-DBA de World00_Character.bak (case-insensitive)
bak_files_by_lower = {}
for fname in sorted(os.listdir(TARGET_DIR)):
    if not fname.endswith(".sql"):
        continue
    base = fname[:-4]
    if base.lower() in DBA_EXCLUDED:
        continue
    bak_files_by_lower[base.lower()] = base

# 3. Construir la unión exacta de 419 SPs oficiales + 6 SPs operativos S2S = 425 SPs
canonical_by_lower = {}
for bname in bak_files_by_lower.values():
    canonical_by_lower[bname.lower()] = bname
for oname in odbc_311:
    canonical_by_lower[oname.lower()] = oname

ODBC_ONLY_14 = [
    "p_Guild_SetMoney",
    "p_IsPromotionUser",
    "p_Item_MiniHouse_Effect_Resource_Del",
    "p_MiniHouseURL_Get",
    "p_MiniHouseURL_Set",
    "p_Promotion_Get",
    "p_Promotion_Set",
    "usp_CharacterLoadingDetailInfo",
    "usp_GuildTournament_setRefundNotice",
    "usp_Guild_getHistoryList",
    "usp_Guild_setHistory",
    "usp_LinkScroll_SetDestMap",
    "usp_Pet_setDelete",
    "usp_Pet_setTendency",
]
for oname in ODBC_ONLY_14:
    canonical_by_lower[oname.lower()] = oname

OPERATIONAL_S2S_6 = [
    "usp_Session_IssueZoneTransfer",
    "usp_Session_ConsumeZoneTransfer",
    "usp_WorldZone_UpsertNode",
    "usp_WorldZone_Heartbeat",
    "usp_WorldZone_GetList",
    "usp_WorldZone_ResolveMapEndpoint",
]
for sname in OPERATIONAL_S2S_6:
    canonical_by_lower[sname.lower()] = sname

TARGET_425 = sorted(canonical_by_lower.values(), key=lambda s: s.lower())


def wrap_nested_call(call_stmt):
    return (
        "SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;\n"
        f"        {call_stmt}\n"
        "        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);"
    )


def map_type(t):
    t = t.strip()
    tl = t.lower()
    if tl == "bit":
        return "TINYINT"
    if tl.startswith("varchar(max)") or tl.startswith("nvarchar(max)"):
        return "TEXT"
    if tl.startswith("nvarchar"):
        return re.sub(r"^nvarchar", "VARCHAR", t, flags=re.IGNORECASE)
    if tl.startswith("nchar"):
        return re.sub(r"^nchar", "CHAR", t, flags=re.IGNORECASE)
    if tl.startswith("varbinary(max)"):
        return "BLOB"
    if tl == "money":
        return "BIGINT"
    if tl == "smallmoney":
        return "INT"
    return t.upper()


def clean_sql_source(raw):
    raw = raw.replace("\r\n", "\n").replace("\r", "\n")
    parts = re.split(r"^\s*GO\s*$", raw, flags=re.IGNORECASE | re.MULTILINE)
    source = parts[0]

    lines = source.split("\n")
    clean_lines = []
    for l in lines:
        s = l.strip()
        if not s:
            clean_lines.append("")
            continue
        if s.isdigit() or re.match(r"^\d+\s*$", s) or s.startswith("<"):
            continue
        if any(ord(c) < 32 and c not in "\t\n" for c in s):
            continue
        if re.match(r"^0[\x00-\x20\x7f-\uffff]", s) or re.match(r"^0\s+[^\w\s]", s) or s.upper() == "GO":
            continue
        clean_lines.append(l)
    return "\n".join(clean_lines)


# ----------------------------------------------------------------------------
# Procedimientos Especializados, los 14 Exclusivos de Character.exe y los 6 S2S
# (100% Pool-Safe: proc_logic: BEGIN ... END proc_logic; + SELECT incondicional)
# ----------------------------------------------------------------------------
SPECIAL_PROCEDURES_RAW = {
    "p_Char_Create": """CREATE PROCEDURE `p_Char_Create`(
    IN in_nUserNo INT,
    IN in_nCreateWorld TINYINT,
    IN in_nAdminLevel TINYINT,
    IN in_nSlotNo TINYINT,
    IN in_sID VARCHAR(40),
    IN in_nRace TINYINT,
    IN in_nClass TINYINT,
    IN in_nGender TINYINT,
    IN in_nHairType TINYINT,
    IN in_nHairColor TINYINT,
    IN in_nFaceShape TINYINT,
    OUT out_nCharNo INT
)
sp_body: BEGIN
    DECLARE v_newCharNo INT DEFAULT 0;
    DECLARE v_authID TINYINT DEFAULT 0;
    DECLARE v_adminLevel TINYINT DEFAULT 0;

    SET out_nCharNo = 0;
    SET v_adminLevel = in_nAdminLevel;

    proc_logic: BEGIN
        SELECT IFNULL(nAuthID, 0) INTO v_authID
        FROM tUser
        WHERE nUserNo = in_nUserNo
        LIMIT 1;

        IF v_authID = 9 THEN
            SET v_adminLevel = 100;
        END IF;

        IF EXISTS (SELECT 1 FROM tCharacter WHERE nUserNo = in_nUserNo AND nSlotNo = in_nSlotNo AND bDeleted = 0) THEN
            SET out_nCharNo = -1;
            LEAVE proc_logic;
        END IF;

        IF EXISTS (SELECT 1 FROM tCharacter WHERE sID = in_sID) THEN
            SET out_nCharNo = -2;
            LEAVE proc_logic;
        END IF;

        START TRANSACTION;

        INSERT INTO tCharacter (nUserNo, nSlotNo, sID, nCreateWorld, nAdminLevel, dCreateDate)
        VALUES (in_nUserNo, in_nSlotNo, in_sID, in_nCreateWorld, v_adminLevel, NOW());

        SET v_newCharNo = LAST_INSERT_ID();

        INSERT INTO tCharacterShape (nCharNo, nRace, nClass, nGender, nHairType, nHairColor, nFaceShape)
        VALUES (v_newCharNo, in_nRace, in_nClass, in_nGender, in_nHairType, in_nHairColor, in_nFaceShape);

        INSERT INTO tCharacterOptions (nCharNo)
        VALUES (v_newCharNo);

        INSERT IGNORE INTO tCharacterKeyMap (nCharNo, nFunctionNo, sys_create_date, sys_update_date, nExtendedKey, nASCIICode)
        SELECT v_newCharNo, nFunctionNo, NOW(), NOW(), nExtendedKey, nASCIICode
        FROM tKeyMapInit;

        COMMIT;
        SET out_nCharNo = v_newCharNo;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCharNo AS nCharNo;
    END IF;
END""",

    "p_Char_Delete": """CREATE PROCEDURE `p_Char_Delete`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_sID VARCHAR(40) DEFAULT '';
    DECLARE v_DeleteMark VARCHAR(2) DEFAULT '~';
    DECLARE v_StartDeleteMark INT DEFAULT 0;
    DECLARE v_CurrentDate DATETIME DEFAULT NOW();
    DECLARE v_sChangedID VARCHAR(40) DEFAULT '';

    SET out_nRet = 0;

    proc_logic: BEGIN
        IF EXISTS(SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted > 0) THEN
            SET out_nRet = 1;
            LEAVE proc_logic;
        END IF;

        SELECT sID INTO v_sID FROM tCharacter WHERE nCharNo = in_nCharNo LIMIT 1;
        IF v_sID IS NULL OR CHAR_LENGTH(v_sID) = 0 THEN
            SET out_nRet = 2;
            LEAVE proc_logic;
        END IF;

        SET v_StartDeleteMark = LOCATE(v_DeleteMark, v_sID);
        IF v_StartDeleteMark > 0 THEN
            SET v_sID = LEFT(v_sID, v_StartDeleteMark - 1);
        END IF;

        SET v_sChangedID = CONCAT(LEFT(RTRIM(v_sID), 28), '~', DATE_FORMAT(v_CurrentDate, '%Y%m%d'), '#0');

        UPDATE tCharacter
        SET bDeleted = 1,
            dDeletedDate = NOW(),
            sID = v_sChangedID
        WHERE nCharNo = in_nCharNo;

        IF ROW_COUNT() = 1 THEN
            SET out_nRet = 0;
        ELSE
            SET out_nRet = 3;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_Char_GetListOfUserChar": """CREATE PROCEDURE `p_Char_GetListOfUserChar`(
    IN in_nUserNo INT
)
sp_body: BEGIN
    SELECT nCharNo FROM tCharacter WHERE nUserNo = in_nUserNo AND bDeleted = 0;
END""",

    "p_Char_GetLoginData": """CREATE PROCEDURE `p_Char_GetLoginData`(
    IN in_nCharNo INT
)
sp_body: BEGIN
    SELECT
        C.sID, C.nLevel, C.nSlotNo, C.sLoginZone, C.bDeleted, C.dDeletedDate,
        S.nRace, S.nClass, S.nGender, S.nHairType, S.nHairColor, S.nFaceShape,
        C.nKQHandle, C.sKQMap, C.nKQX, C.nKQY, C.nKQD, C.dKQDate
    FROM tCharacter AS C
    INNER JOIN tCharacterShape AS S ON S.nCharNo = in_nCharNo
    WHERE C.nCharNo = in_nCharNo;
END""",

    "p_Char_GetAllData": """CREATE PROCEDURE `p_Char_GetAllData`(
    IN in_nCharNo INT
)
sp_body: BEGIN
    SELECT
        nUserNo, nSlotNo, sID, nLevel, nExp, nHP, nSP, nAP, nHPS, nSPS, nPS, nGS,
        nMoney, nFame, nFlags, nAdminLevel, sLoginZone, nLoginZoneX, nLoginZoneY, nLoginZoneD,
        nPrisonMin, nPKCount, nPKYellowTime, nStrength, nConstitute, nDexterity, nIntelligence,
        nMentalPower, nRedistributePoint, nPlayMin, nLoginCount, dLastLoginDate, dCreateDate,
        nCreateWorld, bDeleted, dDeletedDate
    FROM tCharacter
    WHERE nCharNo = in_nCharNo;
END""",

    "p_Item_Create": """CREATE PROCEDURE `p_Item_Create`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    IN in_nItemID INT,
    IN in_nFlags INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        INSERT INTO tItem (nItemKey, nOwner, nStorageType, nStorage, nItemID, nFlags, dDate)
        VALUES (in_nItemKey, in_nOwner, in_nStorageType, in_nStorage, in_nItemID, in_nFlags, NOW())
        ON DUPLICATE KEY UPDATE nItemID = VALUES(nItemID), nFlags = VALUES(nFlags);
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_Item_GetList": """CREATE PROCEDURE `p_Item_GetList`(
    IN in_nCharNo INT
)
sp_body: BEGIN
    SELECT nItemKey, nStorageType, nStorage, nItemID, nFlags, dDate
    FROM tItem
    WHERE nOwner = in_nCharNo;
END""",

    "p_Item_Delete": """CREATE PROCEDURE `p_Item_Delete`(
    IN in_nItemKey BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        INSERT INTO tItem_deleted (nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags, dDate)
        SELECT nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags, NOW()
        FROM tItem
        WHERE nItemKey = in_nItemKey;

        DELETE FROM tItem WHERE nItemKey = in_nItemKey;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_Item_GetStorageEmptyList": """CREATE PROCEDURE `p_Item_GetStorageEmptyList`(
    IN in_nOwner INT,
    IN in_nStorageType INT,
    IN in_nGetStorageCnt INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF in_nOwner = -1 THEN
            SET out_nRet = -1000;
            LEAVE proc_logic;
        END IF;

        SELECT A.nStorage
        FROM tItemStorageFindEmpty A
        LEFT OUTER JOIN (
            SELECT nStorage FROM tItem WHERE nOwner = in_nOwner AND nStorageType = in_nStorageType
        ) B ON A.nStorage = B.nStorage
        WHERE B.nStorage IS NULL
        ORDER BY A.nStorage
        LIMIT in_nGetStorageCnt;
    END proc_logic;
END""",

    "p_Char_Find": """CREATE PROCEDURE `p_Char_Find`(
    IN in_sID VARCHAR(40),
    OUT out_nCharNo INT
)
sp_body: BEGIN
    SET out_nCharNo = 0;
    proc_logic: BEGIN
        SELECT nCharNo INTO out_nCharNo FROM tCharacter WHERE sID = in_sID LIMIT 1;
        IF out_nCharNo IS NULL THEN
            SET out_nCharNo = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCharNo AS nCharNo;
    END IF;
END""",

    "p_Char_GetCharNo": """CREATE PROCEDURE `p_Char_GetCharNo`(
    IN in_sID VARCHAR(40),
    OUT out_nCharNo INT
)
sp_body: BEGIN
    SET out_nCharNo = 0;
    proc_logic: BEGIN
        SELECT nCharNo INTO out_nCharNo FROM tCharacter WHERE sID = in_sID LIMIT 1;
        IF out_nCharNo IS NULL THEN
            SET out_nCharNo = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCharNo AS nCharNo;
    END IF;
END""",

    "p_Char_Login": """CREATE PROCEDURE `p_Char_Login`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_lastLogin DATETIME;
    DECLARE v_loginDiff INT DEFAULT 0;
    SET out_nRet = 0;

    proc_logic: BEGIN
        IF EXISTS (SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo) THEN
            SELECT dLastLoginDate INTO v_lastLogin FROM tCharacter WHERE nCharNo = in_nCharNo LIMIT 1;
            UPDATE tCharacter SET nLoginCount = nLoginCount + 1, dLastLoginDate = NOW() WHERE nCharNo = in_nCharNo;
            SET v_loginDiff = IFNULL(DATEDIFF(NOW(), v_lastLogin), 0);
            SET out_nRet = 0;
        ELSE
            SET out_nRet = -1111;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        IF out_nRet = 0 THEN
            SELECT v_loginDiff AS nLoginDateDiff;
        ELSE
            SELECT out_nRet AS nRet;
        END IF;
    END IF;
END""",

    "p_Char_Logout": """CREATE PROCEDURE `p_Char_Logout`(
    IN in_nCharNo INT,
    IN in_nPlayMin INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF EXISTS (SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo) THEN
            UPDATE tCharacter SET nPlayMin = nPlayMin + in_nPlayMin WHERE nCharNo = in_nCharNo;
            SET out_nRet = 1;
        ELSE
            SET out_nRet = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_View_GetData": """CREATE PROCEDURE `p_View_GetData`(
    IN in_QueryExpr VARCHAR(512)
)
sp_body: BEGIN
    SET @s = in_QueryExpr;
    PREPARE stmt FROM @s;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END""",

    "p_Quest_Set": """CREATE PROCEDURE `p_Quest_Set`(
    IN in_nCharNo INT,
    IN in_nQuestNo INT,
    IN in_nStatus TINYINT,
    IN in_sData VARBINARY(100),
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        UPDATE tQuest
        SET nStatus = in_nStatus, sData = in_sData
        WHERE nCharNo = in_nCharNo AND nQuestNo = in_nQuestNo;

        IF ROW_COUNT() = 0 THEN
            INSERT INTO tQuest (nCharNo, nQuestNo, nStatus, sData)
            VALUES (in_nCharNo, in_nQuestNo, in_nStatus, in_sData);
        END IF;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    # Los 3 SPs compartidos con Etapa 1 (Account) mantienen 100% compatibilidad con procedures_account.sql
    "usp_Tutorial_getCharInfo": """CREATE PROCEDURE `usp_Tutorial_getCharInfo`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        SELECT nState, nStep
        FROM tTutorial
        WHERE nCharNo = in_nCharNo
        LIMIT 1;
    END proc_logic;
END""",

    "usp_Tutorial_getIsSkip": """CREATE PROCEDURE `usp_Tutorial_getIsSkip`(
    IN in_nUserNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF EXISTS(SELECT 1 FROM tTutorial WHERE nUserNo = in_nUserNo AND nState >= 1) THEN
            SET out_nRet = 1;
        ELSE
            SET out_nRet = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_Tutorial_setCharInfo": """CREATE PROCEDURE `usp_Tutorial_setCharInfo`(
    IN in_nUserNo INT,
    IN in_nCharNo INT,
    IN in_nState TINYINT,
    IN in_nStep TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep)
        VALUES (in_nUserNo, in_nCharNo, in_nState, in_nStep)
        ON DUPLICATE KEY UPDATE nState = VALUES(nState), nStep = VALUES(nStep);
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    # Cross-DB wrappers de ChargeItem hacia tChargeItem en xelion_online
    "usp_ChargeItem_Draw": """CREATE PROCEDURE `usp_ChargeItem_Draw`(
    IN in_ROW_NO INT,
    IN in_USER_NO INT,
    OUT out_RET INT
)
sp_body: BEGIN
    SET out_RET = 0;
    proc_logic: BEGIN
        SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
        CALL `usp_Charge_ItemDraw`(in_ROW_NO, in_USER_NO, out_RET);
        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_RET AS RET;
    END IF;
END""",

    "usp_ChargeItem_GetList": """CREATE PROCEDURE `usp_ChargeItem_GetList`(
    IN in_USER_NO INT
)
sp_body: BEGIN
    CALL `usp_Charge_ItemSelectList`(in_USER_NO);
END""",

    "p_Event_Friend_Reward": """CREATE PROCEDURE `p_Event_Friend_Reward`(
    IN in_nFriendCharNo INT,
    IN in_nItemKey BIGINT,
    OUT out_err SMALLINT
)
sp_body: BEGIN
    DECLARE v_nFriendUserNo INT DEFAULT 0;
    DECLARE v_nLevel SMALLINT DEFAULT 0;
    DECLARE v_nRewardItemID INT DEFAULT 8004;
    DECLARE v_nStorage_Min SMALLINT DEFAULT 0;
    DECLARE v_nStorage_Max SMALLINT DEFAULT 143;
    DECLARE v_nStorageType TINYINT DEFAULT 6;
    DECLARE v_nItemID INT DEFAULT 0;
    DECLARE v_nOptionType SMALLINT DEFAULT 1;
    DECLARE v_nOptionData BIGINT DEFAULT 5;
    DECLARE v_nOldItemKey BIGINT DEFAULT 0;

    SET out_err = 0;

    proc_logic: BEGIN
        SELECT nLevel, nUserNo INTO v_nLevel, v_nFriendUserNo
        FROM tCharacter
        WHERE nCharNo = in_nFriendCharNo
        LIMIT 1;

        IF v_nFriendUserNo IS NULL OR v_nFriendUserNo = 0 THEN
            SET out_err = 1;
            LEAVE proc_logic;
        END IF;

        IF v_nLevel > 0 AND v_nLevel < 24 THEN
            SET v_nRewardItemID = 8004;
        ELSEIF v_nLevel >= 24 AND v_nLevel < 48 THEN
            SET v_nRewardItemID = 8013;
        ELSE
            SET v_nRewardItemID = 8022;
        END IF;

        SELECT IFNULL(MIN(nStorage), 0) INTO v_nStorage_Min
        FROM tItem
        WHERE nOwner = v_nFriendUserNo
          AND nStorageType = v_nStorageType
          AND nItemID IN (8004, 8013, 8022);

        IF v_nStorage_Min > 0 AND v_nStorage_Min < 37 THEN
            SET out_err = 2;
            LEAVE proc_logic;
        END IF;

        IF v_nStorage_Min = 0 THEN
            SET v_nStorage_Min = v_nStorage_Max;
            SET v_nOptionData = 5;
        ELSE
            SELECT I.nItemID, O.nOptionData, I.nItemKey
            INTO v_nItemID, v_nOptionData, v_nOldItemKey
            FROM tItem I
            INNER JOIN tItemOptions O ON I.nItemKey = O.nItemKey AND O.nOptionType = v_nOptionType
            WHERE I.nOwner = v_nFriendUserNo
              AND I.nStorageType = v_nStorageType
              AND I.nStorage = v_nStorage_Min
            LIMIT 1;

            IF v_nRewardItemID = v_nItemID AND v_nOptionData > 0 AND v_nOptionData < 6 THEN
                SET v_nOptionData = v_nOptionData + 5;
            ELSE
                SET v_nStorage_Min = v_nStorage_Min - 1;
                SET v_nOptionData = 5;
            END IF;
        END IF;

        IF v_nOptionData = 5 AND EXISTS (
            SELECT 1 FROM tItem
            WHERE nOwner = v_nFriendUserNo AND nStorageType = v_nStorageType AND nStorage = v_nStorage_Min
        ) THEN
            SET out_err = 2;
            LEAVE proc_logic;
        END IF;

        START TRANSACTION;
        IF v_nOptionData = 5 THEN
            INSERT INTO tItem (nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags, dDate)
            VALUES (in_nItemKey, v_nStorageType, v_nFriendUserNo, v_nStorage_Min, v_nRewardItemID, 0, NOW());

            INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
            VALUES (in_nItemKey, v_nOptionType, v_nOptionData);
        ELSE
            UPDATE tItemOptions
            SET nOptionData = v_nOptionData
            WHERE nItemKey = v_nOldItemKey AND nOptionType = v_nOptionType;
        END IF;
        COMMIT;
        SET out_err = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_err AS err;
    END IF;
END""",

    "p_Char_Wedding_Divorce_Cancel": """CREATE PROCEDURE `p_Char_Wedding_Divorce_Cancel`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    IN in_IncreseCount TINYINT,
    OUT out_SubStatusA TINYINT,
    OUT out_SubStatusB TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_SubStatusA = 0;
    SET out_SubStatusB = 0;
    SET out_nRet = 0;

    proc_logic: BEGIN
        START TRANSACTION;

        IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo) THEN
            SELECT nSubStatus INTO out_SubStatusA FROM tMarriage WHERE nCharNo = in_CharANo LIMIT 1;
            UPDATE tMarriage
            SET nSubStatus = CASE nSubStatus WHEN 1 THEN 5 WHEN 3 THEN 7 ELSE nSubStatus END,
                nDivorceCancel = nDivorceCancel + in_IncreseCount,
                dSubWhen = dWhen
            WHERE nCharNo = in_CharANo;
        ELSE
            ROLLBACK;
            SET out_nRet = 1;
            LEAVE proc_logic;
        END IF;

        IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo) THEN
            SELECT nSubStatus INTO out_SubStatusB FROM tMarriage WHERE nCharNo = in_CharBNo LIMIT 1;
            UPDATE tMarriage
            SET nSubStatus = 0,
                dSubWhen = dWhen
            WHERE nCharNo = in_CharBNo;
        ELSE
            ROLLBACK;
            SET out_nRet = 1;
            LEAVE proc_logic;
        END IF;

        COMMIT;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_SubStatusA AS SubStatusA, out_SubStatusB AS SubStatusB, out_nRet AS nRet;
    END IF;
END""",

    "p_Guild_Member_SetGrade": """CREATE PROCEDURE `p_Guild_Member_SetGrade`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nGrade TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_cnt INT DEFAULT 0;
    DECLARE v_guildMaster INT DEFAULT 0;
    DECLARE v_lastChangeDate DATETIME;

    SET out_nRet = 0;

    proc_logic: BEGIN
        IF IFNULL(in_nNo, 0) = 0 THEN
            SET out_nRet = -1311;
            LEAVE proc_logic;
        END IF;

        IF IFNULL(in_nCharNo, 0) = 0 THEN
            SET out_nRet = -1313;
            LEAVE proc_logic;
        END IF;

        SELECT IFNULL(COUNT(nCharNo), 0) INTO v_cnt
        FROM tCharacterChangeID
        WHERE nCharNo = in_nCharNo AND bFlag = 0 AND bInit = 1;

        IF in_nGrade = 0 AND v_cnt = 1 THEN
            SET out_nRet = -1323;
            LEAVE proc_logic;
        END IF;

        START TRANSACTION;

        IF in_nGrade = 0 THEN
            SELECT IFNULL(MAX(nCharNo), 0) INTO v_guildMaster
            FROM tGuildMember WHERE nNo = in_nNo AND nGrade = 0;

            SELECT MAX(dDate) INTO v_lastChangeDate
            FROM tLog_GuildMasterChange WHERE nGuildNo = in_nNo;

            IF v_lastChangeDate IS NOT NULL AND DATE_ADD(v_lastChangeDate, INTERVAL 7 DAY) > NOW() THEN
                ROLLBACK;
                SET out_nRet = -1324;
                LEAVE proc_logic;
            END IF;
        END IF;

        IF in_nGrade = 0 AND v_guildMaster <> 0 THEN
            UPDATE tGuildMember SET nGrade = 6 WHERE nNo = in_nNo AND nGrade = 0;
        END IF;

        UPDATE tGuildMember SET nGrade = in_nGrade WHERE nNo = in_nNo AND nCharNo = in_nCharNo;
        IF ROW_COUNT() <> 1 THEN
            ROLLBACK;
            SET out_nRet = -1325;
        ELSE
            IF in_nGrade = 0 THEN
                INSERT INTO tLog_GuildMasterChange (nGuildNo, nOldGuildMaster, nNewGuildMaster)
                VALUES (in_nNo, v_guildMaster, in_nCharNo);
            END IF;
            COMMIT;
            SET out_nRet = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_GuildAcademy_raiseMember": """CREATE PROCEDURE `usp_GuildAcademy_raiseMember`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nGrade TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_nMax INT DEFAULT 0;
    DECLARE v_nCnt INT DEFAULT 0;

    SET out_nRet = 0;

    proc_logic: BEGIN
        IF NOT EXISTS (SELECT nNo FROM tGuildAcademyMember WHERE nNo = in_nNo AND nCharNo = in_nCharNo) THEN
            SET out_nRet = -1528;
            LEAVE proc_logic;
        END IF;

        SELECT T1.nMaxMembers, T2.nCnt INTO v_nMax, v_nCnt
        FROM tGuild AS T1
        INNER JOIN (
            SELECT G.nNo, IFNULL(COUNT(*), 0) AS nCnt
            FROM tCharacter AS C, tGuildMember AS G
            WHERE C.nCharNo = G.nCharNo AND G.nNo = in_nNo AND C.bDeleted = 0
            GROUP BY G.nNo
        ) AS T2 ON (T1.nNo = T2.nNo);

        IF v_nCnt >= v_nMax THEN
            SET out_nRet = -1530;
            LEAVE proc_logic;
        END IF;

        START TRANSACTION;

        SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
        CALL `usp_GuildAcademy_delMember`(in_nNo, in_nCharNo, out_nRet);
        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
        IF out_nRet <> 0 THEN
            ROLLBACK;
            LEAVE proc_logic;
        END IF;

        SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
        CALL `p_Guild_Member_Join`(in_nNo, in_nCharNo, in_nGrade, out_nRet);
        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
        IF out_nRet <> 1 THEN
            SET out_nRet = -1531;
            ROLLBACK;
            LEAVE proc_logic;
        END IF;

        SET out_nRet = 0;
        COMMIT;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_Item_Upgrade_CelllInAll": """CREATE PROCEDURE `p_Item_Upgrade_CelllInAll`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nUpgradeType TINYINT,
    IN in_nUpgradeValue BIGINT,
    IN in_nRawKey BIGINT,
    IN in_nRawLeftKey BIGINT,
    IN in_nRawRightKey BIGINT,
    IN in_nRawMiddleKey BIGINT,
    IN in_nGiftMoney BIGINT,
    IN in_nRawOptionType SMALLINT,
    IN in_nRawLot SMALLINT,
    IN in_nRawLeftOptionType SMALLINT,
    IN in_nRawLeftLot SMALLINT,
    IN in_nRawRightOptionType SMALLINT,
    IN in_nRawRightLot SMALLINT,
    IN in_nRawMiddleOptionType SMALLINT,
    IN in_nRawMiddleLot SMALLINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_nTinyintRet TINYINT DEFAULT 0;
    DECLARE v_nRetTemp INT DEFAULT 0;

    SET out_nRet = 0;

    proc_logic: BEGIN
        SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
        CALL `p_Item_IsValidOwner`(in_nItemKey, in_nOwner, v_nTinyintRet);
        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
        IF v_nTinyintRet = 0 THEN
            SET out_nRet = 2;
            LEAVE proc_logic;
        END IF;

        START TRANSACTION;
        IF in_nGiftMoney > 0 THEN
            SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
            CALL `p_Char_MoneySub`(in_nOwner, in_nGiftMoney, v_nTinyintRet);
            SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            IF v_nTinyintRet = 0 THEN
                ROLLBACK;
                SET out_nRet = 3;
                LEAVE proc_logic;
            END IF;
        END IF;

        SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
        CALL `p_Item_SetOption`(in_nItemKey, in_nUpgradeType, in_nUpgradeValue, v_nRetTemp);
        SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
        IF v_nRetTemp <> 0 THEN
            ROLLBACK;
            SET out_nRet = 4;
            LEAVE proc_logic;
        END IF;

        IF in_nRawLot = 1 THEN
            SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
            CALL `p_Item_Delete`(in_nRawKey, v_nRetTemp);
            SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 5;
                LEAVE proc_logic;
            END IF;
        ELSE
            SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
            CALL `p_Item_SetOption`(in_nRawKey, in_nRawOptionType, in_nRawLot - 1, v_nRetTemp);
            SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 5;
                LEAVE proc_logic;
            END IF;
        END IF;

        IF in_nRawLeftKey <> 0 THEN
            IF in_nRawLeftLot = 1 THEN
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_Delete`(in_nRawLeftKey, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 6;
                    LEAVE proc_logic;
                END IF;
            ELSE
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_SetOption`(in_nRawLeftKey, in_nRawLeftOptionType, in_nRawLeftLot - 1, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 6;
                    LEAVE proc_logic;
                END IF;
            END IF;
        END IF;

        IF in_nRawRightKey <> 0 THEN
            IF in_nRawRightLot = 1 THEN
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_Delete`(in_nRawRightKey, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 7;
                    LEAVE proc_logic;
                END IF;
            ELSE
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_SetOption`(in_nRawRightKey, in_nRawRightOptionType, in_nRawRightLot - 1, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 7;
                    LEAVE proc_logic;
                END IF;
            END IF;
        END IF;

        IF in_nRawMiddleKey <> 0 THEN
            IF in_nRawMiddleLot = 1 THEN
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_Delete`(in_nRawMiddleKey, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 8;
                    LEAVE proc_logic;
                END IF;
            ELSE
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_SetOption`(in_nRawMiddleKey, in_nRawMiddleOptionType, in_nRawMiddleLot - 1, v_nRetTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
                IF v_nRetTemp <> 0 THEN
                    ROLLBACK;
                    SET out_nRet = 8;
                    LEAVE proc_logic;
                END IF;
            END IF;
        END IF;

        COMMIT;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    # ========================================================================
    # 14 Stored Procedures ODBC exclusivos de Character.exe (SQLBindParameter 1:1)
    # ========================================================================
    "p_Guild_SetMoney": """CREATE PROCEDURE `p_Guild_SetMoney`(
    IN in_nNo INT,
    IN in_nMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        UPDATE tGuild SET nMoney = in_nMoney WHERE nNo = in_nNo;
        IF ROW_COUNT() = 1 THEN
            SET out_nRet = 1;
        ELSE
            SET out_nRet = 0;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_IsPromotionUser": """CREATE PROCEDURE `p_IsPromotionUser`(
    IN in_nUserNo INT,
    OUT out_bIsPromotion TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_bIsPromotion = 0;
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF NOT EXISTS (SELECT 1 FROM tUser WHERE nUserNo = in_nUserNo AND bIsDelete = 0) THEN
            SET out_nRet = -1;
            LEAVE proc_logic;
        END IF;
        IF EXISTS (
            SELECT 1 FROM tCharacterTitle T
            INNER JOIN tCharacter C ON C.nCharNo = T.nCharNo
            WHERE C.nUserNo = in_nUserNo AND T.nType >= 250
        ) THEN
            SET out_bIsPromotion = 1;
        ELSE
            SET out_bIsPromotion = 0;
        END IF;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_bIsPromotion AS bIsPromotion, out_nRet AS nRet;
    END IF;
END""",

    "p_Item_MiniHouse_Effect_Resource_Del": """CREATE PROCEDURE `p_Item_MiniHouse_Effect_Resource_Del`(
    IN in_nOwner INT,
    IN in_nItemKey BIGINT,
    IN in_nDemandMoney BIGINT,
    IN in_nOptionType SMALLINT,
    IN in_nRemainLot INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_valid TINYINT DEFAULT 0;
    DECLARE v_retTemp INT DEFAULT 0;

    SET out_nRet = 0;

    proc_logic: BEGIN
        IF in_nItemKey <> 0 THEN
            SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
            CALL `p_Item_IsValidOwner`(in_nItemKey, in_nOwner, v_valid);
            SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            IF v_valid = 0 THEN
                SET out_nRet = 2;
                LEAVE proc_logic;
            END IF;
        END IF;

        START TRANSACTION;

        IF in_nDemandMoney > 0 THEN
            SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
            CALL `p_Char_MoneySub`(in_nOwner, in_nDemandMoney, v_valid);
            SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            IF v_valid = 0 THEN
                ROLLBACK;
                SET out_nRet = 3;
                LEAVE proc_logic;
            END IF;
        END IF;

        IF in_nItemKey <> 0 THEN
            IF in_nRemainLot <= 0 THEN
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_Delete`(in_nItemKey, v_retTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            ELSE
                SET @_xc_nested_sp = IFNULL(@_xc_nested_sp, 0) + 1;
                CALL `p_Item_SetOption`(in_nItemKey, in_nOptionType, in_nRemainLot, v_retTemp);
                SET @_xc_nested_sp = GREATEST(0, IFNULL(@_xc_nested_sp, 1) - 1);
            END IF;

            IF v_retTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 5;
                LEAVE proc_logic;
            END IF;
        END IF;

        COMMIT;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "p_MiniHouseURL_Get": """CREATE PROCEDURE `p_MiniHouseURL_Get`(
    IN in_nCharNo INT,
    OUT out_sURL VARCHAR(255)
)
sp_body: BEGIN
    SET out_sURL = '';
    proc_logic: BEGIN
        SELECT sContent INTO out_sURL
        FROM tMinihouseGuest
        WHERE nCharNo = in_nCharNo AND nWriterNo = 0 AND sWriterID = '__MINIHOUSE_URL__'
        ORDER BY nArticleNo DESC
        LIMIT 1;

        IF out_sURL IS NULL THEN
            SET out_sURL = '';
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_sURL AS sURL;
    END IF;
END""",

    "p_MiniHouseURL_Set": """CREATE PROCEDURE `p_MiniHouseURL_Set`(
    IN in_nCharNo INT,
    IN in_sURL VARCHAR(255)
)
sp_body: BEGIN
    INSERT IGNORE INTO tMinihouse (nCharNo, nVisitor, dDate)
    VALUES (in_nCharNo, 0, NOW());

    UPDATE tMinihouseGuest
    SET sContent = in_sURL, dDate = NOW()
    WHERE nCharNo = in_nCharNo AND nWriterNo = 0 AND sWriterID = '__MINIHOUSE_URL__';

    IF ROW_COUNT() = 0 THEN
        INSERT INTO tMinihouseGuest (nCharNo, nWriterNo, sWriterID, sContent, dDate)
        VALUES (in_nCharNo, 0, '__MINIHOUSE_URL__', in_sURL, NOW());
    END IF;
END""",

    "p_Promotion_Get": """CREATE PROCEDURE `p_Promotion_Get`(
    IN in_nCharNo INT,
    IN in_nPromotionType TINYINT,
    IN in_nStep TINYINT,
    OUT out_nCount INT,
    OUT out_dDate DATETIME
)
sp_body: BEGIN
    DECLARE v_type TINYINT UNSIGNED DEFAULT 0;
    SET out_nCount = 0;
    SET out_dDate = NOW();
    SET v_type = 250 + in_nPromotionType;

    proc_logic: BEGIN
        SELECT CAST(nElementValue AS SIGNED), dDate
        INTO out_nCount, out_dDate
        FROM tCharacterTitle
        WHERE nCharNo = in_nCharNo AND nType = v_type AND nElementNo = in_nStep
        ORDER BY nNo DESC
        LIMIT 1;

        IF out_nCount IS NULL THEN
            SET out_nCount = 0;
            SET out_dDate = NOW();
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCount AS nCount, out_dDate AS dDate;
    END IF;
END""",

    "p_Promotion_Set": """CREATE PROCEDURE `p_Promotion_Set`(
    IN in_nCharNo INT,
    IN in_nPromotionType TINYINT,
    IN in_nStep TINYINT,
    IN in_nCount INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_type TINYINT UNSIGNED DEFAULT 0;
    SET out_nRet = 0;
    SET v_type = 250 + in_nPromotionType;

    proc_logic: BEGIN
        UPDATE tCharacterTitle
        SET nElementValue = in_nCount, dDate = NOW()
        WHERE nCharNo = in_nCharNo AND nType = v_type AND nElementNo = in_nStep;

        IF ROW_COUNT() = 0 THEN
            INSERT INTO tCharacterTitle (dDate, nCharNo, nType, nElementNo, nElementValue)
            VALUES (NOW(), in_nCharNo, v_type, in_nStep, in_nCount);
        END IF;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_CharacterLoadingDetailInfo": """CREATE PROCEDURE `usp_CharacterLoadingDetailInfo`(
    IN in_nCharNo INT,
    IN in_nSeqNo INT,
    IN in_nFailType SMALLINT,
    IN in_nLoadingFailCnt INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        UPDATE tCharacterLoading
        SET nFailType = in_nFailType,
            nLoadingFailCnt = CAST(in_nLoadingFailCnt AS SIGNED),
            dDBTime = NOW()
        WHERE nCharNo = in_nCharNo AND nSeqNo = in_nSeqNo;

        IF ROW_COUNT() = 0 THEN
            INSERT INTO tCharacterLoading (nCharNo, nSeqNo, dWorldTime, dDBTime, sLoadingMapName, nFailType, nLoadingFailCnt)
            VALUES (in_nCharNo, in_nSeqNo, NOW(), NOW(), '', in_nFailType, CAST(in_nLoadingFailCnt AS SIGNED));
        END IF;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_GuildTournament_setRefundNotice": """CREATE PROCEDURE `usp_GuildTournament_setRefundNotice`(
    IN in_nGTNo INT,
    IN in_nNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF IFNULL(in_nGTNo, 0) <= 0 OR IFNULL(in_nNo, 0) <= 0 THEN
            SET out_nRet = -1411;
            LEAVE proc_logic;
        END IF;

        UPDATE tGuildTournamentRegist
        SET nMoney = 0
        WHERE nGTNo = in_nGTNo AND nNo = in_nNo AND bChoice = 0;

        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_Guild_getHistoryList": """CREATE PROCEDURE `usp_Guild_getHistoryList`(
    IN in_nGuildNo INT,
    IN in_nAmount INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_limit INT DEFAULT 50;
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF in_nAmount > 0 THEN
            SET v_limit = in_nAmount;
        END IF;

        SELECT sID, nAction, dDate
        FROM tGuildAcademyHistory
        WHERE nAcademyID = -in_nGuildNo
        ORDER BY dDate DESC
        LIMIT v_limit;
    END proc_logic;
END""",

    "usp_Guild_setHistory": """CREATE PROCEDURE `usp_Guild_setHistory`(
    IN in_nGuildNo INT,
    IN in_nCharNo INT,
    IN in_sID VARCHAR(40),
    IN in_nAction TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF IFNULL(in_nGuildNo, 0) <= 0 THEN
            SET out_nRet = -1;
            LEAVE proc_logic;
        END IF;

        INSERT INTO tGuildAcademyHistory (nAcademyID, nCharNo, sID, dDate, nAction)
        VALUES (-in_nGuildNo, in_nCharNo, in_sID, NOW(), in_nAction);
        SET out_nRet = 1;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_LinkScroll_SetDestMap": """CREATE PROCEDURE `usp_LinkScroll_SetDestMap`(
    IN in_nCharNo INT,
    IN in_sDestMap VARCHAR(16),
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        IF NOT EXISTS (SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo) THEN
            SET out_nRet = -1;
            LEAVE proc_logic;
        END IF;

        UPDATE tCharacter
        SET sKQMap = in_sDestMap
        WHERE nCharNo = in_nCharNo;
        SET out_nRet = 0;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_Pet_setDelete": """CREATE PROCEDURE `usp_Pet_setDelete`(
    IN in_nOwnerNo INT,
    IN in_nPetNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        DELETE FROM tPetNew
        WHERE nPetNo = in_nPetNo AND (in_nOwnerNo = 0 OR nOwnerNo = in_nOwnerNo);

        IF ROW_COUNT() = 1 THEN
            SET out_nRet = 1;
        ELSE
            SET out_nRet = -1;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    "usp_Pet_setTendency": """CREATE PROCEDURE `usp_Pet_setTendency`(
    IN in_nPetNo INT,
    IN in_nPetMind TINYINT,
    IN in_nPetStress TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    proc_logic: BEGIN
        UPDATE tPetNew
        SET nPetMind = in_nPetMind,
            nPetStress = in_nPetStress
        WHERE nPetNo = in_nPetNo;

        IF ROW_COUNT() = 1 THEN
            SET out_nRet = 1;
        ELSE
            SET out_nRet = -1;
        END IF;
    END proc_logic;

    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nRet AS nRet;
    END IF;
END""",

    # ========================================================================
    # 6 Stored Procedures Operativos S2S (WorldService <-> ZoneService)
    # ========================================================================
    "usp_Session_IssueZoneTransfer": """CREATE PROCEDURE `usp_Session_IssueZoneTransfer`(
    IN in_nWorldNo TINYINT UNSIGNED,
    IN in_nZoneNo SMALLINT UNSIGNED,
    IN in_nUserNo INT,
    IN in_sUserID VARCHAR(60),
    IN in_nCharNo INT,
    IN in_sCharID VARCHAR(40),
    IN in_nRandomID SMALLINT UNSIGNED,
    IN in_nAdminLevel TINYINT UNSIGNED,
    IN in_sClientIP VARCHAR(45),
    IN in_sTargetMap VARCHAR(16),
    IN in_nTargetX INT,
    IN in_nTargetY INT
)
sp_body: BEGIN
    DELETE FROM tActiveZoneTransfer WHERE dExpiresAt <= NOW();

    INSERT INTO tActiveZoneTransfer (
        sCharID, nWorldNo, nZoneNo, nUserNo, sUserID, nCharNo,
        nRandomID, nAdminLevel, sClientIP, sTargetMap, nTargetX, nTargetY,
        dIssuedAt, dExpiresAt
    ) VALUES (
        in_sCharID, in_nWorldNo, in_nZoneNo, in_nUserNo, in_sUserID, in_nCharNo,
        in_nRandomID, in_nAdminLevel, in_sClientIP, in_sTargetMap, in_nTargetX, in_nTargetY,
        NOW(), DATE_ADD(NOW(), INTERVAL 120 SECOND)
    )
    ON DUPLICATE KEY UPDATE
        nWorldNo = VALUES(nWorldNo),
        nZoneNo = VALUES(nZoneNo),
        nUserNo = VALUES(nUserNo),
        sUserID = VALUES(sUserID),
        nCharNo = VALUES(nCharNo),
        nRandomID = VALUES(nRandomID),
        nAdminLevel = VALUES(nAdminLevel),
        sClientIP = VALUES(sClientIP),
        sTargetMap = VALUES(sTargetMap),
        nTargetX = VALUES(nTargetX),
        nTargetY = VALUES(nTargetY),
        dIssuedAt = NOW(),
        dExpiresAt = DATE_ADD(NOW(), INTERVAL 120 SECOND);

    SELECT 0 AS nRet;
END""",

    "usp_Session_ConsumeZoneTransfer": """CREATE PROCEDURE `usp_Session_ConsumeZoneTransfer`(
    IN in_sCharID VARCHAR(40),
    IN in_nRandomID SMALLINT UNSIGNED,
    IN in_sClientIP VARCHAR(45)
)
sp_body: BEGIN
    DECLARE v_worldNo TINYINT UNSIGNED DEFAULT 0;
    DECLARE v_zoneNo SMALLINT UNSIGNED DEFAULT 0;
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_userID VARCHAR(60) DEFAULT '';
    DECLARE v_charNo INT DEFAULT 0;
    DECLARE v_charID VARCHAR(40) DEFAULT '';
    DECLARE v_randomID SMALLINT UNSIGNED DEFAULT 0;
    DECLARE v_adminLevel TINYINT UNSIGNED DEFAULT 0;
    DECLARE v_clientIP VARCHAR(45) DEFAULT '';
    DECLARE v_targetMap VARCHAR(16) DEFAULT '';
    DECLARE v_targetX INT DEFAULT 0;
    DECLARE v_targetY INT DEFAULT 0;
    DECLARE v_expiresAt DATETIME DEFAULT NULL;
    DECLARE v_ret INT DEFAULT -1;

    proc_logic: BEGIN
        DELETE FROM tActiveZoneTransfer WHERE dExpiresAt <= NOW();

        START TRANSACTION;

        SELECT nWorldNo, nZoneNo, nUserNo, sUserID, nCharNo, sCharID,
               nRandomID, nAdminLevel, sClientIP, sTargetMap, nTargetX, nTargetY, dExpiresAt
        INTO v_worldNo, v_zoneNo, v_userNo, v_userID, v_charNo, v_charID,
             v_randomID, v_adminLevel, v_clientIP, v_targetMap, v_targetX, v_targetY, v_expiresAt
        FROM tActiveZoneTransfer
        WHERE sCharID = in_sCharID
        FOR UPDATE;

        IF v_charNo = 0 OR v_expiresAt IS NULL THEN
            ROLLBACK;
            SET v_ret = -1;
            LEAVE proc_logic;
        END IF;

        IF in_nRandomID <> 0 AND v_randomID <> in_nRandomID THEN
            ROLLBACK;
            SET v_ret = -2;
            LEAVE proc_logic;
        END IF;

        IF in_sClientIP <> '' AND v_clientIP <> '' AND v_clientIP <> in_sClientIP THEN
            ROLLBACK;
            SET v_ret = -3;
            LEAVE proc_logic;
        END IF;

        DELETE FROM tActiveZoneTransfer WHERE sCharID = in_sCharID;
        COMMIT;
        SET v_ret = 0;
    END proc_logic;

    SELECT
        v_worldNo AS nWorldNo,
        v_zoneNo AS nZoneNo,
        v_userNo AS nUserNo,
        v_userID AS sUserID,
        v_charNo AS nCharNo,
        v_charID AS sCharID,
        v_adminLevel AS nAdminLevel,
        v_targetMap AS sTargetMap,
        v_targetX AS nTargetX,
        v_targetY AS nTargetY,
        v_ret AS nRet;
END""",

    "usp_WorldZone_UpsertNode": """CREATE PROCEDURE `usp_WorldZone_UpsertNode`(
    IN in_nWorldNo TINYINT UNSIGNED,
    IN in_nZoneNo SMALLINT UNSIGNED,
    IN in_sZoneIP VARCHAR(45),
    IN in_nZonePort SMALLINT UNSIGNED
)
sp_body: BEGIN
    INSERT INTO tWorldZoneNode (nWorldNo, nZoneNo, sZoneIP, nZonePort, nCurrentPlayers, bIsOnline, dLastHeartbeat)
    VALUES (in_nWorldNo, in_nZoneNo, in_sZoneIP, in_nZonePort, 0, 1, NOW())
    ON DUPLICATE KEY UPDATE
        sZoneIP = VALUES(sZoneIP),
        nZonePort = VALUES(nZonePort),
        bIsOnline = 1,
        dLastHeartbeat = NOW();

    SELECT 0 AS nRet;
END""",

    "usp_WorldZone_Heartbeat": """CREATE PROCEDURE `usp_WorldZone_Heartbeat`(
    IN in_nWorldNo TINYINT UNSIGNED,
    IN in_nZoneNo SMALLINT UNSIGNED,
    IN in_nCurrentPlayers INT
)
sp_body: BEGIN
    UPDATE tWorldZoneNode
    SET nCurrentPlayers = in_nCurrentPlayers,
        bIsOnline = 1,
        dLastHeartbeat = NOW()
    WHERE nWorldNo = in_nWorldNo AND nZoneNo = in_nZoneNo;

    SELECT ROW_COUNT() AS nRet;
END""",

    "usp_WorldZone_GetList": """CREATE PROCEDURE `usp_WorldZone_GetList`(
    IN in_nWorldNo TINYINT UNSIGNED
)
sp_body: BEGIN
    SELECT nWorldNo, nZoneNo, sZoneIP, nZonePort, nCurrentPlayers, bIsOnline, dLastHeartbeat
    FROM tWorldZoneNode
    WHERE nWorldNo = in_nWorldNo AND bIsOnline = 1
    ORDER BY nZoneNo ASC;
END""",

    "usp_WorldZone_ResolveMapEndpoint": """CREATE PROCEDURE `usp_WorldZone_ResolveMapEndpoint`(
    IN in_nWorldNo TINYINT UNSIGNED,
    IN in_nZoneNo SMALLINT UNSIGNED
)
sp_body: BEGIN
    DECLARE v_ip VARCHAR(45) DEFAULT '';
    DECLARE v_port SMALLINT UNSIGNED DEFAULT 0;
    DECLARE v_online TINYINT DEFAULT 0;

    SELECT sZoneIP, nZonePort, bIsOnline
    INTO v_ip, v_port, v_online
    FROM tWorldZoneNode
    WHERE nWorldNo = in_nWorldNo AND nZoneNo = in_nZoneNo AND bIsOnline = 1
    LIMIT 1;

    IF v_port = 0 THEN
        SELECT sZoneIP, nZonePort, bIsOnline
        INTO v_ip, v_port, v_online
        FROM tWorldZoneNode
        WHERE nWorldNo = in_nWorldNo AND nZoneNo = 0 AND bIsOnline = 1
        LIMIT 1;
    END IF;

    SELECT v_ip AS sZoneIP, v_port AS nZonePort, v_online AS bIsOnline;
END""",
}

SPECIAL_PROCEDURES = {k.lower(): (k, v) for k, v in SPECIAL_PROCEDURES_RAW.items()}


def tokenize(sql):
    token_spec = [
        ("STRING", r"'([^']|'')*'"),
        ("BACKTICK", r"`[^`]+`"),
        ("NUMBER", r"\b\d+(\.\d*)?\b"),
        ("HEX", r"0x[0-9a-fA-F]+"),
        ("WORD", r"\b[a-zA-Z_][a-zA-Z0-9_]*\b"),
        ("OPERATOR", r"<>|<=|>=|!=|==|:=|[+\-*/=<>!&|^]"),
        ("DELIM", r"[(),;.]"),
        ("WS", r"\s+"),
    ]
    tok_regex = "|".join("(?P<%s>%s)" % pair for pair in token_spec)
    tokens = []
    for mo in re.finditer(tok_regex, sql):
        kind = mo.lastgroup
        val = mo.group()
        if kind != "WS":
            tokens.append((kind, val))
    return tokens


def parse_body_statements(tokens):
    idx = 0
    n = len(tokens)

    STMT_KEYWORDS = {
        "SET", "SELECT", "INSERT", "UPDATE", "DELETE", "CALL",
        "LEAVE", "START", "COMMIT", "ROLLBACK", "PREPARE", "EXECUTE", "DEALLOCATE",
    }
    BLOCK_KEYWORDS = {"IF", "ELSE", "BEGIN", "END"}

    def peek(offset=0):
        if idx + offset < n:
            return tokens[idx + offset][1].upper()
        return None

    def get_token():
        nonlocal idx
        if idx < n:
            t = tokens[idx]
            idx += 1
            return t
        return None

    def parse_statement():
        nonlocal idx
        first = peek()
        if not first:
            return None

        if first == "BEGIN":
            get_token()
            stmts = []
            while peek() and peek() != "END":
                s = parse_statement()
                if s:
                    stmts.append(s)
            if peek() == "END":
                get_token()
            return "\n        ".join(stmts)

        if first == "IF":
            get_token()
            cond_tokens = []
            depth = 0
            while idx < n:
                w = peek()
                if w == "(":
                    depth += 1
                elif w == ")":
                    depth -= 1

                if depth == 0:
                    if w == "BEGIN":
                        get_token()
                        then_stmts = []
                        while peek() and peek() != "END":
                            s = parse_statement()
                            if s:
                                then_stmts.append(s)
                        if peek() == "END":
                            get_token()

                        else_part = ""
                        if peek() == "ELSE":
                            get_token()
                            if peek() == "BEGIN":
                                get_token()
                                e_stmts = []
                                while peek() and peek() != "END":
                                    es = parse_statement()
                                    if es:
                                        e_stmts.append(es)
                                if peek() == "END":
                                    get_token()
                                else_part = "ELSE\n            " + "\n            ".join(e_stmts)
                            else:
                                es = parse_statement()
                                if es:
                                    else_part = f"ELSE\n            {es}"

                        cond_str = " ".join(t[1] for t in cond_tokens)
                        cond_str = re.sub(r"\bCAST\s+\(", "CAST(", cond_str, flags=re.IGNORECASE)
                        then_body = "\n            ".join(then_stmts)
                        if else_part:
                            return f"IF {cond_str} THEN\n            {then_body}\n        {else_part}\n        END IF;"
                        return f"IF {cond_str} THEN\n            {then_body}\n        END IF;"

                    elif w in STMT_KEYWORDS or w == "IF":
                        then_stmt = parse_statement()
                        else_part = ""
                        if peek() == "ELSE":
                            get_token()
                            if peek() == "BEGIN":
                                get_token()
                                e_stmts = []
                                while peek() and peek() != "END":
                                    es = parse_statement()
                                    if es:
                                        e_stmts.append(es)
                                if peek() == "END":
                                    get_token()
                                else_part = "ELSE\n            " + "\n            ".join(e_stmts)
                            else:
                                es = parse_statement()
                                if es:
                                    else_part = f"ELSE\n            {es}"

                        cond_str = " ".join(t[1] for t in cond_tokens)
                        cond_str = re.sub(r"\bCAST\s+\(", "CAST(", cond_str, flags=re.IGNORECASE)
                        if else_part:
                            return f"IF {cond_str} THEN\n            {then_stmt}\n        {else_part}\n        END IF;"
                        return f"IF {cond_str} THEN\n            {then_stmt}\n        END IF;"

                cond_tokens.append(get_token())

        if first == "ELSE":
            get_token()
            stmt = parse_statement()
            return f"ELSE\n        {stmt}"

        stmt_tokens = []
        depth = 0
        in_update = first == "UPDATE"
        in_insert = first == "INSERT"
        saw_set_in_update = False

        while idx < n:
            w = peek()
            if w == "(":
                depth += 1
            elif w == ")":
                depth -= 1
            elif w == ";":
                get_token()
                break

            if depth == 0 and stmt_tokens:
                if in_update and w == "SET" and not saw_set_in_update:
                    saw_set_in_update = True
                elif in_insert and w == "SELECT":
                    pass
                elif w in STMT_KEYWORDS or w in BLOCK_KEYWORDS:
                    break

            stmt_tokens.append(get_token())

        stmt_str = " ".join(t[1] for t in stmt_tokens).strip()
        if stmt_str.rstrip(";").isdigit() or stmt_str.rstrip(";") in ("<", ">", ")", ","):
            return None

        stmt_str = re.sub(r"\s*([(),;.=<>!+*/-])\s*", r"\1", stmt_str)
        stmt_str = re.sub(
            r"\b(INTO|FROM|WHERE|AND|OR|SET|VALUES|ORDER BY|GROUP BY|LIMIT|JOIN|INNER JOIN|LEFT JOIN|AS)\b",
            r" \1 ",
            stmt_str,
            flags=re.IGNORECASE,
        )
        stmt_str = re.sub(r"\bCAST\s+\(", "CAST(", stmt_str, flags=re.IGNORECASE)
        stmt_str = re.sub(r"\s+", " ", stmt_str).strip()

        if not stmt_str.endswith(";"):
            stmt_str += ";"

        if stmt_str.upper().startswith("CALL "):
            return wrap_nested_call(stmt_str)

        return stmt_str

    all_stmts = []
    while idx < n:
        s = parse_statement()
        if s:
            all_stmts.append(s)
    return "\n        ".join(all_stmts)


def inject_not_found_handler(sql_text):
    if " INTO " not in sql_text.upper() or "CONTINUE HANDLER FOR NOT FOUND" in sql_text.upper():
        return sql_text
    lines = sql_text.split("\n")
    begin_idx = -1
    last_decl_idx = -1
    for i, ln in enumerate(lines):
        s = ln.strip().upper()
        if begin_idx == -1 and s.endswith("BEGIN") and "SP_BODY" in s:
            begin_idx = i
            last_decl_idx = i
        elif begin_idx != -1:
            if s.startswith("DECLARE "):
                last_decl_idx = i
            elif s == "":
                continue
            else:
                break
    if last_decl_idx != -1:
        lines.insert(last_decl_idx + 1, "    DECLARE CONTINUE HANDLER FOR NOT FOUND BEGIN END;")
    return "\n".join(lines)


def transpile_procedure(sp_name):
    key = sp_name.lower()
    if key in SPECIAL_PROCEDURES:
        orig_name, spec_sql = SPECIAL_PROCEDURES[key]
        if orig_name != sp_name:
            spec_sql = re.sub(
                rf"CREATE PROCEDURE `{re.escape(orig_name)}`",
                f"CREATE PROCEDURE `{sp_name}`",
                spec_sql,
                count=1,
            )
        return inject_not_found_handler(spec_sql)

    bak_base = bak_files_by_lower.get(key)
    if not bak_base:
        raise ValueError(f"Missing SQL source and SPECIAL_PROCEDURES entry for {sp_name}")

    filepath = os.path.join(TARGET_DIR, bak_base + ".sql")
    with open(filepath, "r", encoding="latin-1") as fp:
        raw = fp.read()

    clean = clean_sql_source(raw)
    # Strip SQL comments BEFORE matching header AS so inline comments like `-- As Master` never split header
    clean_no_comments = re.sub(r"/\*.*?\*/", "", clean, flags=re.DOTALL)
    clean_no_comments = re.sub(r"--[^\n]*", "", clean_no_comments)

    m = re.search(
        r"CREATE\s+PROC(?:EDURE)?\s+(?:\[?dbo\]?\.)?\[?([a-zA-Z0-9_]+)\]?(.*?)\bAS\b(.*)",
        clean_no_comments,
        flags=re.DOTALL | re.IGNORECASE,
    )
    if not m:
        raise ValueError(f"Failed to match procedure header for {sp_name}")

    param_text = m.group(2)
    body = m.group(3)

    # Support BOTH `@p TYPE OUTPUT = 0` AND `@p TYPE = 0 OUTPUT`
    raw_params = re.findall(
        r"@([a-zA-Z0-9_]+)\s+([a-zA-Z0-9_]+(?:\([^\)]+\))?)((?:\s*(?:=\s*[^,\n]+|\bOUTPUT\b|\bOUT\b))*)",
        param_text,
        re.IGNORECASE,
    )

    mysql_params = []
    param_map = {}
    out_params = []

    for pname, ptype, pmods in raw_params:
        is_out = bool(re.search(r"\b(?:OUTPUT|OUT)\b", pmods or "", re.IGNORECASE))
        m_type = map_type(ptype)
        if is_out:
            mapped_name = "out_" + pname
            mysql_params.append(f"OUT {mapped_name} {m_type}")
            param_map[pname.lower()] = mapped_name
            out_params.append(mapped_name)
        else:
            mapped_name = "in_" + pname
            mysql_params.append(f"IN {mapped_name} {m_type}")
            param_map[pname.lower()] = mapped_name

    # Cross-database prefixes (Account..tUser -> tUser, Account.dbo.X -> X)
    body = re.sub(r"\b(?:[a-zA-Z0-9_]+\.)*dbo\.", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\b[a-zA-Z0-9_]+\.\.([a-zA-Z0-9_]+)", r"\1", body)

    # Remove T-SQL PRINT statements
    body = re.sub(r"\bPRINT\s*\([^\n;]*\)\s*;?", "", body, flags=re.IGNORECASE)
    body = re.sub(r"^[ \t]*PRINT\b[^\n]*$", "", body, flags=re.IGNORECASE | re.MULTILINE)

    # Bracketed identifiers [Date] -> `Date`
    body = re.sub(r"\[([a-zA-Z0-9_]+)\]", r"`\1`", body)

    # T-SQL CONVERT(NVARCHAR(10), GETDATE(), 120) -> DATE_FORMAT(NOW(), '%Y-%m-%d')
    body = re.sub(
        r"\bCONVERT\s*\(\s*N?VARCHAR\s*\(\s*10\s*\)\s*,\s*(?:GETDATE\(\)|NOW\(\))\s*,\s*120\s*\)",
        "DATE_FORMAT(NOW(), '%Y-%m-%d')",
        body,
        flags=re.IGNORECASE,
    )

    # INSERT table (...) without INTO -> INSERT INTO table (...)
    body = re.sub(r"\bINSERT\s+(?!INTO\b|IGNORE\b)([`a-zA-Z0-9_]+)", r"INSERT INTO \1", body, flags=re.IGNORECASE)

    # Remove T-SQL settings and hints
    body = re.sub(r"\bSET\s+NOCOUNT\s+(?:ON|OFF)\b;?", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bSET\s+LOCK_TIMEOUT\s+-?\d+;?", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bWITH\s*\([^\)]+\)", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\b\(NOLOCK\)", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGETDATE\(\)", "NOW()", body, flags=re.IGNORECASE)
    body = re.sub(r"\bCURRENT_TIMESTAMP\b", "NOW()", body, flags=re.IGNORECASE)
    body = re.sub(r"\bISNULL\(", "IFNULL(", body, flags=re.IGNORECASE)
    body = re.sub(r"\bLEN\(", "CHAR_LENGTH(", body, flags=re.IGNORECASE)
    body = re.sub(r"@@IDENTITY\b", "LAST_INSERT_ID()", body, flags=re.IGNORECASE)
    body = re.sub(r"\bSCOPE_IDENTITY\(\)", "LAST_INSERT_ID()", body, flags=re.IGNORECASE)
    body = re.sub(r"@@ROWCOUNT\b", "ROW_COUNT()", body, flags=re.IGNORECASE)
    body = re.sub(r"@@TRANCOUNT\b", "1", body, flags=re.IGNORECASE)

    # Cast conversions
    body = re.sub(r"\bCAST\s*\(\s*(.+?)\s+AS\s+INT\s*\)", r"CAST(\1 AS SIGNED)", body, flags=re.IGNORECASE)
    body = re.sub(r"\bCAST\s*\(\s*(.+?)\s+AS\s+N?VARCHAR(?:\([^\)]+\))?\s*\)", r"CAST(\1 AS CHAR)", body, flags=re.IGNORECASE)

    # Redundant FROM in UPDATE fix
    body = re.sub(r"\bUPDATE\s+([a-zA-Z0-9_]+)\s+SET\s+([^;\n]+?)\s+FROM\s+\1\b", r"UPDATE \1 SET \2", body, flags=re.IGNORECASE)

    # Strip TOP inside EXISTS(...)
    body = re.sub(r"\bEXISTS\s*\(\s*SELECT\s+TOP\s+\(?\d+\)?\s+", "EXISTS(SELECT ", body, flags=re.IGNORECASE)

    # DATEADD conversion
    def dateadd_repl(dam):
        unit = dam.group(1).upper()
        val = dam.group(2).strip()
        date_arg = dam.group(3).strip()
        return f"DATE_ADD({date_arg}, INTERVAL {val} {unit})"

    body = re.sub(
        r"\bDATEADD\s*\(\s*([a-zA-Z]+)\s*,\s*([^,]+)\s*,\s*((?:[^()]+|\([^()]*\))+)\)",
        dateadd_repl,
        body,
        flags=re.IGNORECASE,
    )

    # Common T-SQL error patterns (preserve ROW_COUNT() across SET @var = ROW_COUNT())
    body = re.sub(r"IF\s+@@ERROR\s*<>\s*0\s+OR\s+ROW_COUNT\(\)\s*=\s*0", "IF ROW_COUNT() = 0", body, flags=re.IGNORECASE)
    body = re.sub(r"IF\s+@@ERROR\s*<>\s*0\s+OR\s+ROW_COUNT\(\)\s*<>\s*1", "IF ROW_COUNT() <> 1", body, flags=re.IGNORECASE)
    body = re.sub(r"IF\s+@@ERROR\s*=\s*0\s+AND\s+ROW_COUNT\(\)\s*=\s*1", "IF ROW_COUNT() = 1", body, flags=re.IGNORECASE)
    body = re.sub(r"@@ERROR\b", "0", body, flags=re.IGNORECASE)
    body = re.sub(
        r"\bSET\s+(@?[a-zA-Z0-9_]+)\s*=\s*ROW_COUNT\(\)\s*;?\s*IF\s+ROW_COUNT\(\)",
        r"SET \1 = ROW_COUNT(); IF \1",
        body,
        flags=re.IGNORECASE,
    )

    body = re.sub(r"\bBEGIN\s+TRAN(?:SACTION)?\b", "START TRANSACTION;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bCOMMIT\s+TRAN(?:SACTION)?\b", "COMMIT;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bROLLBACK\s+TRAN(?:SACTION)?\b", "ROLLBACK;", body, flags=re.IGNORECASE)

    # Hoist DECLARE statements
    mysql_declares = []
    declare_blocks = re.findall(
        r"\bDECLARE\b(.*?)(?=\b(?:SET|BEGIN|UPDATE|SELECT|INSERT|DELETE|IF|CALL|START|COMMIT|ROLLBACK|EXEC|DECLARE)\b)",
        body,
        flags=re.DOTALL | re.IGNORECASE,
    )
    body = re.sub(
        r"\bDECLARE\b.*?(?=\b(?:SET|BEGIN|UPDATE|SELECT|INSERT|DELETE|IF|CALL|START|COMMIT|ROLLBACK|EXEC|DECLARE)\b)",
        "",
        body,
        flags=re.DOTALL | re.IGNORECASE,
    )

    for decl_str in declare_blocks:
        decl_clean = decl_str.strip().rstrip(";")
        parts = []
        cur = []
        depth = 0
        for c in decl_clean:
            if c == "(":
                depth += 1
            elif c == ")":
                depth -= 1
            elif c == "," and depth == 0:
                parts.append("".join(cur).strip())
                cur = []
                continue
            cur.append(c)
        if cur:
            parts.append("".join(cur).strip())

        for p in parts:
            dm = re.match(
                r"@([a-zA-Z0-9_]+)(?:\s+AS)?\s+([a-zA-Z0-9_]+(?:\([^\)]+\))?)(?:\s*=\s*(.+))?",
                p.strip(),
                re.IGNORECASE,
            )
            if dm:
                dname = dm.group(1)
                dtype = dm.group(2)
                dinit = dm.group(3)
                m_type = map_type(dtype)
                mapped_dname = "v_" + dname
                param_map[dname.lower()] = mapped_dname
                if dinit:
                    init_val = dinit.strip()
                    for r_var, m_var in param_map.items():
                        init_val = re.sub(rf"@{r_var}\b", m_var, init_val, flags=re.IGNORECASE)
                    init_val = init_val.replace("GETDATE()", "NOW()")
                    mysql_declares.append(f"    DECLARE {mapped_dname} {m_type} DEFAULT {init_val};")
                else:
                    mysql_declares.append(f"    DECLARE {mapped_dname} {m_type};")

    # Replace GOTOs and RETURN with LEAVE proc_logic; (never skips final SELECT!)
    body = re.sub(r"\bGOTO\s+ERROR_[ \t]*;?", "ROLLBACK; SET out_nRet = 1; LEAVE proc_logic;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO\s+END_[ \t]*;?", "LEAVE proc_logic;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO\s+ERR[ \t]*;?", "ROLLBACK; SET out_err = 1; LEAVE proc_logic;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO[ \t]+[a-zA-Z0-9_]+[ \t]*;?", "LEAVE proc_logic;", body, flags=re.IGNORECASE)
    body = re.sub(r"^[ \t]*[a-zA-Z0-9_]+:[ \t]*.*$", "", body, flags=re.MULTILINE)

    # EXEC calls -> CALL `proc`(args);
    def exec_repl(em):
        proc = em.group(1)
        args_str = em.group(2) or ""
        args_str = re.sub(r"\s+", " ", args_str).strip()
        raw_args = [a.strip() for a in args_str.split(",") if a.strip()]
        clean_args = []
        for ra in raw_args:
            ra = re.sub(r"\b(?:OUTPUT|OUT)\b", "", ra, flags=re.IGNORECASE).strip()
            for r_var, m_var in param_map.items():
                ra = re.sub(rf"@{r_var}\b", m_var, ra, flags=re.IGNORECASE)
            # If literal passed into an OUT param (e.g. EXEC p_Item_SetOption @k, 300, 1, 0), bind dummy @var
            clean_args.append(ra)
        if proc.lower() == "p_item_setoption" and len(clean_args) == 4 and clean_args[3].isdigit():
            clean_args[3] = "@_xc_dummy_out"
        return f"CALL `{proc}`({', '.join(clean_args)});"

    body = re.sub(
        r"\bEXEC(?:UTE)?\s+([a-zA-Z0-9_]+)\s+([^;]+?)(?=(?:\b(?:SET|IF|ELSE|BEGIN|END|RETURN|GOTO|COMMIT|ROLLBACK|EXEC|START)\b|$))",
        exec_repl,
        body,
        flags=re.DOTALL | re.IGNORECASE,
    )

    # SELECT @var = col FROM ... -> SELECT col INTO var FROM ... (parenthesis-aware for subqueries)
    def convert_single_select_assign(full_stmt):
        full_stmt = full_stmt.strip()
        # Find top-level FROM (depth == 0)
        from_pos = -1
        d = 0
        in_str = False
        i = 0
        L = len(full_stmt)
        while i < L:
            ch = full_stmt[i]
            if in_str:
                if ch == "'":
                    in_str = False
                i += 1
                continue
            if ch == "'":
                in_str = True
            elif ch == "(":
                d += 1
            elif ch == ")":
                d = max(0, d - 1)
            elif d == 0 and full_stmt[i : i + 4].upper() == "FROM":
                before_ok = i == 0 or not (full_stmt[i - 1].isalnum() or full_stmt[i - 1] == "_")
                after_ok = (i + 4 == L) or not (full_stmt[i + 4].isalnum() or full_stmt[i + 4] == "_")
                if before_ok and after_ok:
                    from_pos = i
                    break
            i += 1

        if from_pos != -1:
            head_part = full_stmt[:from_pos]
            rest_part = full_stmt[from_pos:].strip()
            assigns_part = re.sub(r"^\s*SELECT(?:\s+TOP\s+\(?\d+\)?)?\s+", "", head_part, flags=re.IGNORECASE).strip()

            top_match = re.search(r"\bTOP\s+\(?(\d+)\)?\b", head_part, re.IGNORECASE)
            limit_clause = f" LIMIT {top_match.group(1)}" if top_match else ""

            pairs = []
            cur = []
            depth = 0
            for char in assigns_part:
                if char == "(":
                    depth += 1
                elif char == ")":
                    depth -= 1
                elif char == "," and depth == 0:
                    pairs.append("".join(cur).strip())
                    cur = []
                    continue
                cur.append(char)
            if cur:
                pairs.append("".join(cur).strip())

            cols = []
            vars_ = []
            valid = True
            for p in pairs:
                m_pair = re.match(r"@?([a-zA-Z0-9_]+)\s*=\s*(.+)", p, re.DOTALL)
                if m_pair:
                    v_name = m_pair.group(1).lower()
                    mapped_v = param_map.get(v_name, "v_" + v_name)
                    vars_.append(mapped_v)
                    cols.append(m_pair.group(2).strip())
                else:
                    valid = False
                    break
            if valid:
                rest_clean = rest_part.rstrip("; \t\n")
                return f"SELECT {', '.join(cols)} INTO {', '.join(vars_)} {rest_clean}{limit_clause};"
        else:
            assigns_part = re.sub(r"^\s*SELECT\s+", "", full_stmt, flags=re.IGNORECASE).rstrip("; \t\n")
            return f"SET {assigns_part};"
        return full_stmt

    boundary_kw = {
        "SELECT", "SET", "INSERT", "UPDATE", "DELETE", "IF", "ELSE",
        "BEGIN", "END", "RETURN", "GOTO", "COMMIT", "ROLLBACK", "CALL", "START", "LEAVE",
    }
    sel_pat = re.compile(r"\bSELECT(?:\s+TOP\s+\(?\d+\)?)?\s+@?[a-zA-Z0-9_]+\s*=", re.IGNORECASE)
    out_chunks = []
    pos = 0
    L_body = len(body)
    while pos < L_body:
        m_sel = sel_pat.search(body, pos)
        if not m_sel:
            out_chunks.append(body[pos:])
            break
        start_idx = m_sel.start()
        out_chunks.append(body[pos:start_idx])
        # Scan from m_sel.end() respecting parentheses and quotes until top-level boundary
        idx_scan = m_sel.end()
        d = 0
        in_str = False
        end_idx = L_body
        while idx_scan < L_body:
            ch = body[idx_scan]
            if in_str:
                if ch == "'":
                    in_str = False
                idx_scan += 1
                continue
            if ch == "'":
                in_str = True
                idx_scan += 1
                continue
            if ch == "(":
                d += 1
                idx_scan += 1
                continue
            if ch == ")":
                if d == 0:
                    end_idx = idx_scan
                    break
                d -= 1
                idx_scan += 1
                continue
            if d == 0:
                if ch == ";":
                    end_idx = idx_scan + 1
                    break
                if ch.isalpha() or ch == "_":
                    w_end = idx_scan + 1
                    while w_end < L_body and (body[w_end].isalnum() or body[w_end] == "_"):
                        w_end += 1
                    word_up = body[idx_scan:w_end].upper()
                    before_ok = idx_scan == 0 or not (body[idx_scan - 1].isalnum() or body[idx_scan - 1] in ("_", "@"))
                    if before_ok and word_up in boundary_kw:
                        end_idx = idx_scan
                        break
                    idx_scan = w_end
                    continue
            idx_scan += 1
        stmt_slice = body[start_idx:end_idx]
        out_chunks.append(convert_single_select_assign(stmt_slice) + " ")
        pos = end_idx
    body = "".join(out_chunks)

    # SELECT TOP (n) / SELECT TOP n without assign -> SELECT ... LIMIT n
    def select_top_repl(stm):
        n_lim = stm.group(1).strip("() ")
        for r_var, m_var in param_map.items():
            n_lim = re.sub(rf"@{r_var}\b", m_var, n_lim, flags=re.IGNORECASE)
        rest = stm.group(2).rstrip("; \t\n")
        return f"SELECT {rest} LIMIT {n_lim};"

    body = re.sub(
        r"\bSELECT\s+TOP\s+(\(@?[a-zA-Z0-9_]+\)|\d+)\s+([^;]+?)(?=(?:\b(?:SET|INSERT|UPDATE|DELETE|IF|ELSE|BEGIN|END|RETURN|GOTO|COMMIT|ROLLBACK|CALL|START)\b|$))",
        select_top_repl,
        body,
        flags=re.IGNORECASE | re.DOTALL,
    )

    # Variable substitution
    for raw_var, mapped_var in param_map.items():
        body = re.sub(rf"@{raw_var}\b", mapped_var, body, flags=re.IGNORECASE)

    # Initialize OUT parameters at start of proc_logic so NULL is never returned on early LEAVE
    out_inits = []
    for p_decl in mysql_params:
        if p_decl.startswith("OUT "):
            parts = p_decl.split()
            oname = parts[1]
            otype = parts[2].upper()
            if "CHAR" in otype or "TEXT" in otype:
                out_inits.append(f"SET {oname} = '';")
            elif "DATE" in otype or "TIME" in otype:
                out_inits.append(f"SET {oname} = NOW();")
            else:
                out_inits.append(f"SET {oname} = 0;")

    # Replace RETURN with LEAVE proc_logic;
    body = re.sub(r"\bRETURN(?:[ \t]+[^;\n\r]+)?[ \t]*;?", "LEAVE proc_logic;", body, flags=re.IGNORECASE)

    tokens = tokenize(body)
    parsed_body = parse_body_statements(tokens)

    if " INTO " in parsed_body.upper():
        mysql_declares.append("    DECLARE CONTINUE HANDLER FOR NOT FOUND BEGIN END;")

    # Classify resultset-emitting SELECT statements inside proc_logic
    has_table_resultset = False
    has_scalar_resultset = False
    body_lines = parsed_body.split("\n")
    rewritten_lines = []
    for bl in body_lines:
        stripped_up = bl.strip().upper()
        if stripped_up.startswith("SELECT ") and " INTO " not in stripped_up:
            if " FROM " in stripped_up:
                has_table_resultset = True
                rewritten_lines.append(bl)
            else:
                has_scalar_resultset = True
                indent = bl[: len(bl) - len(bl.lstrip())]
                rewritten_lines.append(f"{indent}SET @_xc_emitted_rs = 1;\n{bl}")
        else:
            rewritten_lines.append(bl)
    parsed_body = "\n".join(rewritten_lines)

    params_sql = ",\n    ".join(mysql_params)
    if params_sql:
        params_sql = "\n    " + params_sql + "\n"

    declares_sql = "\n".join(mysql_declares)
    if declares_sql:
        declares_sql = "\n" + declares_sql + "\n"

    if has_scalar_resultset and out_params:
        out_inits.insert(0, "SET @_xc_emitted_rs = 0;")

    init_sql = ""
    if out_inits:
        init_sql = "\n        ".join(out_inits) + "\n        "

    out_select = ""
    if out_params and not has_table_resultset:
        cols = ", ".join(f"{p} AS {p.replace('out_', '')}" for p in out_params)
        if has_scalar_resultset:
            out_select = (
                "\n    IF IFNULL(@_xc_nested_sp, 0) = 0 AND IFNULL(@_xc_emitted_rs, 0) = 0 THEN\n"
                f"        SELECT {cols};\n"
                "    END IF;\n"
            )
        else:
            out_select = (
                "\n    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN\n"
                f"        SELECT {cols};\n"
                "    END IF;\n"
            )

    sql = f"""CREATE PROCEDURE `{sp_name}`({params_sql})
sp_body: BEGIN{declares_sql}
    proc_logic: BEGIN
        {init_sql}{parsed_body}
    END proc_logic;{out_select}
END"""
    return sql


if __name__ == "__main__":
    conn = mysql.connector.connect(
        host="127.0.0.1",
        port=3306,
        user="root",
        password="lerd",
        database="xelion_online",
    )
    cursor = conn.cursor()

    success = 0
    errors = []
    all_sqls = [
        """-- ============================================================================
-- XELION ONLINE - STORED PROCEDURES OFICIALES Y OPERATIVOS DE WORLD00_CHARACTER (MYSQL 8.4)
-- Catálogo exhaustivo de 425 Stored Procedures (419 Oficiales + 6 Operativos S2S)
-- Unión exacta: 311 ODBC de Character.exe U 405 no-DBA de World00_Character.bak + 6 S2S
-- 100% Pool-Safe (proc_logic + SELECT incondicional para *sql.DB en Go)
-- ============================================================================

USE `xelion_online`;
"""
    ]

    for sp in TARGET_425:
        try:
            sp_sql = transpile_procedure(sp)
            cursor.execute(f"DROP PROCEDURE IF EXISTS `{sp}`")
            cursor.execute(sp_sql)
            conn.commit()
            success += 1

            file_sp_sql = f"""-- ----------------------------------------------------------------------------
-- Procedure: {sp}
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `{sp}`;
DELIMITER $$
{sp_sql}$$
DELIMITER ;
"""
            all_sqls.append(file_sp_sql)
        except Exception as e:
            errors.append((sp, str(e)))

    with open("database/procedures_character.sql", "w", encoding="utf-8") as out_fp:
        out_fp.write("\n".join(all_sqls))

    shutil.copyfile("scripts/transpile_character_sps.py", "scripts/build_procedures_character.py")

    print(f"Results: {success}/{len(TARGET_425)} procedures compiled cleanly in MySQL 8.4!")
    print(f"Saved complete script with {success} Stored Procedures to database/procedures_character.sql.")
    if errors:
        print(f"Errors ({len(errors)}):")
        for sp, err in errors:
            print(f"  {sp:40s}: {err.replace(chr(10), ' ')[:120]}")
    conn.close()
    if errors:
        raise SystemExit(1)
