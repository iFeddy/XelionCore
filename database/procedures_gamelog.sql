-- ============================================================================
-- XELION ONLINE - STORED PROCEDURES OFICIALES DE WORLD00_GAMELOG (MYSQL 8.4)
-- Catálogo exhaustivo de 27 Stored Procedures oficiales
-- 22 Inserciones de Logs (p_GameLog_Add, p_GameLog_Add_Type1..9, TypeA..L)
-- 1 Exclusivo de GameLog.exe (usp_Login_setChannel)
-- 4 Consultas Administrativas de World00_GameLog.bak (p_GameLog_Get, p_GameLog_GetCount, p_GameLog_GetLog, p_GameLog_GetLogCount)
-- 100% Pool-Safe (proc_logic + SELECT final incondicional para *sql.DB en Go)
-- ============================================================================

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetCharNo, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type1
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type1`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type1`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type2
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type2`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type2`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type3
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type3`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type3`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT,
    IN in_nInt1 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nTargetID`, `nInt1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetCharNo, in_nTargetID, in_nInt1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type4
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type4`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type4`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nInt1`, `nInt2`, `nInt3`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nInt1, in_nInt2, in_nInt3);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type5
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type5`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type5`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type6
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type6`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type6`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nInt1 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nInt1`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nInt1, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type7
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type7`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type7`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, nTargetCharNo, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type8
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type8`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type8`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nInt1 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nInt1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nInt1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_Type9
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_Type9`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_Type9`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nItemKey`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nItemKey, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeA
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeA`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeA`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `nTargetCharNo`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`)
        VALUES (in_nType, in_nCharNo, in_nTargetCharNo, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeB
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeB`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeB`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetCharNo, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeC
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeC`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeC`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeD
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeD`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeD`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `nItemKey`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_nItemKey, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeE
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeE`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeE`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeF
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeF`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeF`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeG
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeG`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeG`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeH
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeH`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeH`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetCharNo);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeI
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeI`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeI`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetCharNo`, `nTargetID`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetCharNo, in_nTargetID);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeJ
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeJ`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeJ`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeK
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeK`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeK`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_sMap VARCHAR(16),
    IN in_nMapX INT,
    IN in_nMapY INT,
    IN in_nMapZ INT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `sMap`, `nMapX`, `nMapY`, `nMapZ`, `nInt1`, `nInt2`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_sMap, in_nMapX, in_nMapY, in_nMapZ, in_nInt1, in_nInt2, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Add_TypeL
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Add_TypeL`;
DELIMITER $$
CREATE PROCEDURE `p_GameLog_Add_TypeL`(
    OUT out_nLogNo BIGINT,
    IN in_nType INT,
    IN in_nCharNo INT,
    IN in_nTargetCharNo INT,
    IN in_nTargetID INT,
    IN in_nItemKey BIGINT,
    IN in_nInt1 INT,
    IN in_nInt2 INT,
    IN in_nInt3 INT,
    IN in_nBigint1 BIGINT
)
sp_body: BEGIN
    proc_logic: BEGIN
        SET out_nLogNo = 0;
        INSERT INTO `tGameLog` (`nType`, `nCharNo`, `nTargetCharNo`, `nTargetID`, `nItemKey`, `nInt1`, `nInt2`, `nInt3`, `nBigint1`)
        VALUES (in_nType, in_nCharNo, in_nTargetCharNo, in_nTargetID, in_nItemKey, in_nInt1, in_nInt2, in_nInt3, in_nBigint1);
        SET out_nLogNo = LAST_INSERT_ID();
        LEAVE proc_logic;
    END proc_logic;
    IF IFNULL(@_xc_nested_sp, 0) = 0 THEN
        SELECT out_nLogNo AS nLogNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_Get`;
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

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_GetCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_GetCount`;
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

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_GetLog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_GetLog`;
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

-- ----------------------------------------------------------------------------
-- Procedure: p_GameLog_GetLogCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GameLog_GetLogCount`;
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

-- ----------------------------------------------------------------------------
-- Procedure: usp_Login_setChannel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Login_setChannel`;
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
