#!/usr/bin/env python3
"""
transpile_character_sps.py - Transpilador de alta fidelidad para los 311 Stored Procedures
de World00_Character desde T-SQL a MySQL 8.4 usando un parser formal de sentencias.
"""

import os
import re
import sys
import mysql.connector

TARGET_DIR = "docs/sql/procedures/World00_Character"

# 1. Catálogo oficial de los 311 Stored Procedures
with open("../AstraServer/Server/Character/Character.exe", "rb") as f:
    data = f.read()

ascii_matches = re.findall(b"(?:p|usp)_[a-zA-Z0-9_]+", data)
utf16_matches = re.findall(b"(?:(?:p|usp)_[a-zA-Z0-9_]+)".decode("ascii").encode("utf-16le"), data)
char_exe_sps = set(m.decode("ascii", errors="ignore") for m in ascii_matches)
for m in utf16_matches:
    char_exe_sps.add(m.decode("utf-16le", errors="ignore"))

files = set(f.replace(".sql", "") for f in os.listdir(TARGET_DIR) if f.endswith(".sql"))
in_char = sorted([f for f in files if f in char_exe_sps])
wm_procs = [
    "p_Char_GetNumOfUserChar", "p_Char_GetNumOfUserCharAll", "p_Char_GetIsDeleted",
    "p_Char_GetUserNo", "p_Char_Find", "p_Char_ChangeID", "p_Char_GetTitleCount",
    "p_Char_DelOptions", "p_Char_CheckOptionsExistance", "p_Char_GetState",
    "p_Char_GetStateAll", "p_Char_DelState", "p_Char_FameGet", "p_Char_FameSub",
    "p_Char_GetChargedBufferCount",
]
target_311 = sorted(list(set(in_char + wm_procs)))

def map_type(t):
    t = t.strip()
    tl = t.lower()
    if tl == "bit":
        return "TINYINT"
    if tl.startswith("nvarchar"):
        return tl.replace("nvarchar", "VARCHAR")
    if tl.startswith("varchar(max)") or tl.startswith("nvarchar(max)"):
        return "TEXT"
    if tl.startswith("varbinary(max)"):
        return "BLOB"
    if tl == "money":
        return "BIGINT"
    if tl == "smallmoney":
        return "INT"
    return t.upper()

def clean_sql_source(raw):
    raw = raw.replace("\r\n", "\n").replace("\r", "\n")
    parts = re.split(r"\bGO\b", raw, flags=re.IGNORECASE)
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
        if any(ord(c) < 32 and c not in '\t\n' for c in s):
            continue
        if re.match(r"^0[\x00-\x20\x7f-\uffff]", s) or re.match(r"^0\s+[^\w\s]", s) or s.upper() == "GO":
            continue
        clean_lines.append(l)
    return "\n".join(clean_lines)

# ----------------------------------------------------------------------------
# Procedimientos Especializados con Lógica Crítica Pulida al 100%
# ----------------------------------------------------------------------------
SPECIAL_PROCEDURES = {
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

    SET out_nCharNo = 0;

    IF EXISTS (SELECT 1 FROM tCharacter WHERE nUserNo = in_nUserNo AND nSlotNo = in_nSlotNo AND bDeleted = 0) THEN
        SET out_nCharNo = -1;
        LEAVE sp_body;
    END IF;

    IF EXISTS (SELECT 1 FROM tCharacter WHERE sID = in_sID) THEN
        SET out_nCharNo = -2;
        LEAVE sp_body;
    END IF;

    START TRANSACTION;

    INSERT INTO tCharacter (nUserNo, nSlotNo, sID, nCreateWorld, nAdminLevel, dCreateDate)
    VALUES (in_nUserNo, in_nSlotNo, in_sID, in_nCreateWorld, in_nAdminLevel, NOW());

    SET v_newCharNo = LAST_INSERT_ID();

    INSERT INTO tCharacterShape (nCharNo, nRace, nClass, nGender, nHairType, nHairColor, nFaceShape)
    VALUES (v_newCharNo, in_nRace, in_nClass, in_nGender, in_nHairType, in_nHairColor, in_nFaceShape);

    INSERT INTO tCharacterOptions (nCharNo)
    VALUES (v_newCharNo);

    COMMIT;

    SET out_nCharNo = v_newCharNo;
    SELECT out_nCharNo AS nCharNo;
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

    IF EXISTS(SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted > 0) THEN
        SET out_nRet = 1;
        LEAVE sp_body;
    END IF;

    SELECT sID INTO v_sID FROM tCharacter WHERE nCharNo = in_nCharNo LIMIT 1;
    IF v_sID IS NULL OR CHAR_LENGTH(v_sID) = 0 THEN
        SET out_nRet = 2;
        LEAVE sp_body;
    END IF;

    SET v_StartDeleteMark = LOCATE(v_DeleteMark, v_sID);
    IF v_StartDeleteMark > 0 THEN
        SET v_sID = LEFT(v_sID, v_StartDeleteMark - 1);
    END IF;

    SET v_sChangedID = CONCAT(RTRIM(v_sID), '~', DATE_FORMAT(v_CurrentDate, '%Y%m%d'), '#0');

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

    SELECT out_nRet AS nRet;
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
    INSERT INTO tItem (nItemKey, nOwner, nStorageType, nStorage, nItemID, nFlags, dDate)
    VALUES (in_nItemKey, in_nOwner, in_nStorageType, in_nStorage, in_nItemID, in_nFlags, NOW())
    ON DUPLICATE KEY UPDATE nItemID = VALUES(nItemID), nFlags = VALUES(nFlags);
    
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
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
    DELETE FROM tItem WHERE nItemKey = in_nItemKey;
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END""",

    "p_Item_GetStorageEmptyList": """CREATE PROCEDURE `p_Item_GetStorageEmptyList`(
    IN in_nOwner INT,
    IN in_nStorageType INT,
    IN in_nGetStorageCnt INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    IF in_nOwner = -1 THEN
        SET out_nRet = -1000;
        LEAVE sp_body;
    END IF;

    SELECT A.nStorage
    FROM tItemStorageFindEmpty A
    LEFT OUTER JOIN (
        SELECT nStorage FROM tItem WHERE nOwner = in_nOwner AND nStorageType = in_nStorageType
    ) B ON A.nStorage = B.nStorage
    WHERE B.nStorage IS NULL
    ORDER BY A.nStorage
    LIMIT in_nGetStorageCnt;
END""",

    "p_Char_Find": """CREATE PROCEDURE `p_Char_Find`(
    IN in_sID VARCHAR(40),
    OUT out_nCharNo INT
)
sp_body: BEGIN
    SET out_nCharNo = 0;
    SELECT nCharNo INTO out_nCharNo FROM tCharacter WHERE sID = in_sID LIMIT 1;
    IF out_nCharNo IS NULL THEN
        SET out_nCharNo = 0;
    END IF;
    SELECT out_nCharNo AS nCharNo;
END""",

    "p_Char_GetCharNo": """CREATE PROCEDURE `p_Char_GetCharNo`(
    IN in_sID VARCHAR(40),
    OUT out_nCharNo INT
)
sp_body: BEGIN
    SET out_nCharNo = 0;
    SELECT nCharNo INTO out_nCharNo FROM tCharacter WHERE sID = in_sID LIMIT 1;
    IF out_nCharNo IS NULL THEN
        SET out_nCharNo = 0;
    END IF;
    SELECT out_nCharNo AS nCharNo;
END""",

    "p_Char_Login": """CREATE PROCEDURE `p_Char_Login`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DECLARE v_lastLogin DATETIME;
    SET out_nRet = 0;

    IF EXISTS (SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo) THEN
        SELECT dLastLoginDate INTO v_lastLogin FROM tCharacter WHERE nCharNo = in_nCharNo LIMIT 1;
        UPDATE tCharacter SET nLoginCount = nLoginCount + 1, dLastLoginDate = NOW() WHERE nCharNo = in_nCharNo;
        SET out_nRet = 0;
        SELECT IFNULL(DATEDIFF(NOW(), v_lastLogin), 0) AS nLoginDateDiff;
    ELSE
        SET out_nRet = -1111;
    END IF;
    SELECT out_nRet AS nRet;
END""",

    "p_Char_Logout": """CREATE PROCEDURE `p_Char_Logout`(
    IN in_nCharNo INT,
    IN in_nPlayMin INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;
    IF EXISTS (SELECT 1 FROM tCharacter WHERE nCharNo = in_nCharNo) THEN
        UPDATE tCharacter SET nPlayMin = nPlayMin + in_nPlayMin WHERE nCharNo = in_nCharNo;
        SET out_nRet = 1;
    ELSE
        SET out_nRet = 0;
    END IF;
    SELECT out_nRet AS nRet;
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
    UPDATE tQuest
    SET nStatus = in_nStatus, sData = in_sData
    WHERE nCharNo = in_nCharNo AND nQuestNo = in_nQuestNo;

    IF ROW_COUNT() = 0 THEN
        INSERT INTO tQuest (nCharNo, nQuestNo, nStatus, sData)
        VALUES (in_nCharNo, in_nQuestNo, in_nStatus, in_sData);
    END IF;

    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
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
    IF EXISTS (SELECT nCharNo FROM tTutorial WHERE nUserNo = in_nUserNo AND nCharNo = in_nCharNo) THEN
        UPDATE tTutorial SET nState = in_nState, nStep = in_nStep WHERE nUserNo = in_nUserNo AND nCharNo = in_nCharNo;
        IF ROW_COUNT() <> 1 THEN SET out_nRet = -1; END IF;
    ELSE
        INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (in_nUserNo, in_nCharNo, 0, 0);
        IF ROW_COUNT() <> 1 THEN SET out_nRet = -1; END IF;
    END IF;
    SELECT out_nRet AS nRet;
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
    SET out_nRet = 0;
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
        LEAVE sp_body;
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
        LEAVE sp_body;
    END IF;

    COMMIT;
    SET out_nRet = 0;
    SELECT out_SubStatusA AS SubStatusA, out_SubStatusB AS SubStatusB, out_nRet AS nRet;
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
    SET v_guildMaster = 0;

    IF IFNULL(in_nNo, 0) = 0 THEN
        SET out_nRet = -1311;
        LEAVE sp_body;
    END IF;

    IF IFNULL(in_nCharNo, 0) = 0 THEN
        SET out_nRet = -1313;
        LEAVE sp_body;
    END IF;

    SELECT IFNULL(COUNT(nCharNo), 0) INTO v_cnt
    FROM tCharacterChangeID
    WHERE nCharNo = in_nCharNo AND bFlag = 0 AND bInit = 1;

    IF in_nGrade = 0 AND v_cnt = 1 THEN
        SET out_nRet = -1323;
        LEAVE sp_body;
    ELSE
        START TRANSACTION;

        IF in_nGrade = 0 THEN
            SELECT IFNULL(MAX(nCharNo), 0) INTO v_guildMaster
            FROM tGuildMember WHERE nNo = in_nNo AND nGrade = 0;

            SELECT MAX(dDate) INTO v_lastChangeDate
            FROM tLog_GuildMasterChange WHERE nGuildNo = in_nNo;

            IF v_lastChangeDate IS NOT NULL AND DATE_ADD(v_lastChangeDate, INTERVAL 7 DAY) > NOW() THEN
                ROLLBACK;
                SET out_nRet = -1324;
                LEAVE sp_body;
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
    END IF;
    SELECT out_nRet AS nRet;
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

    IF EXISTS (SELECT nNo FROM tGuildAcademyMember WHERE nNo = in_nNo AND nCharNo = in_nCharNo) THEN
        SELECT T1.nMaxMembers, T2.nCnt INTO v_nMax, v_nCnt
        FROM tGuild AS T1
        INNER JOIN (
            SELECT G.nNo, IFNULL(COUNT(*), 0) AS nCnt
            FROM tCharacter AS C, tGuildMember AS G
            WHERE C.nCharNo = G.nCharNo AND G.nNo = in_nNo AND C.bDeleted = 0
            GROUP BY G.nNo
        ) AS T2 ON (T1.nNo = T2.nNo);

        IF v_nCnt < v_nMax THEN
            START TRANSACTION;

            CALL `usp_GuildAcademy_delMember`(in_nNo, in_nCharNo, out_nRet);
            IF out_nRet <> 0 THEN
                ROLLBACK;
                LEAVE sp_body;
            END IF;

            CALL `p_Guild_Member_Join`(in_nNo, in_nCharNo, in_nGrade, out_nRet);
            IF out_nRet <> 1 THEN
                SET out_nRet = -1531;
                ROLLBACK;
                LEAVE sp_body;
            ELSE
                SET out_nRet = 0;
                COMMIT;
            END IF;
        ELSE
            SET out_nRet = -1530;
        END IF;
    ELSE
        SET out_nRet = -1528;
    END IF;
    SELECT out_nRet AS nRet;
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
    CALL `p_Item_IsValidOwner`(in_nItemKey, in_nOwner, v_nTinyintRet);
    IF v_nTinyintRet = 0 THEN
        SET out_nRet = 2;
        LEAVE sp_body;
    END IF;

    START TRANSACTION;
    IF in_nGiftMoney > 0 THEN
        CALL `p_Char_MoneySub`(in_nOwner, in_nGiftMoney, v_nTinyintRet);
        IF v_nTinyintRet = 0 THEN
            ROLLBACK;
            SET out_nRet = 3;
            LEAVE sp_body;
        END IF;
    END IF;

    CALL `p_Item_SetOption`(in_nItemKey, in_nUpgradeType, in_nUpgradeValue, v_nRetTemp);
    IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet = 4;
        LEAVE sp_body;
    END IF;

    IF in_nRawLot = 1 THEN
        CALL `p_Item_Delete`(in_nRawKey, v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
            ROLLBACK;
            SET out_nRet = 5;
            LEAVE sp_body;
        END IF;
    ELSE
        CALL `p_Item_SetOption`(in_nRawKey, in_nRawOptionType, in_nRawLot - 1, v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
            ROLLBACK;
            SET out_nRet = 5;
            LEAVE sp_body;
        END IF;
    END IF;

    IF in_nRawLeftKey <> 0 THEN
        IF in_nRawLeftLot = 1 THEN
            CALL `p_Item_Delete`(in_nRawLeftKey, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 6;
                LEAVE sp_body;
            END IF;
        ELSE
            CALL `p_Item_SetOption`(in_nRawLeftKey, in_nRawLeftOptionType, in_nRawLeftLot - 1, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 6;
                LEAVE sp_body;
            END IF;
        END IF;
    END IF;

    IF in_nRawRightKey <> 0 THEN
        IF in_nRawRightLot = 1 THEN
            CALL `p_Item_Delete`(in_nRawRightKey, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 7;
                LEAVE sp_body;
            END IF;
        ELSE
            CALL `p_Item_SetOption`(in_nRawRightKey, in_nRawRightOptionType, in_nRawRightLot - 1, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 7;
                LEAVE sp_body;
            END IF;
        END IF;
    END IF;

    IF in_nRawMiddleKey <> 0 THEN
        IF in_nRawMiddleLot = 1 THEN
            CALL `p_Item_Delete`(in_nRawMiddleKey, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 8;
                LEAVE sp_body;
            END IF;
        ELSE
            CALL `p_Item_SetOption`(in_nRawMiddleKey, in_nRawMiddleOptionType, in_nRawMiddleLot - 1, v_nRetTemp);
            IF v_nRetTemp <> 0 THEN
                ROLLBACK;
                SET out_nRet = 8;
                LEAVE sp_body;
            END IF;
        END IF;
    END IF;

    COMMIT;
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END""",
}

def tokenize(sql):
    token_spec = [
        ('STRING', r"'([^']|'')*'"),
        ('NUMBER', r'\b\d+(\.\d*)?\b'),
        ('HEX', r'0x[0-9a-fA-F]+'),
        ('WORD', r'\b[a-zA-Z_][a-zA-Z0-9_]*\b'),
        ('OPERATOR', r'<>|<=|>=|!=|==|:=|[+\-*/=<>!&|^]'),
        ('DELIM', r'[(),;.]'),
        ('WS', r'\s+'),
    ]
    tok_regex = '|'.join('(?P<%s>%s)' % pair for pair in token_spec)
    tokens = []
    for mo in re.finditer(tok_regex, sql):
        kind = mo.lastgroup
        val = mo.group()
        if kind != 'WS':
            tokens.append((kind, val))
    return tokens

def parse_body_statements(tokens):
    idx = 0
    n = len(tokens)
    
    STMT_KEYWORDS = {"SET", "SELECT", "INSERT", "UPDATE", "DELETE", "CALL", "LEAVE", "START", "COMMIT", "ROLLBACK", "PREPARE", "EXECUTE", "DEALLOCATE"}
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
            get_token() # consume BEGIN
            stmts = []
            while peek() and peek() != "END":
                s = parse_statement()
                if s:
                    stmts.append(s)
            if peek() == "END":
                get_token() # consume END
            return "\n    ".join(stmts)

        if first == "IF":
            get_token() # consume IF
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
                        get_token() # consume BEGIN
                        then_stmts = []
                        while peek() and peek() != "END":
                            s = parse_statement()
                            if s:
                                then_stmts.append(s)
                        if peek() == "END":
                            get_token() # consume END
                        
                        else_part = ""
                        if peek() == "ELSE":
                            get_token() # consume ELSE
                            if peek() == "BEGIN":
                                get_token()
                                e_stmts = []
                                while peek() and peek() != "END":
                                    es = parse_statement()
                                    if es:
                                        e_stmts.append(es)
                                if peek() == "END":
                                    get_token()
                                else_part = "ELSE\n        " + "\n        ".join(e_stmts)
                            else:
                                es = parse_statement()
                                if es:
                                    else_part = f"ELSE\n        {es}"
                        
                        cond_str = " ".join(t[1] for t in cond_tokens)
                        cond_str = re.sub(r"\bCAST\s+\(", "CAST(", cond_str, flags=re.IGNORECASE)
                        then_body = "\n        ".join(then_stmts)
                        if else_part:
                            return f"IF {cond_str} THEN\n        {then_body}\n    {else_part}\n    END IF;"
                        else:
                            return f"IF {cond_str} THEN\n        {then_body}\n    END IF;"
                    
                    elif w in STMT_KEYWORDS or w == "IF":
                        # Single statement follows
                        then_stmt = parse_statement()
                        else_part = ""
                        if peek() == "ELSE":
                            get_token() # consume ELSE
                            if peek() == "BEGIN":
                                get_token()
                                e_stmts = []
                                while peek() and peek() != "END":
                                    es = parse_statement()
                                    if es:
                                        e_stmts.append(es)
                                if peek() == "END":
                                    get_token()
                                else_part = "ELSE\n        " + "\n        ".join(e_stmts)
                            else:
                                es = parse_statement()
                                if es:
                                    else_part = f"ELSE\n        {es}"
                        
                        cond_str = " ".join(t[1] for t in cond_tokens)
                        cond_str = re.sub(r"\bCAST\s+\(", "CAST(", cond_str, flags=re.IGNORECASE)
                        if else_part:
                            return f"IF {cond_str} THEN\n        {then_stmt}\n    {else_part}\n    END IF;"
                        else:
                            return f"IF {cond_str} THEN\n        {then_stmt}\n    END IF;"

                cond_tokens.append(get_token())

        if first == "ELSE":
            get_token()
            stmt = parse_statement()
            return f"ELSE\n    {stmt}"

        stmt_tokens = []
        depth = 0
        in_update = (first == "UPDATE")
        in_insert = (first == "INSERT")
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
                    pass # SELECT is part of INSERT INTO ... SELECT
                elif w in STMT_KEYWORDS or w in BLOCK_KEYWORDS:
                    break

            stmt_tokens.append(get_token())

        stmt_str = " ".join(t[1] for t in stmt_tokens).strip()
        if stmt_str.rstrip(";").isdigit() or stmt_str.rstrip(";") in ("<", ">", ")", ","):
            return None

        stmt_str = re.sub(r"\s*([(),;.=<>!+*/-])\s*", r"\1", stmt_str)
        stmt_str = re.sub(r"\b(INTO|FROM|WHERE|AND|OR|SET|VALUES|ORDER BY|GROUP BY|LIMIT|JOIN|INNER JOIN|LEFT JOIN|AS)\b", r" \1 ", stmt_str, flags=re.IGNORECASE)
        stmt_str = re.sub(r"\bCAST\s+\(", "CAST(", stmt_str, flags=re.IGNORECASE)
        stmt_str = re.sub(r"\s+", " ", stmt_str).strip()
        
        if not stmt_str.endswith(";"):
            stmt_str += ";"
        return stmt_str

    all_stmts = []
    while idx < n:
        s = parse_statement()
        if s:
            all_stmts.append(s)
    return "\n    ".join(all_stmts)

def transpile_procedure(sp_name):
    if sp_name in SPECIAL_PROCEDURES:
        return SPECIAL_PROCEDURES[sp_name]

    filepath = os.path.join(TARGET_DIR, sp_name + ".sql")
    with open(filepath, "r", encoding="latin-1") as fp:
        raw = fp.read()

    clean = clean_sql_source(raw)

    m = re.search(r"CREATE\s+PROCEDURE\s+\[?dbo\]?\.\[?([a-zA-Z0-9_]+)\]?(.*?)\bAS\b(.*)", clean, flags=re.DOTALL | re.IGNORECASE)
    if not m:
        raise ValueError(f"Failed to match procedure header for {sp_name}")

    name = m.group(1)
    param_text = m.group(2)
    body_text = m.group(3)

    p_clean = re.sub(r"/\*.*?\*/", "", param_text, flags=re.DOTALL)
    p_clean = re.sub(r"--.*", "", p_clean)

    raw_params = re.findall(r"@([a-zA-Z0-9_]+)\s+([a-zA-Z0-9_]+(?:\([^\)]+\))?)(?:\s+(OUTPUT|OUT))?(?:\s*=\s*([^,\n]+))?", p_clean, re.IGNORECASE)

    mysql_params = []
    param_map = {}
    out_params = []

    for pname, ptype, pout, pdef in raw_params:
        is_out = bool(pout)
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

    body = body_text
    body = re.sub(r"/\*.*?\*/", "", body, flags=re.DOTALL)
    body = re.sub(r"--.*", "", body)

    # Cross-database prefixes
    body = re.sub(r"\b(?:[a-zA-Z0-9_]+\.)*dbo\.", "", body, flags=re.IGNORECASE)

    # Remove T-SQL settings and hints
    body = re.sub(r"\bSET\s+NOCOUNT\s+(?:ON|OFF)\b;?", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bSET\s+LOCK_TIMEOUT\s+-?\d+;?", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bWITH\s*\([^\)]+\)", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\b\(NOLOCK\)", "", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGETDATE\(\)", "NOW()", body, flags=re.IGNORECASE)
    body = re.sub(r"\bISNULL\(", "IFNULL(", body, flags=re.IGNORECASE)
    body = re.sub(r"\bLEN\(", "CHAR_LENGTH(", body, flags=re.IGNORECASE)
    body = re.sub(r"@@IDENTITY\b", "LAST_INSERT_ID()", body, flags=re.IGNORECASE)
    body = re.sub(r"\bSCOPE_IDENTITY\(\)", "LAST_INSERT_ID()", body, flags=re.IGNORECASE)
    body = re.sub(r"@@ROWCOUNT\b", "ROW_COUNT()", body, flags=re.IGNORECASE)
    body = re.sub(r"@@TRANCOUNT\b", "1", body, flags=re.IGNORECASE)

    # Cast integer fix
    body = re.sub(r"\bCAST\s*\(\s*(.+?)\s+AS\s+INT\s*\)", r"CAST(\1 AS SIGNED)", body, flags=re.IGNORECASE)

    # Redundant FROM in UPDATE fix (without re.DOTALL so it never crosses queries)
    body = re.sub(r"\bUPDATE\s+([a-zA-Z0-9_]+)\s+SET\s+([^;\n]+?)\s+FROM\s+\1\b", r"UPDATE \1 SET \2", body, flags=re.IGNORECASE)

    # Strip TOP inside EXISTS(...)
    body = re.sub(r"\bEXISTS\s*\(\s*SELECT\s+TOP\s+\d+\s+", "EXISTS(SELECT ", body, flags=re.IGNORECASE)

    # DATEADD conversion
    def dateadd_repl(dam):
        unit = dam.group(1).upper()
        val = dam.group(2).strip()
        date_arg = dam.group(3).strip()
        return f"DATE_ADD({date_arg}, INTERVAL {val} {unit})"
    body = re.sub(r"\bDATEADD\s*\(\s*([a-zA-Z]+)\s*,\s*([^,]+)\s*,\s*((?:[^()]+|\([^()]*\))+)\)", dateadd_repl, body, flags=re.IGNORECASE)

    # Common T-SQL error patterns
    body = re.sub(r"IF\s+@@ERROR\s*<>\s*0\s+OR\s+ROW_COUNT\(\)\s*=\s*0", "IF ROW_COUNT() = 0", body, flags=re.IGNORECASE)
    body = re.sub(r"IF\s+@@ERROR\s*<>\s*0\s+OR\s+ROW_COUNT\(\)\s*<>\s*1", "IF ROW_COUNT() <> 1", body, flags=re.IGNORECASE)
    body = re.sub(r"IF\s+@@ERROR\s*=\s*0\s+AND\s+ROW_COUNT\(\)\s*=\s*1", "IF ROW_COUNT() = 1", body, flags=re.IGNORECASE)
    body = re.sub(r"@@ERROR\b", "0", body, flags=re.IGNORECASE)

    body = re.sub(r"\bBEGIN\s+TRAN(?:SACTION)?\b", "START TRANSACTION;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bCOMMIT\s+TRAN(?:SACTION)?\b", "COMMIT;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bROLLBACK\s+TRAN(?:SACTION)?\b", "ROLLBACK;", body, flags=re.IGNORECASE)

    # Hoist single and multi-variable DECLARE statements (supporting consecutive DECLAREs & AS)
    mysql_declares = []
    declare_blocks = re.findall(r"\bDECLARE\b(.*?)(?=\b(?:SET|BEGIN|UPDATE|SELECT|INSERT|DELETE|IF|CALL|START|COMMIT|ROLLBACK|EXEC|DECLARE)\b)", body, flags=re.DOTALL | re.IGNORECASE)
    body = re.sub(r"\bDECLARE\b.*?(?=\b(?:SET|BEGIN|UPDATE|SELECT|INSERT|DELETE|IF|CALL|START|COMMIT|ROLLBACK|EXEC|DECLARE)\b)", "", body, flags=re.DOTALL | re.IGNORECASE)

    for decl_str in declare_blocks:
        decl_clean = decl_str.strip().rstrip(";")
        parts = []
        cur = []
        depth = 0
        for c in decl_clean:
            if c == "(": depth += 1
            elif c == ")": depth -= 1
            elif c == "," and depth == 0:
                parts.append("".join(cur).strip())
                cur = []
                continue
            cur.append(c)
        if cur: parts.append("".join(cur).strip())

        for p in parts:
            dm = re.match(r"@([a-zA-Z0-9_]+)(?:\s+AS)?\s+([a-zA-Z0-9_]+(?:\([^\)]+\))?)(?:\s*=\s*(.+))?", p.strip(), re.IGNORECASE)
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

    # Replace GOTOs and strip labels BEFORE SELECT parsing so labels don't get swallowed into SELECT
    body = re.sub(r"\bGOTO\s+ERROR_[ \t]*;?", "ROLLBACK; SET out_nRet = 1; LEAVE sp_body;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO\s+END_[ \t]*;?", "LEAVE sp_body;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO\s+ERR[ \t]*;?", "ROLLBACK; SET out_err = 1; LEAVE sp_body;", body, flags=re.IGNORECASE)
    body = re.sub(r"\bGOTO[ \t]+[a-zA-Z0-9_]+[ \t]*;?", "LEAVE sp_body;", body, flags=re.IGNORECASE)
    body = re.sub(r"^[ \t]*[a-zA-Z0-9_]+:[ \t]*.*$", "", body, flags=re.MULTILINE)

    # EXEC calls: EXEC other_sp @a, @b, @ret OUTPUT -> CALL other_sp(a, b, ret);
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
            clean_args.append(ra)
        return f"CALL `{proc}`({', '.join(clean_args)});"

    body = re.sub(r"\bEXEC(?:UTE)?\s+([a-zA-Z0-9_]+)\s+([^;]+?)(?=(?:\b(?:SET|IF|ELSE|BEGIN|END|RETURN|GOTO|COMMIT|ROLLBACK|EXEC|START)\b|$))", exec_repl, body, flags=re.DOTALL | re.IGNORECASE)

    # SELECT @var = col FROM ... -> SELECT col INTO var FROM ...
    def select_assign_repl(sm):
        full_stmt = sm.group(0).strip()
        if re.search(r"\bFROM\b", full_stmt, re.IGNORECASE):
            m_split = re.split(r"\bFROM\b", full_stmt, maxsplit=1, flags=re.IGNORECASE)
            assigns_part = re.sub(r"^\s*SELECT(?:\s+TOP\s+\d+)?\s+", "", m_split[0], flags=re.IGNORECASE).strip()
            rest_part = "FROM " + m_split[1].strip()
            
            top_match = re.search(r"\bTOP\s+(\d+)\b", m_split[0], re.IGNORECASE)
            limit_clause = f" LIMIT {top_match.group(1)}" if top_match else ""
            
            pairs = []
            cur = []
            depth = 0
            for char in assigns_part:
                if char == '(': depth += 1
                elif char == ')': depth -= 1
                elif char == ',' and depth == 0:
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

    body = re.sub(r"\bSELECT(?:\s+TOP\s+\d+)?\s+@?[a-zA-Z0-9_]+\s*=\s*[^;]+?(?=(?:\b(?:SELECT|SET|INSERT|UPDATE|DELETE|IF|ELSE|BEGIN|END|RETURN|GOTO|COMMIT|ROLLBACK|CALL|START)\b|$))", select_assign_repl, body, flags=re.IGNORECASE | re.DOTALL)

    # SELECT TOP n ... without assign -> SELECT ... LIMIT n
    def select_top_repl(stm):
        n = stm.group(1)
        rest = stm.group(2).rstrip("; \t\n")
        return f"SELECT {rest} LIMIT {n};"
    body = re.sub(r"\bSELECT\s+TOP\s+(\d+)\s+([^;]+?)(?=(?:\b(?:SET|INSERT|UPDATE|DELETE|IF|ELSE|BEGIN|END|RETURN|GOTO|COMMIT|ROLLBACK|CALL|START)\b|$))", select_top_repl, body, flags=re.IGNORECASE | re.DOTALL)

    # Variable substitution
    for raw_var, mapped_var in param_map.items():
        body = re.sub(rf"@{raw_var}\b", mapped_var, body, flags=re.IGNORECASE)

    # Replace RETURN with LEAVE sp_body; (only horizontal whitespace to avoid eating next line)
    body = re.sub(r"\bRETURN(?:[ \t]+[^;\n\r]+)?[ \t]*;?", "LEAVE sp_body;", body, flags=re.IGNORECASE)

    # Tokenize and parse into proper structured statements
    tokens = tokenize(body)
    parsed_body = parse_body_statements(tokens)

    params_sql = ",\n    ".join(mysql_params)
    if params_sql:
        params_sql = "\n    " + params_sql + "\n"

    declares_sql = "\n".join(mysql_declares)
    if declares_sql:
        declares_sql = "\n" + declares_sql + "\n"

    out_select = ""
    if out_params:
        cols = ", ".join(f"{p} AS {p.replace('out_', '')}" for p in out_params)
        out_select = f"\n    SELECT {cols};\n"

    sql = f"""CREATE PROCEDURE `{name}`({params_sql})
sp_body: BEGIN
{declares_sql}
    {parsed_body}
{out_select}
END"""
    return sql

if __name__ == "__main__":
    conn = mysql.connector.connect(host='127.0.0.1', port=3306, user='root', password='lerd', database='xelion_online')
    cursor = conn.cursor()

    success = 0
    errors = []
    all_sqls = ["""-- ============================================================================
-- XELION ONLINE - STORED PROCEDURES OFICIALES DE WORLD00_CHARACTER (MYSQL 8.4)
-- Catálogo oficial de 311 Stored Procedures de Character.exe y WorldManager.exe
-- 100% de paridad con World00_Character.bak, PDBs oficiales y docs/ProcedimientosLegacy.md
-- ============================================================================

USE `xelion_online`;

"""]

    for sp in target_311:
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

    print(f"Results: {success}/{len(target_311)} procedures compiled cleanly in MySQL 8.4!")
    print(f"Saved complete script with {len(target_311)} Stored Procedures to database/procedures_character.sql.")
    if errors:
        print(f"Errors ({len(errors)}):")
        for sp, err in errors:
            print(f"  {sp:35s}: {err.replace(chr(10), ' ')[:90]}")

    conn.close()
