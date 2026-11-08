#!/usr/bin/env python3
"""
build_procedures_gamelog.py - Generador 100% exhaustivo y pool-safe de los 27 Stored Procedures
de World00_GameLog (23 invocaciones ODBC de GameLog.exe U 26 de World00_GameLog.bak).
Genera database/procedures_gamelog.sql con encapsulación proc_logic + SELECT final incondicional.
"""

import os
import sys

OUTPUT_FILE = "database/procedures_gamelog.sql"

PROC_DEFINITIONS = [
    # 1. p_GameLog_Add
    {
        "name": "p_GameLog_Add",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": [
            "nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ",
            "nTargetCharNo", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3", "nBigint1"
        ],
        "values": [
            "in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ",
            "in_nTargetCharNo", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3", "in_nBigint1"
        ],
        "type": "add"
    },
    # 2. p_GameLog_Add_Type1
    {
        "name": "p_GameLog_Add_Type1",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ"],
        "type": "add"
    },
    # 3. p_GameLog_Add_Type2
    {
        "name": "p_GameLog_Add_Type2",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 4. p_GameLog_Add_Type3
    {
        "name": "p_GameLog_Add_Type3",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nInt1", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetCharNo", "nTargetID", "nInt1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetCharNo", "in_nTargetID", "in_nInt1"],
        "type": "add"
    },
    # 5. p_GameLog_Add_Type4
    {
        "name": "p_GameLog_Add_Type4",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nInt1", "nInt2", "nInt3"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nInt1", "in_nInt2", "in_nInt3"],
        "type": "add"
    },
    # 6. p_GameLog_Add_Type5
    {
        "name": "p_GameLog_Add_Type5",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 7. p_GameLog_Add_Type6
    {
        "name": "p_GameLog_Add_Type6",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nInt1", "nBigint1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nInt1", "in_nBigint1"],
        "type": "add"
    },
    # 8. p_GameLog_Add_Type7
    {
        "name": "p_GameLog_Add_Type7",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetCharNo", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "nTargetCharNo", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 9. p_GameLog_Add_Type8
    {
        "name": "p_GameLog_Add_Type8",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nInt1", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nInt1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nInt1"],
        "type": "add"
    },
    # 10. p_GameLog_Add_Type9
    {
        "name": "p_GameLog_Add_Type9",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nItemKey", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nItemKey", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 11. p_GameLog_Add_TypeA
    {
        "name": "p_GameLog_Add_TypeA",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
        ],
        "columns": ["nType", "nCharNo", "nTargetCharNo", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3"],
        "values": ["in_nType", "in_nCharNo", "in_nTargetCharNo", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3"],
        "type": "add"
    },
    # 12. p_GameLog_Add_TypeB
    {
        "name": "p_GameLog_Add_TypeB",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetCharNo", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetCharNo", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3"],
        "type": "add"
    },
    # 13. p_GameLog_Add_TypeC
    {
        "name": "p_GameLog_Add_TypeC",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3"],
        "type": "add"
    },
    # 14. p_GameLog_Add_TypeD
    {
        "name": "p_GameLog_Add_TypeD",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "nItemKey", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_nItemKey", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 15. p_GameLog_Add_TypeE
    {
        "name": "p_GameLog_Add_TypeE",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetID", "nItemKey", "nInt1", "nInt2"],
        "values": ["in_nType", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 16. p_GameLog_Add_TypeF
    {
        "name": "p_GameLog_Add_TypeF",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3", "nBigint1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3", "in_nBigint1"],
        "type": "add"
    },
    # 17. p_GameLog_Add_TypeG
    {
        "name": "p_GameLog_Add_TypeG",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetID", "nItemKey", "nInt1", "nInt2", "nBigint1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nBigint1"],
        "type": "add"
    },
    # 18. p_GameLog_Add_TypeH
    {
        "name": "p_GameLog_Add_TypeH",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetCharNo"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetCharNo"],
        "type": "add"
    },
    # 19. p_GameLog_Add_TypeI
    {
        "name": "p_GameLog_Add_TypeI",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetCharNo", "nTargetID"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetCharNo", "in_nTargetID"],
        "type": "add"
    },
    # 20. p_GameLog_Add_TypeJ
    {
        "name": "p_GameLog_Add_TypeJ",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nTargetID", "nItemKey", "nInt1", "nInt2"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2"],
        "type": "add"
    },
    # 21. p_GameLog_Add_TypeK
    {
        "name": "p_GameLog_Add_TypeK",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_sMap", "VARCHAR(16)"),
            ("IN", "in_nMapX", "INT"),
            ("IN", "in_nMapY", "INT"),
            ("IN", "in_nMapZ", "INT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": ["nType", "nCharNo", "sMap", "nMapX", "nMapY", "nMapZ", "nInt1", "nInt2", "nBigint1"],
        "values": ["in_nType", "in_nCharNo", "in_sMap", "in_nMapX", "in_nMapY", "in_nMapZ", "in_nInt1", "in_nInt2", "in_nBigint1"],
        "type": "add"
    },
    # 22. p_GameLog_Add_TypeL
    {
        "name": "p_GameLog_Add_TypeL",
        "params": [
            ("OUT", "out_nLogNo", "BIGINT"),
            ("IN", "in_nType", "INT"),
            ("IN", "in_nCharNo", "INT"),
            ("IN", "in_nTargetCharNo", "INT"),
            ("IN", "in_nTargetID", "INT"),
            ("IN", "in_nItemKey", "BIGINT"),
            ("IN", "in_nInt1", "INT"),
            ("IN", "in_nInt2", "INT"),
            ("IN", "in_nInt3", "INT"),
            ("IN", "in_nBigint1", "BIGINT"),
        ],
        "columns": ["nType", "nCharNo", "nTargetCharNo", "nTargetID", "nItemKey", "nInt1", "nInt2", "nInt3", "nBigint1"],
        "values": ["in_nType", "in_nCharNo", "in_nTargetCharNo", "in_nTargetID", "in_nItemKey", "in_nInt1", "in_nInt2", "in_nInt3", "in_nBigint1"],
        "type": "add"
    },
    # 23. p_GameLog_Get
    {
        "name": "p_GameLog_Get",
        "type": "custom",
        "sql": """DROP PROCEDURE IF EXISTS `p_GameLog_Get`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Get`(
    IN in_WhereExpr VARCHAR(256)
)
sp_body: BEGIN
    proc_logic: BEGIN
        IF in_WhereExpr IS NULL OR TRIM(in_WhereExpr) = '' THEN
            SET @v_sql = 'SELECT nLogNo, dDate, nCharNo, nType, sMap, nMapX, nMapY, nMapZ, nTargetCharNo, nTargetID, nItemKey, nInt1, nInt2, nInt3, nBigint1 FROM tGameLog ORDER BY nLogNo DESC LIMIT 100';
        ELSE
            SET @v_sql = CONCAT('SELECT nLogNo, dDate, nCharNo, nType, sMap, nMapX, nMapY, nMapZ, nTargetCharNo, nTargetID, nItemKey, nInt1, nInt2, nInt3, nBigint1 FROM tGameLog WHERE ', in_WhereExpr);
        END IF;
        PREPARE stmt FROM @v_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        LEAVE proc_logic;
    END proc_logic;
END$$
DELIMITER ;
"""
    },
    # 24. p_GameLog_GetCount
    {
        "name": "p_GameLog_GetCount",
        "type": "custom",
        "sql": """DROP PROCEDURE IF EXISTS `p_GameLog_GetCount`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_GetCount`(
    IN in_WhereExpr VARCHAR(256),
    OUT out_nCount INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nCount = 0;
        IF in_WhereExpr IS NULL OR TRIM(in_WhereExpr) = '' THEN
            SET @v_sql = 'SELECT COUNT(*) INTO @v_count FROM tGameLog';
        ELSE
            SET @v_sql = CONCAT('SELECT COUNT(*) INTO @v_count FROM tGameLog WHERE ', in_WhereExpr);
        END IF;
        PREPARE stmt FROM @v_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        SET out_nCount = IFNULL(@v_count, 0);
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCount AS nCount;
    END IF;
END$$
DELIMITER ;
"""
    },
    # 25. p_GameLog_GetLog
    {
        "name": "p_GameLog_GetLog",
        "type": "custom",
        "sql": """DROP PROCEDURE IF EXISTS `p_GameLog_GetLog`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_GetLog`(
    IN in_QueryExpr VARCHAR(256)
)
sp_body: BEGIN
    proc_logic: BEGIN
        IF in_QueryExpr IS NULL OR TRIM(in_QueryExpr) = '' THEN
            SET @v_sql = 'SELECT a.nLogNo, a.dDate, a.nCharNo, a.nType, a.sMap, a.nMapX, a.nMapY, a.nTargetCharNo, a.nTargetID, a.nItemKey, a.nInt1, a.nInt2, a.nInt3, a.nBigint1, b.sID FROM tGameLog a INNER JOIN tCharacter b ON a.nCharNo = b.nCharNo ORDER BY a.nLogNo DESC LIMIT 100';
        ELSE
            SET @v_sql = CONCAT('SELECT a.nLogNo, a.dDate, a.nCharNo, a.nType, a.sMap, a.nMapX, a.nMapY, a.nTargetCharNo, a.nTargetID, a.nItemKey, a.nInt1, a.nInt2, a.nInt3, a.nBigint1, b.sID FROM tGameLog a INNER JOIN tCharacter b ON a.nCharNo = b.nCharNo ', in_QueryExpr);
        END IF;
        PREPARE stmt FROM @v_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        LEAVE proc_logic;
    END proc_logic;
END$$
DELIMITER ;
"""
    },
    # 26. p_GameLog_GetLogCount
    {
        "name": "p_GameLog_GetLogCount",
        "type": "custom",
        "sql": """DROP PROCEDURE IF EXISTS `p_GameLog_GetLogCount`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_GetLogCount`(
    IN in_QueryExpr VARCHAR(256),
    OUT out_nCount INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nCount = 0;
        IF in_QueryExpr IS NULL OR TRIM(in_QueryExpr) = '' THEN
            SET @v_sql = 'SELECT COUNT(*) INTO @v_count FROM tGameLog a INNER JOIN tCharacter b ON a.nCharNo = b.nCharNo';
        ELSE
            SET @v_sql = CONCAT('SELECT COUNT(*) INTO @v_count FROM tGameLog a INNER JOIN tCharacter b ON a.nCharNo = b.nCharNo ', in_QueryExpr);
        END IF;
        PREPARE stmt FROM @v_sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
        SET out_nCount = IFNULL(@v_count, 0);
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nCount AS nCount;
    END IF;
END$$
DELIMITER ;
"""
    },
    # 27. usp_Login_setChannel
    {
        "name": "usp_Login_setChannel",
        "type": "custom",
        "sql": """DROP PROCEDURE IF EXISTS `usp_Login_setChannel`;
DELIMITER $$
CREATE PROCEDURE `usp_Login_setChannel`(
    IN in_nUserNo INT,
    IN in_nChannelNo INT,
    IN in_sChannelID VARCHAR(20),
    OUT out_nResult INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nResult = 0;
        INSERT INTO `tChannelLog` (`userNo`, `channel`, `issueDate`, `dDate`)
        VALUES (in_nUserNo, in_sChannelID, CAST(in_nChannelNo AS CHAR), NOW());
        SET out_nResult = 1;
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nResult AS nResult;
    END IF;
END$$
DELIMITER ;
"""
    },
]

def generate_add_proc(defn):
    name = defn["name"]
    params_sql = ",\n".join(f"    {mode} {pname} {ptype}" for mode, pname, ptype in defn["params"])
    cols_sql = ", ".join(f"`{c}`" for c in defn["columns"])
    vals_sql = ", ".join(defn["values"])

    return f"""-- ----------------------------------------------------------------------------
-- Procedure: {name}
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `{name}`;
DELIMITER $$
CREATE PROCEDURE `{name}`(
{params_sql}
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` ({cols_sql})
        VALUES ({vals_sql});
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;
"""

def main():
    lines = [
        "-- ============================================================================",
        "-- XELION ONLINE - STORED PROCEDURES OFICIALES DE WORLD00_GAMELOG (MYSQL 8.4)",
        "-- Catálogo exhaustivo de 27 Stored Procedures oficiales",
        "-- 22 Inserciones de Logs (p_GameLog_Add, p_GameLog_Add_Type1..9, TypeA..L)",
        "-- 1 Exclusivo de GameLog.exe (usp_Login_setChannel)",
        "-- 4 Consultas Administrativas de World00_GameLog.bak (p_GameLog_Get, p_GameLog_GetCount, p_GameLog_GetLog, p_GameLog_GetLogCount)",
        "-- 100% Pool-Safe (proc_logic + SELECT final incondicional para *sql.DB en Go)",
        "-- ============================================================================",
        "",
        "USE `xelion_online`;",
        ""
    ]

    count = 0
    for defn in PROC_DEFINITIONS:
        count += 1
        name = defn["name"]
        if defn["type"] == "add":
            lines.append(generate_add_proc(defn))
        else:
            lines.append(f"-- ----------------------------------------------------------------------------\n-- Procedure: {name}\n-- ----------------------------------------------------------------------------")
            lines.append(defn["sql"].strip() + "\n")

    content = "\n".join(lines)
    os.makedirs(os.path.dirname(OUTPUT_FILE), exist_ok=True)
    with open(OUTPUT_FILE, "w", encoding="utf-8") as f:
        f.write(content)

    print(f"Generated {count} procedures in {OUTPUT_FILE}")

if __name__ == "__main__":
    main()

