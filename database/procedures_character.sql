-- ============================================================================
-- XELION ONLINE - STORED PROCEDURES OFICIALES DE WORLD00_CHARACTER (MYSQL 8.4)
-- Catálogo oficial de 311 Stored Procedures de Character.exe y WorldManager.exe
-- 100% de paridad con World00_Character.bak, PDBs oficiales y docs/ProcedimientosLegacy.md
-- ============================================================================

USE `xelion_online`;


-- ----------------------------------------------------------------------------
-- Procedure: p_Char_AddExp
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_AddExp`;
DELIMITER $$
CREATE PROCEDURE `p_Char_AddExp`(
    IN in_nCharNo INT,
    IN in_nExp BIGINT
)
sp_body: BEGIN

    UPDATE tCharacter SET nExp=nExp+in_nExp WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_ChangeClass
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_ChangeClass`;
DELIMITER $$
CREATE PROCEDURE `p_Char_ChangeClass`(
    IN in_nCharNo INT,
    IN in_nClass TINYINT
)
sp_body: BEGIN

    UPDATE tCharacterShape SET nClass=in_nClass WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_ChangeID
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_ChangeID`;
DELIMITER $$
CREATE PROCEDURE `p_Char_ChangeID`(
    IN in_nCharNo INT,
    IN in_sID VARCHAR(40),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacter SET sID=in_sID WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_CheckOptionsExistance
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_CheckOptionsExistance`;
DELIMITER $$
CREATE PROCEDURE `p_Char_CheckOptionsExistance`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT * FROM tCharacterOptions WHERE nCharNo = in_nCharNo ) THEN
        SET out_nRet=1;
    ELSE
        SET out_nRet=0;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Create
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Create`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Create`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_CreateSetDefaultData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_CreateSetDefaultData`;
DELIMITER $$
CREATE PROCEDURE `p_Char_CreateSetDefaultData`(
    IN in_nCharNo INT,
    IN in_sLoginZone VARCHAR(16),
    IN in_nLoginZoneX INT,
    IN in_nLoginZoneY INT,
    IN in_nHP SMALLINT,
    IN in_nSP SMALLINT,
    IN in_nHPS SMALLINT,
    IN in_nSPS SMALLINT,
    IN in_nMoney BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacter SET sLoginZone=in_sLoginZone,nLoginZoneX=in_nLoginZoneX,nLoginZoneY=in_nLoginZoneY,nHP=in_nHP,nSP=in_nSP,nHPS=in_nHPS,nSPS=in_nSPS,nMoney=in_nMoney WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_DelChargedBuff
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_DelChargedBuff`;
DELIMITER $$
CREATE PROCEDURE `p_Char_DelChargedBuff`(
    IN in_nCharNo INT,
    IN in_nKey INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tCharacterChargedBuffer WHERE nCharNo=in_nCharNo AND nKey=in_nKey;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_DelOptions
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_DelOptions`;
DELIMITER $$
CREATE PROCEDURE `p_Char_DelOptions`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tCharacterOptions WHERE nCharNo=in_nCharNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_DelState
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_DelState`;
DELIMITER $$
CREATE PROCEDURE `p_Char_DelState`(
    IN in_nCharNo INT,
    IN in_nStateNo SMALLINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tCharacterState WHERE nCharNo=in_nCharNo AND nStateNo=in_nStateNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_DelStateAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_DelStateAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_DelStateAll`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tCharacterState WHERE nCharNo=in_nCharNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Delete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Delete`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Delete`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_FameAdd
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_FameAdd`;
DELIMITER $$
CREATE PROCEDURE `p_Char_FameAdd`(
    IN in_nCharNo INT,
    IN in_nAddFame INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nCharFame INT;

    CALL p_Char_FameGet(in_nCharNo,v_nCharFame,out_nRet);
    IF out_nRet = 0 THEN
        LEAVE sp_body;
    END IF;
    SET v_nCharFame=v_nCharFame+in_nAddFame;
    CALL p_Char_FameSet(in_nCharNo,v_nCharFame,out_nRet);

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_FameGet
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_FameGet`;
DELIMITER $$
CREATE PROCEDURE `p_Char_FameGet`(
    IN in_nCharNo INT,
    OUT out_nCharFame INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT nFame INTO out_nCharFame FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 1 THEN
        SET out_nRet=1;
    END IF;
    LEAVE sp_body;

    SELECT out_nCharFame AS nCharFame, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_FameSet
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_FameSet`;
DELIMITER $$
CREATE PROCEDURE `p_Char_FameSet`(
    IN in_nCharNo INT,
    IN in_nSetFame INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tCharacter SET nFame=in_nSetFame WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 1 THEN
        SET out_nRet=1;
    END IF;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_FameSub
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_FameSub`;
DELIMITER $$
CREATE PROCEDURE `p_Char_FameSub`(
    IN in_nCharNo INT,
    IN in_nSubFame INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nCharFame INT;

    CALL p_Char_FameGet(in_nCharNo,v_nCharFame,out_nRet);
    IF out_nRet = 0 THEN
        LEAVE sp_body;
    END IF;
    SET v_nCharFame=v_nCharFame-in_nSubFame;
    IF v_nCharFame < 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    CALL p_Char_FameSet(in_nCharNo,v_nCharFame,out_nRet);

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Find
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Find`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Find`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetAllData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetAllData`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetAllData`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetCharID
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetCharID`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetCharID`(
    IN in_nCharNo INT,
    OUT out_sID VARCHAR(40),
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT sID INTO out_sID FROM tCharacter WHERE nCharNo=in_nCharNo;
    SET out_nRet=ROW_COUNT();

    SELECT out_sID AS sID, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetCharNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetCharNo`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetCharNo`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetChargedBufferAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetChargedBufferAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetChargedBufferAll`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nKey,nID,dUseDate,dEndDate FROM tCharacterChargedBuffer WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetChargedBufferCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetChargedBufferCount`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetChargedBufferCount`(
    IN in_nCharNo INT,
    OUT out_nCount INT
)
sp_body: BEGIN

    SELECT COUNT(*) INTO out_nCount FROM tCharacterChargedBuffer WHERE nCharNo=in_nCharNo;

    SELECT out_nCount AS nCount;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetIsDeleted
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetIsDeleted`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetIsDeleted`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted > 0 ) THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetKQMap
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetKQMap`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetKQMap`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nKQHandle,sKQMap,nKQX,nKQY,nKQD,dKQDate FROM tCharacter WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetListOfUserChar
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetListOfUserChar`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetListOfUserChar`(
    IN in_nUserNo INT
)
sp_body: BEGIN
    SELECT nCharNo FROM tCharacter WHERE nUserNo = in_nUserNo AND bDeleted = 0;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetListOfUserCharAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetListOfUserCharAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetListOfUserCharAll`(
    IN in_nUserNo INT
)
sp_body: BEGIN

    SELECT nCharNo FROM tCharacter WHERE nUserNo=in_nUserNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetLoginData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetLoginData`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetLoginData`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetNumOfUserChar
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetNumOfUserChar`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetNumOfUserChar`(
    IN in_nUserNo INT,
    OUT out_nNumOfUserChar INT
)
sp_body: BEGIN

    SET out_nNumOfUserChar=0;
    SELECT COUNT(nCharNo) INTO out_nNumOfUserChar FROM tCharacter WHERE nUserNo=in_nUserNo AND bDeleted=0;

    SELECT out_nNumOfUserChar AS nNumOfUserChar;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetNumOfUserCharAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetNumOfUserCharAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetNumOfUserCharAll`(
    IN in_nUserNo INT,
    OUT out_nNumOfUserChar INT
)
sp_body: BEGIN

    SET out_nNumOfUserChar=0;
    SELECT COUNT(nCharNo) INTO out_nNumOfUserChar FROM tCharacter WHERE nUserNo=in_nUserNo;

    SELECT out_nNumOfUserChar AS nNumOfUserChar;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptGame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptGame`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptGame`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(64)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sGame,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sGame IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptKeyMapping
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptKeyMapping`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptKeyMapping`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(308)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sKeyMapping,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sKeyMapping IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptShortCutData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptShortCutData`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptShortCutData`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(1024)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sShortCutData,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sShortCutData IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptShortCutSize
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptShortCutSize`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptShortCutSize`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(24)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sShortCutSize,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sShortCutSize IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptSound
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptSound`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptSound`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(1)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sSound,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sSound IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptVideo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptVideo`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptVideo`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(60)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sVideo,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sVideo IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetOptWindowsPos
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetOptWindowsPos`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetOptWindowsPos`(
    IN in_nCharNo INT,
    OUT out_nRet TINYINT,
    OUT out_sData VARBINARY(392)
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT sWindowsPos,1 INTO out_sData,out_nRet FROM tCharacterOptions WHERE nCharNo=in_nCharNo AND sWindowsPos IS NOT NULL;

    SELECT out_nRet AS nRet, out_sData AS sData;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetShape
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetShape`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetShape`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nRace,nClass,nGender,nHairType,nHairColor,nFaceShape FROM tCharacterShape WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetState
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetState`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetState`(
    IN in_nCharNo INT,
    IN in_nStateNo SMALLINT,
    OUT out_nStateLevel SMALLINT,
    OUT out_nStateTime INT,
    OUT out_nStateCount SMALLINT,
    OUT out_nStatePowerDemage TINYINT,
    OUT out_nStatePowerSP TINYINT,
    OUT out_nStatePowerKeepTime TINYINT,
    OUT out_nStatePowerCoolTime TINYINT,
    OUT out_nCasterLevel SMALLINT,
    OUT out_nCasterAttackPower INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT nStateLevel,nStateTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nCasterLevel,out_nCasterAttackPower INTO out_nStateLevel,out_nStateTime,out_nStateCount,out_nStatePowerDemage,out_nStatePowerSP,out_nStatePowerKeepTime,out_nStatePowerCoolTime,out_nCasterLevel,out_nCasterAttackPower FROM tCharacterState WHERE nCharNo=in_nCharNo AND nStateNo=in_nStateNo;
    SET out_nRet=ROW_COUNT();

    SELECT out_nStateLevel AS nStateLevel, out_nStateTime AS nStateTime, out_nStateCount AS nStateCount, out_nStatePowerDemage AS nStatePowerDemage, out_nStatePowerSP AS nStatePowerSP, out_nStatePowerKeepTime AS nStatePowerKeepTime, out_nStatePowerCoolTime AS nStatePowerCoolTime, out_nCasterLevel AS nCasterLevel, out_nCasterAttackPower AS nCasterAttackPower, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetStateAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetStateAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetStateAll`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nStateNo,nStateLevel,nStateTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nCasterLevel,nCasterAttackPower FROM tCharacterState WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetStateAllWithSetItemAndRestTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetStateAllWithSetItemAndRestTime`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetStateAllWithSetItemAndRestTime`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nStateNo,nStateLevel,nStateTime,nRestTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nSetItemIndexA,nSetItemValueA,nSetItemIndexB,nSetItemValueB,nSetItemIndexC,nSetItemValueC,nSetItemIndexD,nSetItemValueD,nCasterLevel,nCasterAttackPower FROM tCharacterState WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetTitleAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetTitleAll`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetTitleAll`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nType,nElementNo,nElementValue FROM tCharacterTitle WHERE nCharNo=in_nCharNo ORDER BY nNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetTitleCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetTitleCount`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetTitleCount`(
    IN in_nCharNo INT,
    OUT out_nCount INT
)
sp_body: BEGIN

    SET out_nCount=0;
    SELECT COUNT(*) INTO out_nCount FROM tCharacterTitle WHERE nCharNo=in_nCharNo;

    SELECT out_nCount AS nCount;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_GetUserNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_GetUserNo`;
DELIMITER $$
CREATE PROCEDURE `p_Char_GetUserNo`(
    IN in_nCharNo INT,
    OUT out_nUserNo INT
)
sp_body: BEGIN

    SET out_nUserNo=0;
    SELECT nUserNo INTO out_nUserNo FROM tCharacter WHERE nCharNo=in_nCharNo;

    SELECT out_nUserNo AS nUserNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Login
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Login`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Login`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Logout
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Logout`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Logout`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_MoneyAdd
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_MoneyAdd`;
DELIMITER $$
CREATE PROCEDURE `p_Char_MoneyAdd`(
    IN in_nCharNo INT,
    IN in_nAddMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nCharMoney BIGINT;

    CALL p_Char_MoneyGet(in_nCharNo,v_nCharMoney,out_nRet);
    IF out_nRet = 0 THEN
        LEAVE sp_body;
    END IF;
    SET v_nCharMoney=v_nCharMoney+in_nAddMoney;
    CALL p_Char_MoneySet(in_nCharNo,v_nCharMoney,out_nRet);

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_MoneyGet
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_MoneyGet`;
DELIMITER $$
CREATE PROCEDURE `p_Char_MoneyGet`(
    IN in_nCharNo INT,
    OUT out_nCharMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT nMoney INTO out_nCharMoney FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 1 THEN
        SET out_nRet=1;
    END IF;
    LEAVE sp_body;

    SELECT out_nCharMoney AS nCharMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_MoneyMove
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_MoneyMove`;
DELIMITER $$
CREATE PROCEDURE `p_Char_MoneyMove`(
    IN in_nCharNoFrom INT,
    IN in_nCharNoTo INT,
    IN in_nMoveMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_money BIGINT;

    IF in_nMoveMoney < 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    IF in_nMoveMoney = 0 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNoFrom AND bDeleted = 0 ) THEN
        SET out_nRet=0;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNoTo AND bDeleted = 0 ) THEN
        SET out_nRet=0;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    SELECT nMoney INTO v_money FROM tCharacter WHERE nCharNo=in_nCharNoFrom AND bDeleted=0;
    IF v_money >= in_nMoveMoney THEN
        UPDATE tCharacter SET nMoney=v_money-in_nMoveMoney WHERE nCharNo=in_nCharNoFrom AND bDeleted=0;
        IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMoney=nMoney+in_nMoveMoney WHERE nCharNo=in_nCharNoTo AND bDeleted=0;
        IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SET out_nRet=1;
        COMMIT;
    ELSE
        SET out_nRet=0;
        ROLLBACK;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_MoneySub
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_MoneySub`;
DELIMITER $$
CREATE PROCEDURE `p_Char_MoneySub`(
    IN in_nCharNo INT,
    IN in_nSubMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nCharMoney BIGINT;

    CALL p_Char_MoneyGet(in_nCharNo,v_nCharMoney,out_nRet);
    IF out_nRet = 0 THEN
        LEAVE sp_body;
    END IF;
    SET v_nCharMoney=v_nCharMoney-in_nSubMoney;
    IF v_nCharMoney < 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneySet(in_nCharNo,v_nCharMoney,out_nRet);

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Rename
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Rename`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Rename`(
    IN in_nCharNo INT,
    IN in_sID VARCHAR(40),
    OUT out_nRet INT
)
sp_body: BEGIN

    IF NOT EXISTS ( SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = in_nCharNo ) THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    UPDATE tCharacter SET sID=in_sID WHERE nCharNo=in_nCharNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SaveLevelExpFrame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SaveLevelExpFrame`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SaveLevelExpFrame`(
    IN in_nCharNo INT,
    IN in_nLevel SMALLINT,
    IN in_nExp BIGINT,
    IN in_nFame INT
)
sp_body: BEGIN

    UPDATE tCharacter SET nLevel=in_nLevel,nExp=in_nExp,nFame=in_nFame WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SaveLocation
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SaveLocation`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SaveLocation`(
    IN in_nCharNo INT,
    IN in_sLoginZone VARCHAR(16),
    IN in_nLoginZoneX INT,
    IN in_nLoginZoneY INT,
    IN in_nKQHandle INT,
    IN in_sKQMap VARCHAR(16),
    IN in_nKQX INT,
    IN in_nKQY INT
)
sp_body: BEGIN

    UPDATE tCharacter SET sLoginZone=in_sLoginZone,nLoginZoneX=in_nLoginZoneX,nLoginZoneY=in_nLoginZoneY,nKQHandle=in_nKQHandle,sKQMap=in_sKQMap,nKQX=in_nKQX,nKQY=in_nKQY,dKQDate=NOW() WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SavePKCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SavePKCount`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SavePKCount`(
    IN in_nCharNo INT,
    IN in_nPKCount INT
)
sp_body: BEGIN

    UPDATE tCharacter SET nPKCount=in_nPKCount WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SaveStat
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SaveStat`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SaveStat`(
    IN in_nCharNo INT,
    IN in_nPS INT,
    IN in_nGS INT,
    IN in_nHPS SMALLINT,
    IN in_nSPS SMALLINT,
    IN in_nHP SMALLINT,
    IN in_nSP SMALLINT
)
sp_body: BEGIN

    UPDATE tCharacter SET nPS=in_nPS,nGS=in_nGS,nHPS=in_nHPS,nSPS=in_nSPS,nHP=in_nHP,nSP=in_nSP WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetChargedBuffer
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetChargedBuffer`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetChargedBuffer`(
    IN in_nCharNo INT,
    IN in_nKey INT,
    IN in_nID SMALLINT,
    IN in_dUseDate DATETIME,
    IN in_dEndDate DATETIME,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacterChargedBuffer WHERE nCharNo = in_nCharNo AND nKey = in_nKey ) THEN
        UPDATE tCharacterChargedBuffer SET nID=in_nID,dUseDate=in_dUseDate,dEndDate=in_dEndDate WHERE nCharNo=in_nCharNo AND nKey=in_nKey;
    ELSE
        INSERT tCharacterChargedBuffer(nCharNo,nKey,nID,dUseDate,dEndDate) VALUES (in_nCharNo,in_nKey,in_nID,in_dUseDate,in_dEndDate);
    END IF;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetLevelExp
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetLevelExp`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetLevelExp`(
    IN in_nCharNo INT,
    IN in_nLevel SMALLINT,
    IN in_nExp BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacter SET nLevel=in_nLevel,nExp=in_nExp WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptGame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptGame`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptGame`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(64),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sGame=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sGame) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptKeyMapping
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptKeyMapping`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptKeyMapping`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(308),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sKeyMapping=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sKeyMapping) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptShortCutData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptShortCutData`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptShortCutData`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(1024),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sShortCutData=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sShortCutData) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptShortCutSize
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptShortCutSize`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptShortCutSize`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(24),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sShortCutSize=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sShortCutSize) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptSound
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptSound`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptSound`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(1),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sSound=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sSound) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptVideo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptVideo`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptVideo`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(60),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sVideo=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sVideo) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetOptWindowsPos
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetOptWindowsPos`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetOptWindowsPos`(
    IN in_nCharNo INT,
    IN in_sData VARBINARY(392),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tCharacterOptions SET sWindowsPos=in_sData WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tCharacterOptions(nCharNo,sWindowsPos) VALUES (in_nCharNo,in_sData);
        SET out_nRet=ROW_COUNT();
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetShape
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetShape`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetShape`(
    IN in_nCharNo INT,
    IN in_nRace TINYINT,
    IN in_nClass TINYINT,
    IN in_nGender TINYINT,
    IN in_nHairType TINYINT,
    IN in_nHairColor TINYINT,
    IN in_nFaceShape TINYINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    UPDATE tCharacterShape SET nRace=in_nRace,nClass=in_nClass,nGender=in_nGender,nHairType=in_nHairType,nHairColor=in_nHairColor,nFaceShape=in_nFaceShape WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetState
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetState`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetState`(
    IN in_nCharNo INT,
    IN in_nStateNo SMALLINT,
    IN in_nStateLevel SMALLINT,
    IN in_nStateTime INT,
    IN in_nStateCount SMALLINT,
    IN in_nStatePowerDemage TINYINT,
    IN in_nStatePowerSP TINYINT,
    IN in_nStatePowerKeepTime TINYINT,
    IN in_nStatePowerCoolTime TINYINT,
    IN in_nCasterLevel SMALLINT,
    IN in_nCasterAttackPower INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacterState WHERE nCharNo = in_nCharNo AND nStateNo = in_nStateNo ) THEN
        UPDATE tCharacterState SET nStateLevel=in_nStateLevel,nStateTime=in_nStateTime,nStateCount=in_nStateCount,nStatePowerDemage=in_nStatePowerDemage,nStatePowerSP=in_nStatePowerSP,nStatePowerKeepTime=in_nStatePowerKeepTime,nStatePowerCoolTime=in_nStatePowerCoolTime,nCasterLevel=in_nCasterLevel,nCasterAttackPower=in_nCasterAttackPower WHERE nCharNo=in_nCharNo AND nStateNo=in_nStateNo;
        SET out_nRet=0;
    ELSE
        INSERT tCharacterState(nCharNo,nStateNo,nStateLevel,nStateTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nCasterLevel,nCasterAttackPower) VALUES (in_nCharNo,in_nStateNo,in_nStateLevel,in_nStateTime,in_nStateCount,in_nStatePowerDemage,in_nStatePowerSP,in_nStatePowerKeepTime,in_nStatePowerCoolTime,in_nCasterLevel,in_nCasterAttackPower);
        SET out_nRet=0;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetStateWithSetItem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetStateWithSetItem`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetStateWithSetItem`(
    IN in_nCharNo INT,
    IN in_nStateNo SMALLINT,
    IN in_nStateLevel SMALLINT,
    IN in_nStateTime INT,
    IN in_nStateCount SMALLINT,
    IN in_nStatePowerDemage TINYINT,
    IN in_nStatePowerSP TINYINT,
    IN in_nStatePowerKeepTime TINYINT,
    IN in_nStatePowerCoolTime TINYINT,
    IN in_nSetItemIndexA SMALLINT,
    IN in_nSetItemValueA SMALLINT,
    IN in_nSetItemIndexB SMALLINT,
    IN in_nSetItemValueB SMALLINT,
    IN in_nSetItemIndexC SMALLINT,
    IN in_nSetItemValueC SMALLINT,
    IN in_nSetItemIndexd SMALLINT,
    IN in_nSetItemValueD SMALLINT,
    IN in_nCasterLevel SMALLINT,
    IN in_nCasterAttackPower INT,
    IN in_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacterState WHERE nCharNo = in_nCharNo AND nStateNo = in_nStateNo ) THEN
        UPDATE tCharacterState SET nStateLevel=in_nStateLevel,nStateTime=in_nStateTime,nStateCount=in_nStateCount,nStatePowerDemage=in_nStatePowerDemage,nStatePowerSP=in_nStatePowerSP,nStatePowerKeepTime=in_nStatePowerKeepTime,nStatePowerCoolTime=in_nStatePowerCoolTime,nSetItemIndexA=in_nSetItemIndexA,nSetItemValueA=in_nSetItemValueA,nSetItemIndexB=in_nSetItemIndexB,nSetItemValueB=in_nSetItemValueB,nSetItemIndexC=in_nSetItemIndexC,nSetItemValueC=in_nSetItemValueC,nSetItemIndexD=in_nSetItemIndexd,nSetItemValueD=in_nSetItemValueD,nCasterLevel=in_nCasterLevel,nCasterAttackPower=in_nCasterAttackPower WHERE nCharNo=in_nCharNo AND nStateNo=in_nStateNo;
        SET in_nRet=0;
    ELSE
        INSERT tCharacterState(nCharNo,nStateNo,nStateLevel,nStateTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nSetItemIndexA,nSetItemValueA,nSetItemIndexB,nSetItemValueB,nSetItemIndexC,nSetItemValueC,nSetItemIndexD,nSetItemValueD,nCasterLevel,nCasterAttackPower) VALUES (in_nCharNo,in_nStateNo,in_nStateLevel,in_nStateTime,in_nStateCount,in_nStatePowerDemage,in_nStatePowerSP,in_nStatePowerKeepTime,in_nStatePowerCoolTime,in_nSetItemIndexA,in_nSetItemValueA,in_nSetItemIndexB,in_nSetItemValueB,in_nSetItemIndexC,in_nSetItemValueC,in_nSetItemIndexd,in_nSetItemValueD,in_nCasterLevel,in_nCasterAttackPower);
        SET in_nRet=0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetStateWithSetItemAndRestTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetStateWithSetItemAndRestTime`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetStateWithSetItemAndRestTime`(
    IN in_nCharNo INT,
    IN in_nStateNo SMALLINT,
    IN in_nStateLevel SMALLINT,
    IN in_nStateTime INT,
    IN in_nRestTime INT,
    IN in_nStateCount SMALLINT,
    IN in_nStatePowerDemage TINYINT,
    IN in_nStatePowerSP TINYINT,
    IN in_nStatePowerKeepTime TINYINT,
    IN in_nStatePowerCoolTime TINYINT,
    IN in_nSetItemIndexA SMALLINT,
    IN in_nSetItemValueA SMALLINT,
    IN in_nSetItemIndexB SMALLINT,
    IN in_nSetItemValueB SMALLINT,
    IN in_nSetItemIndexC SMALLINT,
    IN in_nSetItemValueC SMALLINT,
    IN in_nSetItemIndexd SMALLINT,
    IN in_nSetItemValueD SMALLINT,
    IN in_nCasterLevel SMALLINT,
    IN in_nCasterAttackPower INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacterState WHERE nCharNo = in_nCharNo AND nStateNo = in_nStateNo ) THEN
        UPDATE tCharacterState SET nStateLevel=in_nStateLevel,nStateTime=in_nStateTime,nRestTime=in_nRestTime,nStateCount=in_nStateCount,nStatePowerDemage=in_nStatePowerDemage,nStatePowerSP=in_nStatePowerSP,nStatePowerKeepTime=in_nStatePowerKeepTime,nStatePowerCoolTime=in_nStatePowerCoolTime,nSetItemIndexA=in_nSetItemIndexA,nSetItemValueA=in_nSetItemValueA,nSetItemIndexB=in_nSetItemIndexB,nSetItemValueB=in_nSetItemValueB,nSetItemIndexC=in_nSetItemIndexC,nSetItemValueC=in_nSetItemValueC,nSetItemIndexD=in_nSetItemIndexd,nSetItemValueD=in_nSetItemValueD,nCasterLevel=in_nCasterLevel,nCasterAttackPower=in_nCasterAttackPower WHERE nCharNo=in_nCharNo AND nStateNo=in_nStateNo;
        SET out_nRet=0;
    ELSE
        INSERT tCharacterState(nCharNo,nStateNo,nStateLevel,nStateTime,nRestTime,nStateCount,nStatePowerDemage,nStatePowerSP,nStatePowerKeepTime,nStatePowerCoolTime,nSetItemIndexA,nSetItemValueA,nSetItemIndexB,nSetItemValueB,nSetItemIndexC,nSetItemValueC,nSetItemIndexD,nSetItemValueD,nCasterLevel,nCasterAttackPower) VALUES (in_nCharNo,in_nStateNo,in_nStateLevel,in_nStateTime,in_nRestTime,in_nStateCount,in_nStatePowerDemage,in_nStatePowerSP,in_nStatePowerKeepTime,in_nStatePowerCoolTime,in_nSetItemIndexA,in_nSetItemValueA,in_nSetItemIndexB,in_nSetItemValueB,in_nSetItemIndexC,in_nSetItemValueC,in_nSetItemIndexd,in_nSetItemValueD,in_nCasterLevel,in_nCasterAttackPower);
        SET out_nRet=0;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SetTitle
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SetTitle`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SetTitle`(
    IN in_nCharNo INT,
    IN in_nType TINYINT,
    IN in_nElementNo TINYINT,
    IN in_nElementValue BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterTitle WHERE nCharNo = in_nCharNo AND nType = in_nType AND nElementNo = in_nElementNo ) THEN
        UPDATE tCharacterTitle SET nElementValue=in_nElementValue WHERE nCharNo=in_nCharNo AND nType=in_nType AND nElementNo=in_nElementNo;
        IF ROW_COUNT ( ) = 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT tCharacterTitle(nCharNo,nType,nElementNo,nElementValue) VALUES (in_nCharNo,in_nType,in_nElementNo,in_nElementValue);
        IF 0 = 0 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_StatPoint_Dec
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_StatPoint_Dec`;
DELIMITER $$
CREATE PROCEDURE `p_Char_StatPoint_Dec`(
    IN in_nCharNo INT,
    IN in_nStatPointType TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nStatPoint TINYINT;

    IF in_nStatPointType = 0 THEN
        SELECT nStrength INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nStrength=nStrength-1 WHERE nCharNo=in_nCharNo AND nStrength>0;
    END IF;
    IF in_nStatPointType = 1 THEN
        SELECT nConstitute INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nConstitute=nConstitute-1 WHERE nCharNo=in_nCharNo AND nConstitute>0;
    END IF;
    IF in_nStatPointType = 2 THEN
        SELECT nDexterity INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nDexterity=nDexterity-1 WHERE nCharNo=in_nCharNo AND nDexterity>0;
    END IF;
    IF in_nStatPointType = 3 THEN
        SELECT nIntelligence INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nIntelligence=nIntelligence-1 WHERE nCharNo=in_nCharNo AND nIntelligence>0;
    END IF;
    IF in_nStatPointType = 4 THEN
        SELECT nMentalPower INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMentalPower=nMentalPower-1 WHERE nCharNo=in_nCharNo AND nMentalPower>0;
    END IF;
    IF 0 <> 0 OR in_nStatPointType > 4 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_StatPoint_Inc
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_StatPoint_Inc`;
DELIMITER $$
CREATE PROCEDURE `p_Char_StatPoint_Inc`(
    IN in_nCharNo INT,
    IN in_nStatPointType TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF in_nStatPointType = 0 THEN
        UPDATE tCharacter SET nStrength=nStrength+1 WHERE nCharNo=in_nCharNo;
    END IF;
    IF in_nStatPointType = 1 THEN
        UPDATE tCharacter SET nConstitute=nConstitute+1 WHERE nCharNo=in_nCharNo;
    END IF;
    IF in_nStatPointType = 2 THEN
        UPDATE tCharacter SET nDexterity=nDexterity+1 WHERE nCharNo=in_nCharNo;
    END IF;
    IF in_nStatPointType = 3 THEN
        UPDATE tCharacter SET nIntelligence=nIntelligence+1 WHERE nCharNo=in_nCharNo;
    END IF;
    IF in_nStatPointType = 4 THEN
        UPDATE tCharacter SET nMentalPower=nMentalPower+1 WHERE nCharNo=in_nCharNo;
    END IF;
    IF 0 <> 0 OR in_nStatPointType > 4 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_SubExp
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_SubExp`;
DELIMITER $$
CREATE PROCEDURE `p_Char_SubExp`(
    IN in_nCharNo INT,
    IN in_nExp BIGINT
)
sp_body: BEGIN

    UPDATE tCharacter SET nExp=nExp-in_nExp WHERE nCharNo=in_nCharNo;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_WeddingStatus_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_WeddingStatus_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Char_WeddingStatus_Get`(
    IN in_CharNo INT
)
sp_body: BEGIN

    IF NOT EXISTS ( SELECT * FROM tMarriage WHERE nCharNo = in_CharNo ) THEN
        INSERT tMarriage(nCharNo,nPartnerCharNo,dWhen,dSubWhen) VALUES (in_CharNo,0,NOW(),NOW());
    END IF;
    SELECT nPartnerCharNo,nStatus,nSubStatus,dWhen,dSubWhen,nMarriageCount,nDivorceCancel FROM tMarriage WHERE nCharNo=in_CharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_WeddingStatus_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_WeddingStatus_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Char_WeddingStatus_Set`(
    IN in_CharNo INT,
    IN in_PartnerCharNo INT,
    IN in_Status SMALLINT,
    IN in_SubStatus SMALLINT,
    IN in_When DATETIME,
    IN in_SubWhen DATETIME,
    IN in_MarriageCount TINYINT,
    IN in_DivorceCancel TINYINT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharNo ) THEN
        UPDATE tMarriage SET nPartnerCharNo=in_PartnerCharNo,nStatus=in_Status,nSubStatus=in_SubStatus,dWhen=in_When,dSubWhen=in_SubWhen,nMarriageCount=in_MarriageCount,nDivorceCancel=in_DivorceCancel WHERE nCharNo=in_CharNo;
    ELSE
        INSERT tMarriage(nCharNo,nPartnerCharNo,nStatus,nSubStatus,dWhen,dSubWhen) VALUES (in_CharNo,in_PartnerCharNo,in_Status,in_SubStatus,in_When,in_SubWhen);
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Cancel_Propose
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Cancel_Propose`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Cancel_Propose`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nStatus=0,nSubStatus=0 WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nStatus=0,nSubStatus=0 WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Divorce_Cancel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Divorce_Cancel`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Divorce_Cancel`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Divorce_Do
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Divorce_Do`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Divorce_Do`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nStatus=3,nSubStatus=1,dSubWhen=DATE_ADD(NOW(),INTERVAL 7 DAY) WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nStatus=3,nSubStatus=1,dSubWhen=DATE_ADD(NOW(),INTERVAL 7 DAY) WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Divorce_Do_by_Item
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Divorce_Do_by_Item`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Divorce_Do_by_Item`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nStatus=3,nSubStatus=1,dSubWhen=DATE_ADD(NOW(),INTERVAL 5 SECOND) WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nStatus=3,nSubStatus=2,dSubWhen=DATE_ADD(NOW(),INTERVAL 5 SECOND) WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Divorce_by_Agree
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Divorce_by_Agree`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Divorce_by_Agree`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nSubStatus=1,dSubWhen=DATE_ADD(NOW(),INTERVAL 10 MINUTE) WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nSubStatus=2+4,dSubWhen=DATE_ADD(NOW(),INTERVAL 10 MINUTE) WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Divorce_by_Force
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Divorce_by_Force`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Divorce_by_Force`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nSubStatus=3,dSubWhen=DATE_ADD(NOW(),INTERVAL 3 DAY) WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nSubStatus=4,dSubWhen=DATE_ADD(NOW(),INTERVAL 3 DAY) WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Do
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Do`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Do`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo AND nPartnerCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nStatus=2,nSubStatus=0,dWhen=NOW(),dSubWhen=DATE_ADD(NOW(),INTERVAL 7 DAY),nMarriageCount=nMarriageCount+1,nDivorceCancel=0 WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo AND nPartnerCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nStatus=2,nSubStatus=0,dWhen=NOW(),dSubWhen=DATE_ADD(NOW(),INTERVAL 7 DAY),nMarriageCount=nMarriageCount+1,nDivorceCancel=0 WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        ROLLBACK;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Escape_Divorce_State
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Escape_Divorce_State`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Escape_Divorce_State`(
    IN in_CharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharNo ) THEN
        SET out_nRet=0;
        UPDATE tMarriage SET nStatus=0,nSubStatus=0,nDivorceCancel=0 WHERE nCharNo=in_CharNo;
    ELSE
        SET out_nRet=1;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Get_Partner_Info
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Get_Partner_Info`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Get_Partner_Info`(
    IN in_CharNo INT
)
sp_body: BEGIN

    SELECT nLevel,dLastLoginDate FROM tCharacter WHERE nCharNo=in_CharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wedding_Propose
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wedding_Propose`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wedding_Propose`(
    IN in_CharANo INT,
    IN in_CharBNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharANo ) THEN
        UPDATE tMarriage SET nPartnerCharNo=in_CharBNo,nStatus=1,nSubStatus=0,dWhen=NOW(),dSubWhen=DATE_ADD(NOW(),INTERVAL 3 DAY) WHERE nCharNo=in_CharANo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        INSERT tMarriage(nCharNo,nPartnerCharNo,nStatus,nSubStatus,dWhen,dSubWhen) VALUES (in_CharANo,in_CharBNo,1,0,NOW(),DATE_ADD(NOW(),INTERVAL 3 DAY));
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_CharBNo ) THEN
        UPDATE tMarriage SET nPartnerCharNo=in_CharANo,nStatus=1,nSubStatus=0,dWhen=NOW(),dSubWhen=DATE_ADD(NOW(),INTERVAL 3 DAY) WHERE nCharNo=in_CharBNo;
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    ELSE
        INSERT tMarriage(nCharNo,nPartnerCharNo,nStatus,nSubStatus,dWhen,dSubWhen) VALUES (in_CharBNo,in_CharANo,1,0,NOW(),DATE_ADD(NOW(),INTERVAL 3 DAY));
        IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Char_Wipe
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Char_Wipe`;
DELIMITER $$
CREATE PROCEDURE `p_Char_Wipe`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tFriend WHERE nCharNo=in_nCharNo;
    DELETE FROM tFriend WHERE nFriendCharNo=in_nCharNo;
    DELETE FROM tCharacter WHERE nCharNo=in_nCharNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_Del
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_Del`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_Del`(
    IN in_nCharNo INT,
    IN in_nFriendCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tFriend WHERE nCharNo=in_nCharNo AND nFriendCharNo=in_nFriendCharNo;
    SET out_nRet=0;
    DELETE FROM tFriend WHERE nCharNo=in_nFriendCharNo AND nFriendCharNo=in_nCharNo;
    SET out_nRet=out_nRet+0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_Del_All
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_Del_All`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_Del_All`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tFriend WHERE nCharNo=in_nCharNo OR nFriendCharNo=in_nCharNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_Get`(
    IN in_nCharNo INT,
    IN in_nFriendNo INT
)
sp_body: BEGIN

    SELECT F.nFriendCharNo,F.nFlags,C.sID,S.nClass,C.nLevel,C.sLoginZone,C.dLastLoginDate FROM tFriend AS F INNER JOIN tCharacterShape AS S ON S.nCharNo=in_nFriendNo INNER JOIN tCharacter AS C ON C.nCharNo=in_nFriendNo AND C.bDeleted=0 WHERE F.nCharNo=in_nCharNo AND F.nFriendCharNo=in_nFriendNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_GetAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_GetAll`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_GetAll`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT F.nFriendCharNo,F.nFlags,C.sID FROM tFriend AS F INNER JOIN tCharacter AS C ON C.nCharNo=F.nFriendCharNo AND C.bDeleted=0 WHERE F.nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_GetNumOfFriend
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_GetNumOfFriend`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_GetNumOfFriend`(
    IN in_nCharNo INT,
    OUT out_nNumOfFriend INT
)
sp_body: BEGIN

    SET out_nNumOfFriend=0;
    SELECT COUNT(*) INTO out_nNumOfFriend FROM tFriend AS F INNER JOIN tCharacter AS C ON C.nCharNo=F.nFriendCharNo AND C.bDeleted=0 WHERE F.nCharNo=in_nCharNo;

    SELECT out_nNumOfFriend AS nNumOfFriend;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Friend_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Friend_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Friend_Set`(
    IN in_nMAXFRIEND SMALLINT,
    IN in_nCharNo INT,
    IN in_nFriendCharNo INT,
    IN in_nFlags TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nFriendCount INT;

    IF in_nCharNo = in_nFriendCharNo THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Friend_GetNumOfFriend(in_nCharNo,v_nFriendCount);
    IF v_nFriendCount >= in_nMAXFRIEND THEN
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    CALL p_Friend_GetNumOfFriend(in_nFriendCharNo,v_nFriendCount);
    IF v_nFriendCount >= in_nMAXFRIEND THEN
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    INSERT tFriend(nCharNo,nFriendCharNo,nFlags) VALUES (in_nCharNo,in_nFriendCharNo,in_nFlags);
    SET out_nRet=0;
    IF ( out_nRet <> 0 ) THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    INSERT tFriend(nCharNo,nFriendCharNo,nFlags) VALUES (in_nFriendCharNo,in_nCharNo,in_nFlags);
    SET out_nRet=0;
    IF ( out_nRet <> 0 ) THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GuildTournament_GetLastMatch
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GuildTournament_GetLastMatch`;
DELIMITER $$
CREATE PROCEDURE `p_GuildTournament_GetLastMatch`()
sp_body: BEGIN

    SELECT nGTNo,nType,nData FROM tGuildTournament WHERE (nGTNo IN(SELECT MAX(nGTNo) FROM tGuildTournament)) ORDER BY nType;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_GuildTournament_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_GuildTournament_Set`;
DELIMITER $$
CREATE PROCEDURE `p_GuildTournament_Set`(
    IN in_nGTNo INT,
    IN in_nType SMALLINT,
    IN in_nData INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tGuildTournament SET nData=in_nData WHERE nGTNo=in_nGTNo AND nType=in_nType;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tGuildTournament(nGTNo,nType,nData) VALUES (in_nGTNo,in_nType,in_nData);
    END IF;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Delete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Delete`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Delete`(
    IN in_nNo INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    DELETE FROM tGuild WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Get`(
    IN in_nNo INT
)
sp_body: BEGIN

    SELECT sName,sPassword,nMaxMembers,nMoney,nType,nGrade,nFame,nStoneLevel,nExp,dCreateDate,dWarRequestDate,dWarStartDate,dWarEndDate,nWarEnemy,nWarWinCount,nWarLoseCount,nWarDrawCount,nDismissStatus,dDismissDate,dNotifyDate,sNotifyCharID,sNotify,sIntro FROM tGuild WHERE in_nNo=nNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_GetAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_GetAll`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_GetAll`()
sp_body: BEGIN

    SELECT nNo,sName,sPassword,nMaxMembers,nMoney,nType,nGrade,nFame,nStoneLevel,nExp,dCreateDate,dWarRequestDate,dWarStartDate,dWarEndDate,nWarEnemy,nWarWinCount,nWarLoseCount,nWarDrawCount,nDismissStatus,dDismissDate,dNotifyDate,sNotifyCharID,sNotify,sIntro FROM tGuild;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_GetMyNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_GetMyNo`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_GetMyNo`(
    IN in_nCharNo INT,
    OUT out_nNo INT
)
sp_body: BEGIN

    SET out_nNo=0;
    SELECT nNo INTO out_nNo FROM tGuildMember WHERE nCharNo=in_nCharNo;

    SELECT out_nNo AS nNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Make2
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Make2`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Make2`(
    IN in_sName VARCHAR(18),
    IN in_sPassword VARCHAR(9),
    IN in_nMaxMembers SMALLINT,
    IN in_nType TINYINT,
    IN in_nMakerCharNo INT,
    IN in_nMakerPay BIGINT,
    IN in_nMakerNeedItemID INT,
    OUT out_nNo INT,
    OUT out_nRemainCharMoney BIGINT
)
sp_body: BEGIN

    DECLARE v_nRet TINYINT;

    SET out_nNo=0;
    IF EXISTS ( SELECT sName FROM tGuild WHERE sName = in_sName ) THEN
        SET out_nNo=-1;
        LEAVE sp_body;
    END IF;
    IF in_nMakerNeedItemID <> 65535 THEN
        IF NOT EXISTS ( SELECT nItemID FROM tItem WHERE nOwner = in_nMakerCharNo AND nItemID = in_nMakerNeedItemID ) THEN
        SET out_nNo=-2;
        LEAVE sp_body;
    END IF;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tGuildMember WHERE nCharNo = in_nMakerCharNo ) THEN
        SET out_nNo=-3;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneySub(in_nMakerCharNo,in_nMakerPay,v_nRet);
    IF v_nRet = 0 THEN
        ROLLBACK;
        SET out_nNo=-4;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneyGet(in_nMakerCharNo,out_nRemainCharMoney,v_nRet);
    IF v_nRet = 0 THEN
        ROLLBACK;
        SET out_nNo=-4;
        LEAVE sp_body;
    END IF;
    INSERT tGuild(sName,sPassword,nMaxMembers,nType) VALUES (in_sName,in_sPassword,in_nMaxMembers,in_nType);
    SET out_nNo=LAST_INSERT_ID();
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nNo=-5;
        LEAVE sp_body;
    END IF;
    CALL p_Guild_Member_Join(out_nNo,in_nMakerCharNo,0,v_nRet);
    IF v_nRet = 0 THEN
        ROLLBACK;
        SET out_nNo=-6;
        LEAVE sp_body;
    END IF;
    CALL usp_GuildAcademy_setAcademy(out_nNo,in_nMakerCharNo,50,v_nRet);
    IF v_nRet <> 0 THEN
        ROLLBACK;
        SET out_nNo=v_nRet;
        LEAVE sp_body;
    END IF;
    COMMIT;
    LEAVE sp_body;

    SELECT out_nNo AS nNo, out_nRemainCharMoney AS nRemainCharMoney;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Member_GetAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Member_GetAll`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Member_GetAll`(
    IN in_nNo INT
)
sp_body: BEGIN

    SELECT G.nCharNo,G.dJoinDate,G.nGrade,G.nFlags,G.nExpRatio,G.sIntro,C.sID,S.nClass,C.nLevel,C.sLoginZone,C.dLastLoginDate FROM tGuildMember AS G INNER JOIN tCharacter AS C ON C.nCharNo=G.nCharNo AND C.bDeleted=0 INNER JOIN tCharacterShape AS S ON S.nCharNo=G.nCharNo WHERE G.nNo=in_nNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Member_Join
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Member_Join`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Member_Join`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nGrade TINYINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    INSERT tGuildMember(nNo,nCharNo,nGrade) VALUES (in_nNo,in_nCharNo,in_nGrade);
    IF 0 <> 0 THEN
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Member_Leave
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Member_Leave`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Member_Leave`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    DELETE FROM tGuildMember WHERE nNo=in_nNo AND nCharNo=in_nCharNo;
    SET out_nRet=ROW_COUNT();

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Member_SetGrade
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Member_SetGrade`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Member_SetGrade`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_Member_SetIntro
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_Member_SetIntro`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_Member_SetIntro`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_sIntro VARCHAR(64),
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuildMember SET sIntro=in_sIntro WHERE nNo=in_nNo AND nCharNo=in_nCharNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_MoneyAdd
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_MoneyAdd`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_MoneyAdd`(
    IN in_nNo INT,
    IN in_nMemberCharNo INT,
    IN in_nMoneyAdd BIGINT,
    OUT out_nMemberMoney BIGINT,
    OUT out_nTotalMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    CALL p_Char_MoneyGet(in_nMemberCharNo,out_nMemberMoney,out_nRet);
    IF out_nRet <> 1 THEN
        LEAVE sp_body;
    END IF;
    IF in_nMoneyAdd > out_nMemberMoney THEN
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneySub(in_nMemberCharNo,in_nMoneyAdd,out_nRet);
    IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneyGet(in_nMemberCharNo,out_nMemberMoney,out_nRet);
    IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    UPDATE tGuild SET nMoney=nMoney+in_nMoneyAdd WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SELECT nMoney INTO out_nTotalMoney FROM tGuild WHERE nNo=in_nNo;
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    SET out_nRet=1;
    COMMIT;

    SELECT out_nMemberMoney AS nMemberMoney, out_nTotalMoney AS nTotalMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_MoneySub
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_MoneySub`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_MoneySub`(
    IN in_nNo INT,
    IN in_nMemberCharNo INT,
    IN in_nMoneySub BIGINT,
    OUT out_nMemberMoney BIGINT,
    OUT out_nTotalMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    SET out_nMemberMoney=0;
    SET out_nTotalMoney=0;
    SELECT nMoney INTO out_nTotalMoney FROM tGuild WHERE nNo=in_nNo;
    IF 0 <> 0 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    IF in_nMoneySub > out_nTotalMoney THEN
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    UPDATE tGuild SET nMoney=nMoney-in_nMoneySub WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneyAdd(in_nMemberCharNo,in_nMoneySub,out_nRet);
    IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=4;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneyGet(in_nMemberCharNo,out_nMemberMoney,out_nRet);
    IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    SELECT nMoney INTO out_nTotalMoney FROM tGuild WHERE nNo=in_nNo;
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET out_nRet=6;
        LEAVE sp_body;
    END IF;
    SET out_nRet=0;
    COMMIT;

    SELECT out_nMemberMoney AS nMemberMoney, out_nTotalMoney AS nTotalMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetDismiss
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetDismiss`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetDismiss`(
    IN in_nNo INT,
    IN in_nDismissStatus TINYINT,
    IN in_dDismissDate DATETIME,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET nDismissStatus=in_nDismissStatus,dDismissDate=in_dDismissDate WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetExp
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetExp`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetExp`(
    IN in_nNo INT,
    IN in_nExp BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET nExp=in_nExp WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetFame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetFame`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetFame`(
    IN in_nNo INT,
    IN in_nFame INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET nFame=in_nFame WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetIntro
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetIntro`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetIntro`(
    IN in_nNo INT,
    IN in_sIntro VARCHAR(128),
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET sIntro=in_sIntro WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetNotify
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetNotify`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetNotify`(
    IN in_nNo INT,
    IN in_sNotifyCharID VARCHAR(16),
    IN in_sNotify VARCHAR(512),
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET dNotifyDate=NOW(),sNotifyCharID=in_sNotifyCharID,sNotify=in_sNotify WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetType
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetType`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetType`(
    IN in_nNo INT,
    IN in_nType TINYINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET nType=in_nType WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetWar
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetWar`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetWar`(
    IN in_nNo INT,
    IN in_dWarRequestDate DATETIME,
    IN in_dWarStartDate DATETIME,
    IN in_dWarEndDate DATETIME,
    IN in_nWarEnemy INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET dWarRequestDate=in_dWarRequestDate,dWarStartDate=in_dWarStartDate,dWarEndDate=in_dWarEndDate,nWarEnemy=in_nWarEnemy WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_SetWarResult
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_SetWarResult`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_SetWarResult`(
    IN in_nNo INT,
    IN in_nWarWinCount INT,
    IN in_nWarLoseCount INT,
    IN in_nWarDrawCount INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    UPDATE tGuild SET nWarWinCount=in_nWarWinCount,nWarLoseCount=in_nWarLoseCount,nWarDrawCount=in_nWarDrawCount WHERE nNo=in_nNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Guild_setName
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Guild_setName`;
DELIMITER $$
CREATE PROCEDURE `p_Guild_setName`(
    IN in_nNo INT,
    IN in_sName VARCHAR(18),
    IN in_nRet INT
)
sp_body: BEGIN

    IF in_nNo <= 0 THEN
        SET in_nRet=-1311;
        LEAVE sp_body;
    END IF;
    IF CHAR_LENGTH ( IFNULL ( in_sName , '' ) ) = 0 THEN
        SET in_nRet=-1312;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nNo FROM tGuild WHERE nNo = in_nNo ) THEN
        IF EXISTS ( SELECT sName FROM tGuild WHERE sName = in_sName ) THEN
        SET in_nRet=-1322;
    ELSE
        START TRANSACTION;
        UPDATE tGuild SET sName=in_sName WHERE nNo=in_nNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-1393;
    ELSE
        COMMIT;
        SET in_nRet=0;
    END IF;
    END IF;
    ELSE
        SET in_nRet=-1321;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_Del
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_Del`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_Del`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tHolyPromise WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tHolyPromise SET nUpCharNo=0,nFlags=0,dRegDate=NOW() WHERE nCharNo=in_nCharNo;
        SET out_nRet=0;
    ELSE
        SET out_nRet=0;
    END IF;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_DelChar
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_DelChar`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_DelChar`(
    IN in_nDelCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tHolyPromise WHERE nUpCharNo=in_nDelCharNo;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_Get`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_Get`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT COUNT(nUpCharNo) INTO out_nRet FROM tHolyPromise WHERE nCharNo=in_nCharNo;
    IF out_nRet > 0 THEN
        IF EXISTS ( SELECT P . nUpCharNo , P . nFlags , P . dRegDate , P . dRejoinableDate , P . dMemberAcceptableDate , C . sID , C . nLevel , C . dLastLoginDate FROM tHolyPromise AS P INNER JOIN tCharacter AS C ON C . nCharNo = P . nUpCharNo AND C . bDeleted = 0 WHERE P . nCharNo = in_nCharNo ) THEN
        SELECT P.nUpCharNo,P.nFlags,P.dRegDate,P.dRejoinableDate,P.dMemberAcceptableDate,C.sID,C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P INNER JOIN tCharacter AS C ON C.nCharNo=P.nUpCharNo AND C.bDeleted=0 WHERE P.nCharNo=in_nCharNo;
    ELSE
        SELECT P.nUpCharNo,P.nFlags,P.dRegDate,P.dRejoinableDate,P.dMemberAcceptableDate,'',1,NOW() FROM tHolyPromise AS P WHERE P.nCharNo=in_nCharNo;
    END IF;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_GetMember
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_GetMember`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_GetMember`(
    IN in_nCharNo INT,
    OUT out_nNumOfMember INT
)
sp_body: BEGIN

    SELECT COUNT(P.nCharNo) INTO out_nNumOfMember FROM tHolyPromise AS P INNER JOIN tCharacter AS C ON C.nCharNo=P.nCharNo AND C.bDeleted=0 WHERE P.nUpCharNo=in_nCharNo;
    IF out_nNumOfMember > 0 THEN
        SELECT P.nCharNo,P.nFlags,P.dRegDate,P.dRejoinableDate,P.dMemberAcceptableDate,C.sID,C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P INNER JOIN tCharacter AS C ON C.nCharNo=P.nCharNo AND C.bDeleted=0 WHERE P.nUpCharNo=in_nCharNo;
    END IF;

    SELECT out_nNumOfMember AS nNumOfMember;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_Set`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_Set`(
    IN in_nMAXPROMISEMEMBER SMALLINT,
    IN in_nCharNo INT,
    IN in_nUpCharNo INT,
    IN in_nFlags TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nPromiseCount INT;
    DECLARE v_TempUpCharNo INT;

    CALL p_HolyPromise_CheckRing(in_nCharNo,in_nUpCharNo,out_nRet);
    IF out_nRet = 1 THEN
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_HolyPromise_GetMemberCount(in_nUpCharNo,v_nPromiseCount);
    IF v_nPromiseCount >= in_nMAXPROMISEMEMBER THEN
        SET out_nRet=2;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tHolyPromise WHERE nCharNo = in_nCharNo ) THEN
        SELECT nUpCharNo INTO v_TempUpCharNo FROM tHolyPromise WHERE nCharNo=in_nCharNo;
        IF v_TempUpCharNo <> 0 THEN
        SET out_nRet=4;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tHolyPromise SET nUpCharNo=in_nUpCharNo,nFlags=in_nFlags,dRegDate=NOW() WHERE nCharNo=in_nCharNo;
        SET out_nRet=0;
        IF out_nRet <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SET out_nRet=3;
    ELSE
        INSERT tHolyPromise(nCharNo,nUpCharNo,nFlags) VALUES (in_nCharNo,in_nUpCharNo,in_nFlags);
        SET out_nRet=0;
        IF out_nRet <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_HolyPromise_SetDate
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HolyPromise_SetDate`;
DELIMITER $$
CREATE PROCEDURE `p_HolyPromise_SetDate`(
    IN in_nCharNo INT,
    IN in_dRejoinableDate DATETIME,
    IN in_dMemberAcceptableDate DATETIME,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tHolyPromise WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tHolyPromise SET dRejoinableDate=in_dRejoinableDate,dMemberAcceptableDate=in_dMemberAcceptableDate WHERE nCharNo=in_nCharNo;
        SET out_nRet=0;
    ELSE
        INSERT tHolyPromise(nCharNo,nUpCharNo,nFlags,dRejoinableDate,dMemberAcceptableDate) VALUES (in_nCharNo,0,0,in_dRejoinableDate,in_dMemberAcceptableDate);
        SET out_nRet=0;
    END IF;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_BuyAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_BuyAll`;
DELIMITER $$
CREATE PROCEDURE `p_Item_BuyAll`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    IN in_nItemID INT,
    IN in_nFlags INT,
    IN in_nPrice BIGINT,
    IN in_nFame INT,
    IN in_nCoin BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_CreateItemRet INT;

    IF in_nPrice <= 0 AND in_nFame <= 0 AND in_nCoin <= 0 THEN
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_nPrice > 0 THEN
        CALL p_Char_MoneySub(in_nOwner,in_nPrice,out_nRet);
        IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    IF in_nFame > 0 THEN
        CALL p_Char_FameSub(in_nOwner,in_nFame,out_nRet);
        IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    IF in_nCoin > 0 THEN
        CALL usp_Gamble_subCoin(in_nOwner,in_nCoin,out_nRet);
        IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    CALL p_Item_Create(in_nItemKey,in_nOwner,in_nStorageType,in_nStorage,in_nItemID,in_nFlags,v_CreateItemRet);
    IF v_CreateItemRet <> 0 OR 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=v_CreateItemRet;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_BuyLot
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_BuyLot`;
DELIMITER $$
CREATE PROCEDURE `p_Item_BuyLot`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nOptionType SMALLINT,
    IN in_nOptionData BIGINT,
    IN in_nPrice BIGINT,
    IN in_nFame INT,
    IN in_nCoin BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF in_nPrice <= 0 AND in_nFame <= 0 AND in_nCoin <= 0 THEN
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_nPrice > 0 THEN
        CALL p_Char_MoneySub(in_nOwner,in_nPrice,out_nRet);
        IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    IF in_nFame > 0 THEN
        CALL p_Char_FameSub(in_nOwner,in_nFame,out_nRet);
        IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    IF in_nCoin > 0 THEN
        CALL usp_Gamble_subCoin(in_nOwner,in_nCoin,out_nRet);
        IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    END IF;
    CALL p_Item_SetOption(in_nItemKey,in_nOptionType,in_nOptionData,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_ClrWeaponTitle
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_ClrWeaponTitle`;
DELIMITER $$
CREATE PROCEDURE `p_Item_ClrWeaponTitle`(
    IN in_nWeaponItemKey BIGINT,
    IN in_nMobID SMALLINT,
    IN in_nLicenseClearItemKey BIGINT,
    IN in_bLicenseClearItemDelete TINYINT,
    IN in_nLicenseClearItemOptionType SMALLINT,
    IN in_nLicenseClearItemOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nRetTemp INT;

    IF NOT EXISTS ( SELECT nMobID FROM tItemMobList WHERE nItemKey = in_nWeaponItemKey AND nMobID = in_nMobID ) THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    DELETE FROM tItemMobList WHERE nItemKey=in_nWeaponItemKey AND nMobID=in_nMobID;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    IF in_bLicenseClearItemDelete = 1 THEN
        CALL p_Item_Delete(in_nLicenseClearItemKey,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    ELSE
        CALL p_Item_SetOption(in_nLicenseClearItemKey,in_nLicenseClearItemOptionType,in_nLicenseClearItemOptionData,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=4;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_CostumWeapon_Endure_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_CostumWeapon_Endure_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Item_CostumWeapon_Endure_Set`(
    IN in_nCostumWeaponKey BIGINT,
    IN in_nCostumWeaponEndureOptionType SMALLINT,
    IN in_nCostumWeaponEndureOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nRetTemp INT;

    START TRANSACTION;
    SET v_nRetTemp=1;
    CALL p_Item_SetOption(in_nCostumWeaponKey,in_nCostumWeaponEndureOptionType,in_nCostumWeaponEndureOptionData,v_nRetTemp);
    IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_CostumWeapon_Endure_Update
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_CostumWeapon_Endure_Update`;
DELIMITER $$
CREATE PROCEDURE `p_Item_CostumWeapon_Endure_Update`(
    IN in_nOwner INT,
    IN in_nEndureKitKey BIGINT,
    IN in_nEndureKitOptionType SMALLINT,
    IN in_nEndureKitLot SMALLINT,
    IN in_nCostumWeaponKey BIGINT,
    IN in_nCostumWeaponEndureOptionType SMALLINT,
    IN in_nCostumWeaponEndureOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nRetTemp INT;

    CALL p_Item_IsValidOwner(in_nEndureKitKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    CALL p_Item_IsValidOwner(in_nCostumWeaponKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_nEndureKitLot = 0 THEN
        SET v_nRetTemp=1;
        CALL p_Item_Delete(in_nEndureKitKey,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    ELSE
        SET v_nRetTemp=1;
        CALL p_Item_SetOption(in_nEndureKitKey,in_nEndureKitOptionType,in_nEndureKitLot,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    END IF;
    SET v_nRetTemp=1;
    CALL p_Item_SetOption(in_nCostumWeaponKey,in_nCostumWeaponEndureOptionType,in_nCostumWeaponEndureOptionData,v_nRetTemp);
    IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Create
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Create`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Create`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Create_UseEmptyStorage
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Create_UseEmptyStorage`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Create_UseEmptyStorage`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    OUT out_nStorage SMALLINT,
    IN in_nItemID INT,
    IN in_nFlags INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    START TRANSACTION;
    CALL p_Item_GetStorageEmpty(in_nOwner,in_nStorageType,out_nStorage);
    IF out_nStorage < 0 THEN
        SET out_nRet=-1;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    INSERT tItem(nItemKey,nOwner,nStorageType,nStorage,nItemID,nFlags) VALUES (in_nItemKey,in_nOwner,in_nStorageType,out_nStorage,in_nItemID,in_nFlags);
    SET out_nRet=0;
    IF out_nRet = 0 THEN
        COMMIT;
    ELSE
        ROLLBACK;
    END IF;

    SELECT out_nStorage AS nStorage, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Delete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Delete`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Delete`(
    IN in_nItemKey BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN
    DELETE FROM tItem WHERE nItemKey = in_nItemKey;
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_ExchangeStorage
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_ExchangeStorage`;
DELIMITER $$
CREATE PROCEDURE `p_Item_ExchangeStorage`(
    IN in_nItemKey1 BIGINT,
    IN in_nStorageType1 TINYINT,
    IN in_nStorage1 SMALLINT,
    IN in_nItemKey2 BIGINT,
    IN in_nStorageType2 TINYINT,
    IN in_nStorage2 SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    UPDATE tItem SET nStorageType=11 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=1;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType2,nStorage=in_nStorage2 WHERE nItemKey=in_nItemKey2;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=2;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType1,nStorage=in_nStorage1 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=3;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_ExchangeStorageAndOwner
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_ExchangeStorageAndOwner`;
DELIMITER $$
CREATE PROCEDURE `p_Item_ExchangeStorageAndOwner`(
    IN in_nItemKey1 BIGINT,
    IN in_nOwner1 INT,
    IN in_nStorageType1 TINYINT,
    IN in_nStorage1 SMALLINT,
    IN in_nItemKey2 BIGINT,
    IN in_nOwner2 INT,
    IN in_nStorageType2 TINYINT,
    IN in_nStorage2 SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    UPDATE tItem SET nStorageType=11 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=1;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType2,nStorage=in_nStorage2,nOwner=in_nOwner2 WHERE nItemKey=in_nItemKey2;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=2;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType1,nStorage=in_nStorage1,nOwner=in_nOwner1 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=3;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Furniture_Endure_Update
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Furniture_Endure_Update`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Furniture_Endure_Update`(
    IN in_nOwner INT,
    IN in_nEndureKitKey BIGINT,
    IN in_nEndureKitOptionType SMALLINT,
    IN in_nEndureKitLot SMALLINT,
    IN in_nFurnitureKey BIGINT,
    IN in_nFurnitureEndureOptionType SMALLINT,
    IN in_nFurnitureEndureOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nRetTemp INT;

    CALL p_Item_IsValidOwner(in_nEndureKitKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    CALL p_Item_IsValidOwner(in_nFurnitureKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_nEndureKitLot = 0 THEN
        SET v_nRetTemp=1;
        CALL p_Item_Delete(in_nEndureKitKey,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    ELSE
        SET v_nRetTemp=1;
        CALL p_Item_SetOption(in_nEndureKitKey,in_nEndureKitOptionType,in_nEndureKitLot,v_nRetTemp);
        IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    END IF;
    SET v_nRetTemp=1;
    CALL p_Item_SetOption(in_nFurnitureKey,in_nFurnitureEndureOptionType,in_nFurnitureEndureOptionData,v_nRetTemp);
    IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetData
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetData`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetData`(
    IN in_nItemKey BIGINT,
    OUT out_nOwner INT,
    OUT out_nStorageType TINYINT,
    OUT out_nStorage SMALLINT,
    OUT out_nItemID INT,
    OUT out_nFlags INT,
    OUT out_dDate DATETIME,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT nOwner,nStorageType,nStorage,nItemID,nFlags,dDate INTO out_nOwner,out_nStorageType,out_nStorage,out_nItemID,out_nFlags,out_dDate FROM tItem WHERE nItemKey=in_nItemKey;
    SET out_nRet=ROW_COUNT();

    SELECT out_nOwner AS nOwner, out_nStorageType AS nStorageType, out_nStorage AS nStorage, out_nItemID AS nItemID, out_nFlags AS nFlags, out_dDate AS dDate, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetListType
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetListType`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetListType`(
    IN in_nCharNo INT,
    IN in_nStorageType TINYINT
)
sp_body: BEGIN

    SELECT nItemKey,nStorage,nItemID,nFlags,dDate FROM tItem WHERE nOwner=in_nCharNo AND nStorageType=in_nStorageType;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetMobListAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetMobListAll`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetMobListAll`(
    IN in_nItemKey BIGINT
)
sp_body: BEGIN

    SELECT nMobID,nGrade,nKillCount FROM tItemMobList WHERE nItemKey=in_nItemKey;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetOption`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetOption`(
    IN in_nItemKey BIGINT,
    IN in_nOptionType SMALLINT,
    OUT out_nOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT nOptionData INTO out_nOptionData FROM tItemOptions WHERE nItemKey=in_nItemKey AND nOptionType=in_nOptionType;
    SET out_nRet=ROW_COUNT();

    SELECT out_nOptionData AS nOptionData, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetOptionAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetOptionAll`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetOptionAll`(
    IN in_nItemKey BIGINT
)
sp_body: BEGIN

    SELECT nOptionType,nOptionData FROM tItemOptions WHERE nItemKey=in_nItemKey;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetStorageAndOwner
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetStorageAndOwner`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetStorageAndOwner`(
    IN in_nItemKey BIGINT,
    OUT out_nOwner INT,
    OUT out_nStorageType TINYINT,
    OUT out_nStorage SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SELECT nOwner,nStorageType,nStorage INTO out_nOwner,out_nStorageType,out_nStorage FROM tItem WHERE nItemKey=in_nItemKey;
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=0;
    ELSE
        SET out_nRet=1;
    END IF;

    SELECT out_nOwner AS nOwner, out_nStorageType AS nStorageType, out_nStorage AS nStorage, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetStorageEmptyList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetStorageEmptyList`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetStorageEmptyList`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetStorageListType
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetStorageListType`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetStorageListType`(
    IN in_nCharNo INT,
    IN in_nStorageType TINYINT
)
sp_body: BEGIN

    SELECT nStorage FROM tItem WHERE nOwner=in_nCharNo AND nStorageType=in_nStorageType;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_GetTitle
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_GetTitle`;
DELIMITER $$
CREATE PROCEDURE `p_Item_GetTitle`(
    IN in_nItemKey BIGINT,
    OUT out_nTitleType SMALLINT,
    OUT out_nTitleCharNo INT,
    OUT out_sTitle VARCHAR(24),
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT nTitleType,nTitleCharNo,sTitle INTO out_nTitleType,out_nTitleCharNo,out_sTitle FROM tItemTitle WHERE nItemKey=in_nItemKey;
    SET out_nRet=ROW_COUNT();

    SELECT out_nTitleType AS nTitleType, out_nTitleCharNo AS nTitleCharNo, out_sTitle AS sTitle, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_IsValidOwner
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_IsValidOwner`;
DELIMITER $$
CREATE PROCEDURE `p_Item_IsValidOwner`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nDBOwner INT;

    SET v_nDBOwner=0;
    SET out_nRet=0;
    SELECT nOwner INTO v_nDBOwner FROM tItem WHERE nItemKey=in_nItemKey;
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    IF v_nDBOwner <> in_nOwner THEN
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Merge_New
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Merge_New`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Merge_New`(
    IN in_nOwnerFrom INT,
    IN in_nOwnerTo INT,
    IN in_nItemKeyFrom BIGINT,
    IN in_nItemKeyTo BIGINT,
    IN in_nOptionType SMALLINT,
    OUT out_nOptionDataLotSub BIGINT,
    OUT out_nOptionDataLotResult BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nOptionDataLotAdd BIGINT;

    CALL p_Item_IsValidOwner(in_nItemKeyFrom,in_nOwnerFrom,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    CALL p_Item_IsValidOwner(in_nItemKeyTo,in_nOwnerTo,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet+2;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nItemKeyFrom,in_nOptionType,out_nOptionDataLotSub,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nItemKeyTo,in_nOptionType,v_nOptionDataLotAdd,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=6;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SET out_nOptionDataLotResult=out_nOptionDataLotSub+v_nOptionDataLotAdd;
    CALL p_Item_SetOption(in_nItemKeyTo,in_nOptionType,out_nOptionDataLotResult,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=7;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nItemKeyFrom,out_nRet);
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=8;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nOptionDataLotSub AS nOptionDataLotSub, out_nOptionDataLotResult AS nOptionDataLotResult, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SellAll_LotCheck
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SellAll_LotCheck`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SellAll_LotCheck`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nOptionType SMALLINT,
    IN in_nLot BIGINT,
    IN in_nPrice BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nGetOptionDataLot BIGINT;
    DECLARE v_nGetOptionRet INT;

    CALL p_Item_IsValidOwner(in_nItemKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    IF in_nLot > 1 THEN
        CALL p_Item_GetOption(in_nItemKey,in_nOptionType,v_nGetOptionDataLot,v_nGetOptionRet);
        IF v_nGetOptionRet <> 1 THEN
        SET out_nRet=4;
        LEAVE sp_body;
    END IF;
        IF v_nGetOptionDataLot < in_nLot THEN
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneyAdd(in_nOwner,in_nPrice,v_nTinyintRet);
    IF v_nTinyintRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nItemKey,out_nRet);
    IF 0 <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SellLot_LotCheck
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SellLot_LotCheck`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SellLot_LotCheck`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nOptionType SMALLINT,
    IN in_nOptionData BIGINT,
    IN in_nPrice BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nGetOptionDataLot BIGINT;
    DECLARE v_nGetOptionRet INT;

    CALL p_Item_IsValidOwner(in_nItemKey,in_nOwner,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nItemKey,in_nOptionType,v_nGetOptionDataLot,v_nGetOptionRet);
    IF v_nGetOptionRet <> 1 THEN
        SET out_nRet=4;
        LEAVE sp_body;
    END IF;
    IF v_nGetOptionDataLot < in_nOptionData THEN
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneyAdd(in_nOwner,in_nPrice,out_nRet);
    IF out_nRet <> 1 THEN
        ROLLBACK;
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nItemKey,in_nOptionType,in_nOptionData,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetMobList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetMobList`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetMobList`(
    IN in_nItemKey BIGINT,
    IN in_nMobID SMALLINT,
    IN in_nGrade TINYINT,
    IN in_nKillCount INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tItemMobList SET nGrade=in_nGrade,nKillCount=in_nKillCount WHERE nItemKey=in_nItemKey AND nMobID=in_nMobID;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tItemMobList(nItemKey,nMobID,nGrade,nKillCount) VALUES (in_nItemKey,in_nMobID,in_nGrade,in_nKillCount);
    END IF;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetOption`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetOption`(
    IN in_nItemKey BIGINT,
    IN in_nOptionType SMALLINT,
    IN in_nOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tItemOptions SET nOptionData=in_nOptionData WHERE nItemKey=in_nItemKey AND nOptionType=in_nOptionType;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tItemOptions(nItemKey,nOptionType,nOptionData) VALUES (in_nItemKey,in_nOptionType,in_nOptionData);
    END IF;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetStorage
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetStorage`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetStorage`(
    IN in_nItemKey BIGINT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    UPDATE tItem SET nStorageType=in_nStorageType,nStorage=in_nStorage WHERE nItemKey=in_nItemKey;
    SET out_nRet=ROW_COUNT();

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetStorageAndOwner
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetStorageAndOwner`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetStorageAndOwner`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    UPDATE tItem SET nOwner=in_nOwner,nStorageType=in_nStorageType,nStorage=in_nStorage WHERE nItemKey=in_nItemKey;
    SET out_nRet=ROW_COUNT();

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetStorageAndOwner_UseEmptyStorage
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetStorageAndOwner_UseEmptyStorage`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetStorageAndOwner_UseEmptyStorage`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    OUT out_nStorage SMALLINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    CALL p_Item_GetStorageEmpty(in_nOwner,in_nStorageType,out_nStorage);
    IF out_nStorage < 0 THEN
        SET out_nRet=-1;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nOwner=in_nOwner,nStorageType=in_nStorageType,nStorage=out_nStorage WHERE nItemKey=in_nItemKey;
    SET out_nRet=ROW_COUNT();

    SELECT out_nStorage AS nStorage, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetTitle
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetTitle`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetTitle`(
    IN in_nItemKey BIGINT,
    IN in_nTitleType SMALLINT,
    IN in_nTitleCharNo INT,
    IN in_sTitle VARCHAR(24),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tItemTitle SET nTitleType=in_nTitleType,nTitleCharNo=in_nTitleCharNo,sTitle=in_sTitle WHERE nItemKey=in_nItemKey;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tItemTitle(nItemKey,nTitleType,nTitleCharNo,sTitle) VALUES (in_nItemKey,in_nTitleType,in_nTitleCharNo,in_sTitle);
    END IF;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_SetWeaponTitle
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_SetWeaponTitle`;
DELIMITER $$
CREATE PROCEDURE `p_Item_SetWeaponTitle`(
    IN in_nWeaponItemKey BIGINT,
    IN in_nLicenseItemKey BIGINT,
    IN in_nMobID SMALLINT,
    IN in_nTitleCharNo INT,
    IN in_sTitle VARCHAR(24),
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nRetTemp INT;

    START TRANSACTION;
    CALL p_Item_SetMobList(in_nWeaponItemKey,in_nMobID,0,0,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    SELECT COUNT(*) INTO v_nRetTemp FROM tItemMobList WHERE nItemKey=in_nWeaponItemKey;
    IF v_nRetTemp = 1 THEN
        CALL p_Item_SetTitle(in_nWeaponItemKey,0,in_nTitleCharNo,in_sTitle,out_nRet);
        IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    END IF;
    SET v_nRetTemp=1;
    CALL p_Item_Delete(in_nLicenseItemKey,v_nRetTemp);
    IF v_nRetTemp <> 0 THEN
        ROLLBACK;
        SET out_nRet=3;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Split_N_Merge_New
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Split_N_Merge_New`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Split_N_Merge_New`(
    IN in_nItemKeyFrom BIGINT,
    IN in_nOwnerFrom INT,
    IN in_nOptionTypeFrom SMALLINT,
    OUT out_nOptionDataFrom BIGINT,
    IN in_nItemKeyTo BIGINT,
    IN in_nOwnerTo INT,
    IN in_nOptionTypeTo SMALLINT,
    OUT out_nOptionDataTo BIGINT,
    IN in_nSplitLot BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTinyintRet TINYINT;
    DECLARE v_nOptionDataLotAdd BIGINT;

    CALL p_Item_IsValidOwner(in_nItemKeyFrom,in_nOwnerFrom,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet;
        LEAVE sp_body;
    END IF;
    CALL p_Item_IsValidOwner(in_nItemKeyTo,in_nOwnerTo,v_nTinyintRet);
    IF v_nTinyintRet <> 0 THEN
        SET out_nRet=v_nTinyintRet+2;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nItemKeyFrom,in_nOptionTypeFrom,out_nOptionDataFrom,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=5;
        LEAVE sp_body;
    END IF;
    IF out_nOptionDataFrom < in_nSplitLot THEN
        SET out_nRet=6;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nItemKeyTo,in_nOptionTypeTo,out_nOptionDataTo,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=7;
        LEAVE sp_body;
    END IF;
    SET out_nOptionDataFrom=out_nOptionDataFrom-in_nSplitLot;
    SET out_nOptionDataTo=out_nOptionDataTo+in_nSplitLot;
    START TRANSACTION;
    CALL p_Item_SetOption(in_nItemKeyFrom,in_nOptionTypeFrom,out_nOptionDataFrom,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=8;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nItemKeyTo,in_nOptionTypeTo,out_nOptionDataTo,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=9;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nOptionDataFrom AS nOptionDataFrom, out_nOptionDataTo AS nOptionDataTo, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Item_Upgrade_CelllInAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Item_Upgrade_CelllInAll`;
DELIMITER $$
CREATE PROCEDURE `p_Item_Upgrade_CelllInAll`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_Add`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_Add`(
    IN in_sGmUserID VARCHAR(16),
    IN in_sImprison VARCHAR(16),
    IN in_nMinute SMALLINT,
    IN in_sReason VARCHAR(16),
    IN in_sRemark VARCHAR(64),
    IN in_nImprisonRegNum INT,
    OUT out_err SMALLINT
)
sp_body: BEGIN

    DECLARE v_nRet INT;
    DECLARE v_nPrisonMin INT;

    SET out_err=0;
    START TRANSACTION;
    UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=896,nLoginZoneY=1122,sLoginZone='EldPri' WHERE nCharNo=in_nImprisonRegNum;
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    INSERT tPrison(nCharID,nMinute,sReason,sGmID,sRemark) VALUES (in_nImprisonRegNum,in_nMinute,in_sReason,in_sGmUserID,in_sRemark);
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    LEAVE sp_body;

    SELECT out_err AS err;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_Add2
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_Add2`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_Add2`(
    IN in_sGmUserID VARCHAR(16),
    IN in_sImprison VARCHAR(16),
    IN in_nMinute SMALLINT,
    IN in_sReason VARCHAR(16),
    IN in_sRemark VARCHAR(64),
    OUT out_err SMALLINT
)
sp_body: BEGIN

    DECLARE v_nRet INT;
    DECLARE v_nPrisonMin INT;
    DECLARE v_nImprisonRegNum INT;

    SET out_err=0;
    SELECT nCharNo,nPrisonMin INTO v_nImprisonRegNum,v_nPrisonMin FROM tCharacter WHERE sID=in_sImprison;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_err=1;
        LEAVE sp_body;
    END IF;
    IF v_nPrisonMin <> 0 THEN
        SET out_err=2;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=896,nLoginZoneY=1122,sLoginZone='EldPri' WHERE nCharNo=v_nImprisonRegNum;
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    INSERT tPrison(nCharID,nMinute,sReason,sGmID,sRemark) VALUES (v_nImprisonRegNum,in_nMinute,in_sReason,in_sGmUserID,in_sRemark);
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    LEAVE sp_body;

    SELECT out_err AS err;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_Alter
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_Alter`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_Alter`(
    IN in_sGmUserID VARCHAR(16),
    IN in_sImprison VARCHAR(16),
    IN in_nMinute SMALLINT,
    IN in_sReason VARCHAR(16),
    IN in_sRemark VARCHAR(64),
    IN in_nImprisonRegNum INT,
    OUT out_err SMALLINT
)
sp_body: BEGIN

    DECLARE v_nRet INT;
    DECLARE v_nPrisonMin INT;

    SET out_err=0;
    START TRANSACTION;
    IF in_nMinute = 0 THEN
        UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=5963,nLoginZoneY=9793,sLoginZone='Rou' WHERE nCharNo=in_nImprisonRegNum;
    ELSE
        UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=896,nLoginZoneY=1122,sLoginZone='EldPri' WHERE nCharNo=in_nImprisonRegNum;
    END IF;
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    INSERT tPrison(nCharID,nMinute,sReason,sGmID,sRemark) VALUES (in_nImprisonRegNum,in_nMinute,in_sReason,in_sGmUserID,in_sRemark);
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    LEAVE sp_body;

    SELECT out_err AS err;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_Alter2
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_Alter2`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_Alter2`(
    IN in_sGmUserID VARCHAR(16),
    IN in_sImprison VARCHAR(16),
    IN in_nMinute SMALLINT,
    IN in_sReason VARCHAR(16),
    IN in_sRemark VARCHAR(64),
    OUT out_err SMALLINT
)
sp_body: BEGIN

    DECLARE v_nRet INT;
    DECLARE v_nPrisonMin INT;
    DECLARE v_nImprisonRegNum INT;

    SET out_err=0;
    SELECT nCharNo,nPrisonMin INTO v_nImprisonRegNum,v_nPrisonMin FROM tCharacter WHERE sID=in_sImprison;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_err=1;
        LEAVE sp_body;
    END IF;
    IF v_nPrisonMin = 0 THEN
        SET out_err=2;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_nMinute = 0 THEN
        UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=5963,nLoginZoneY=9793,sLoginZone='Rou' WHERE nCharNo=v_nImprisonRegNum;
    ELSE
        UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=896,nLoginZoneY=1122,sLoginZone='EldPri' WHERE nCharNo=v_nImprisonRegNum;
    END IF;
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    INSERT tPrison(nCharID,nMinute,sReason,sGmID,sRemark) VALUES (v_nImprisonRegNum,in_nMinute,in_sReason,in_sGmUserID,in_sRemark);
    IF 0 <> 0 THEN
        ROLLBACK;
    END IF;
    SET out_err=1;
    LEAVE sp_body;
    COMMIT;
    LEAVE sp_body;
    ROLLBACK;
    LEAVE sp_body;

    SELECT out_err AS err;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_Get`(
    IN in_nCharID INT
)
sp_body: BEGIN

    SELECT sReason,sRemark FROM tPrison WHERE nCharID=in_nCharID ORDER BY dPrisonDate DESC LIMIT 1;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Prison_UpdateCharPrisonMin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Prison_UpdateCharPrisonMin`;
DELIMITER $$
CREATE PROCEDURE `p_Prison_UpdateCharPrisonMin`(
    IN in_nImprisonRegNum INT,
    IN in_nMinute SMALLINT
)
sp_body: BEGIN

    IF in_nMinute = 0 THEN
        UPDATE tCharacter SET nPrisonMin=in_nMinute,nLoginZoneX=5963,nLoginZoneY=9793,sLoginZone='Rou' WHERE nCharNo=in_nImprisonRegNum;
    ELSE
        UPDATE tCharacter SET nPrisonMin=in_nMinute,sLoginZone='EldPri' WHERE nCharNo=in_nImprisonRegNum;
    END IF;
    LEAVE sp_body;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Quest_Del
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Quest_Del`;
DELIMITER $$
CREATE PROCEDURE `p_Quest_Del`(
    IN in_nCharNo INT,
    IN in_nQuestNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tQuest WHERE nCharNo=in_nCharNo AND nQuestNo=in_nQuestNo;
    IF ( 0 = 0 ) AND ( ROW_COUNT ( ) = 1 ) THEN
        SET out_nRet=0;
    ELSE
        SET out_nRet=1;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Quest_GetAllDoing
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Quest_GetAllDoing`;
DELIMITER $$
CREATE PROCEDURE `p_Quest_GetAllDoing`(
    IN in_nCharNo INT,
    IN in_nStatusDone TINYINT
)
sp_body: BEGIN

    SELECT nQuestNo,nStatus,sData FROM tQuest WHERE nCharNo=in_nCharNo AND nStatus=in_nStatusDone;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Quest_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Quest_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Quest_Set`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_Del
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_Del`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_Del`(
    IN in_nCharNo INT,
    IN in_nSkillNo SMALLINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DELETE FROM tSkill WHERE nCharNo=in_nCharNo AND nSkillNo=in_nSkillNo;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_Get`(
    IN in_nCharNo INT,
    IN in_nSkillNo SMALLINT,
    OUT out_nSkillLevel TINYINT,
    OUT out_nSkillExp INT,
    OUT out_nSkillWriteTime BIGINT,
    OUT out_nSkillCoolTime INT,
    OUT out_nSkillPowerDemage TINYINT,
    OUT out_nSkillPowerSP TINYINT,
    OUT out_nSkillPowerKeepTime TINYINT,
    OUT out_nSkillPowerCoolTime TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT nSkillLevel,nSkillExp,nSkillWriteTime,nSkillCoolTime,nSkillPowerDemage,nSkillPowerSP,nSkillPowerKeepTime,nSkillPowerCoolTime INTO out_nSkillLevel,out_nSkillExp,out_nSkillWriteTime,out_nSkillCoolTime,out_nSkillPowerDemage,out_nSkillPowerSP,out_nSkillPowerKeepTime,out_nSkillPowerCoolTime FROM tSkill WHERE nCharNo=in_nCharNo AND nSkillNo=in_nSkillNo;
    SET out_nRet=ROW_COUNT();

    SELECT out_nSkillLevel AS nSkillLevel, out_nSkillExp AS nSkillExp, out_nSkillWriteTime AS nSkillWriteTime, out_nSkillCoolTime AS nSkillCoolTime, out_nSkillPowerDemage AS nSkillPowerDemage, out_nSkillPowerSP AS nSkillPowerSP, out_nSkillPowerKeepTime AS nSkillPowerKeepTime, out_nSkillPowerCoolTime AS nSkillPowerCoolTime, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_GetAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_GetAll`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_GetAll`(
    IN in_nCharNo INT
)
sp_body: BEGIN

    SELECT nSkillNo,nSkillLevel,nSkillExp,nSkillWriteTime,nSkillCoolTime,nSkillPowerDemage,nSkillPowerSP,nSkillPowerKeepTime,nSkillPowerCoolTime FROM tSkill WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_ResetPower
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_ResetPower`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_ResetPower`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nFlagsTemp INT;

    SET out_nRet=0;
    SELECT nFlags INTO v_nFlagsTemp FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) = 0 THEN
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SET v_nFlagsTemp=v_nFlagsTemp & 0x7fffffff;
    UPDATE tCharacter SET nFlags=v_nFlagsTemp WHERE nCharNo=in_nCharNo;
    UPDATE tSkill SET nSkillPowerDemage=0,nSkillPowerSP=0,nSkillPowerKeepTime=0,nSkillPowerCoolTime=0 WHERE nCharNo=in_nCharNo;
    IF 0 <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=1;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_Set`(
    IN in_nCharNo INT,
    IN in_nSkillNo SMALLINT,
    IN in_nSkillLevel TINYINT,
    IN in_nSkillExp INT,
    IN in_nSkillWriteTime BIGINT,
    IN in_nSkillCoolTime INT,
    IN in_nSkillPowerDemage TINYINT,
    IN in_nSkillPowerSP TINYINT,
    IN in_nSkillPowerKeepTime TINYINT,
    IN in_nSkillPowerCoolTime TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tSkill SET nSkillLevel=in_nSkillLevel,nSkillExp=in_nSkillExp,nSkillWriteTime=in_nSkillWriteTime,nSkillCoolTime=in_nSkillCoolTime,nSkillPowerDemage=in_nSkillPowerDemage,nSkillPowerSP=in_nSkillPowerSP,nSkillPowerKeepTime=in_nSkillPowerKeepTime,nSkillPowerCoolTime=in_nSkillPowerCoolTime WHERE nCharNo=in_nCharNo AND nSkillNo=in_nSkillNo;
    IF ROW_COUNT ( ) = 0 THEN
        INSERT tSkill(nCharNo,nSkillNo,nSkillLevel,nSkillExp,nSkillWriteTime,nSkillCoolTime,nSkillPowerDemage,nSkillPowerSP,nSkillPowerKeepTime,nSkillPowerCoolTime) VALUES (in_nCharNo,in_nSkillNo,in_nSkillLevel,in_nSkillExp,in_nSkillWriteTime,in_nSkillCoolTime,in_nSkillPowerDemage,in_nSkillPowerSP,in_nSkillPowerKeepTime,in_nSkillPowerCoolTime);
    END IF;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_SetPassive
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_SetPassive`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_SetPassive`(
    IN in_nCharNo INT,
    IN in_nSkillNewNo SMALLINT,
    IN in_bSkillOld TINYINT,
    IN in_nSkillOldNo SMALLINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF in_bSkillOld = 1 THEN
        UPDATE tSkill SET nSkillNo=in_nSkillNewNo WHERE nCharNo=in_nCharNo AND nSkillNo=in_nSkillOldNo;
        IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=1;
    ELSE
        SET out_nRet=0;
    END IF;
        LEAVE sp_body;
    END IF;
    INSERT tSkill(nCharNo,nSkillNo) VALUES (in_nCharNo,in_nSkillNewNo);
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_Skill_SetPowerAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Skill_SetPowerAll`;
DELIMITER $$
CREATE PROCEDURE `p_Skill_SetPowerAll`(
    IN in_nCharNo INT,
    IN in_nSkillNo SMALLINT,
    IN in_nSkillPowerDemage TINYINT,
    IN in_nSkillPowerSP TINYINT,
    IN in_nSkillPowerKeepTime TINYINT,
    IN in_nSkillPowerCoolTime TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    IF in_nSkillPowerDemage > 5 OR in_nSkillPowerSP > 5 OR in_nSkillPowerKeepTime > 5 OR in_nSkillPowerCoolTime > 5 THEN
        SET out_nRet=0;
        LEAVE sp_body;
    END IF;
    UPDATE tSkill SET nSkillPowerDemage=in_nSkillPowerDemage,nSkillPowerSP=in_nSkillPowerSP,nSkillPowerKeepTime=in_nSkillPowerKeepTime,nSkillPowerCoolTime=in_nSkillPowerCoolTime WHERE nCharNo=in_nCharNo AND nSkillNo=in_nSkillNo;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=0;
    ELSE
        SET out_nRet=1;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_User_Deposit
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_User_Deposit`;
DELIMITER $$
CREATE PROCEDURE `p_User_Deposit`(
    IN in_nUserNo INT,
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    OUT out_nResultUserMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nDBUserNo INT;

    SET out_nRet=0;
    SET v_nDBUserNo=0;
    SELECT nUserNo INTO v_nDBUserNo FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF v_nDBUserNo <> in_nUserNo THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneySub(in_nCharNo,in_nMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=2;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    CALL p_User_MoneyAdd(in_nUserNo,in_nMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=3;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    CALL p_User_MoneyGet(in_nUserNo,out_nResultUserMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=4;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nResultUserMoney AS nResultUserMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_User_MoneyGet
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_User_MoneyGet`;
DELIMITER $$
CREATE PROCEDURE `p_User_MoneyGet`(
    IN in_nUserNo INT,
    OUT out_nUserMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    SET out_nRet=0;
    IF EXISTS ( SELECT nUserMoney FROM tUserMoney WHERE nUserNo = in_nUserNo ) THEN
        SELECT nUserMoney INTO out_nUserMoney FROM tUserMoney WHERE nUserNo=in_nUserNo;
        IF ROW_COUNT ( ) = 1 THEN
        SET out_nRet=1;
    END IF;
        LEAVE sp_body;
    END IF;
    SET out_nUserMoney=0;
    INSERT INTO tUserMoney VALUES (in_nUserNo,out_nUserMoney);
    IF ROW_COUNT ( ) = 1 AND 0 = 0 THEN
        SET out_nRet=1;
    END IF;
    LEAVE sp_body;

    SELECT out_nUserMoney AS nUserMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: p_User_Withdraw
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_User_Withdraw`;
DELIMITER $$
CREATE PROCEDURE `p_User_Withdraw`(
    IN in_nUserNo INT,
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    OUT out_nResultUserMoney BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nDBUserNo INT;

    SET out_nRet=0;
    SET v_nDBUserNo=0;
    SELECT nUserNo INTO v_nDBUserNo FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF v_nDBUserNo <> in_nUserNo THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    CALL p_User_MoneySub(in_nUserNo,in_nMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=2;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    CALL p_Char_MoneyAdd(in_nCharNo,in_nMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=3;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    CALL p_User_MoneyGet(in_nUserNo,out_nResultUserMoney,out_nRet);
    IF out_nRet <> 1 THEN
        SET out_nRet=4;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nResultUserMoney AS nResultUserMoney, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Attendance_getCallDateList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Attendance_getCallDateList`;
DELIMITER $$
CREATE PROCEDURE `usp_Attendance_getCallDateList`(
    IN in_nCharNo INT,
    IN in_dStartDate DATETIME,
    IN in_dEndDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT dAttendanceDate FROM tAttendance WHERE nCharNo=in_nCharNo AND dAttendanceDate>=in_dStartDate AND dAttendanceDate<=in_dEndDate;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Attendance_getCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Attendance_getCount`;
DELIMITER $$
CREATE PROCEDURE `usp_Attendance_getCount`(
    IN in_nCharNo INT,
    IN in_dStartDate DATETIME,
    IN in_dEndDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT COUNT(dAttendanceDate) AS nAttendanceCount FROM tAttendance WHERE nCharNo=in_nCharNo AND dAttendanceDate>=in_dStartDate AND dAttendanceDate<=in_dEndDate;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Attendance_getRewardList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Attendance_getRewardList`;
DELIMITER $$
CREATE PROCEDURE `usp_Attendance_getRewardList`(
    IN in_nCharNo INT,
    IN in_dDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT dEventDate,nRewardID FROM tAttendanceReward WHERE nCharNo=in_nCharNo AND dEventDate>in_dDate;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Attendance_setCall
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Attendance_setCall`;
DELIMITER $$
CREATE PROCEDURE `usp_Attendance_setCall`(
    IN in_nCharNo INT,
    IN in_dAttendaceDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tAttendance WHERE nCharNo = in_nCharNo AND dAttendanceDate = in_dAttendaceDate ) THEN
        SET in_nRet=-1111;
    ELSE
        INSERT INTO tAttendance(nCharNo,dAttendanceDate) VALUES (in_nCharNo,in_dAttendaceDate);
        IF 0 <> 0 THEN
        SET in_nRet=-9999;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Attendance_setReward
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Attendance_setReward`;
DELIMITER $$
CREATE PROCEDURE `usp_Attendance_setReward`(
    IN in_nCharNo INT,
    IN in_dEventDate DATETIME,
    IN in_nRewardID INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tAttendanceReward WHERE nCharNo = in_nCharNo AND dEventDate = CAST( in_dEventDate as datetime ) AND nRewardID = in_nRewardID ) THEN
        SET in_nRet=-1111;
    ELSE
        INSERT INTO tAttendanceReward(nCharNo,dEventDate,dRewardDate,nRewardID) VALUES (in_nCharNo,in_dEventDate,NOW(),in_nRewardID);
        IF 0 <> 0 THEN
        SET in_nRet=-9999;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Char_chkCharNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Char_chkCharNo`;
DELIMITER $$
CREATE PROCEDURE `usp_Char_chkCharNo`(
    IN in_nCharNo INT,
    IN in_sID VARCHAR(40),
    IN in_nRet INT
)
sp_body: BEGIN

    IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND sID = in_sID ) THEN
        SET in_nRet=0;
    ELSE
        SET in_nRet=-1;
    END IF;
    SELECT nCharNo,sID FROM tCharacter WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_CharacterLoading
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_CharacterLoading`;
DELIMITER $$
CREATE PROCEDURE `usp_CharacterLoading`(
    IN in_nCharNo INT,
    IN in_dWorldTime DATETIME,
    IN in_sLoadingMapName VARCHAR(16),
    IN in_nFailType SMALLINT,
    IN in_nLoadingFailCnt SMALLINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nSeqNo INT;

    SET v_nSeqNo=0;
    START TRANSACTION;
    SELECT nSeqNo INTO v_nSeqNo FROM tCharacterLoading WHERE nCharNo=in_nCharNo ORDER BY nSeqNo DESC LIMIT 1;
    SET v_nSeqNo=v_nSeqNo+1;
    INSERT INTO tCharacterLoading(nCharNo,nSeqNo,dWorldTime,sLoadingMapName,nFailType,nLoadingFailCnt) VALUES (in_nCharNo,v_nSeqNo,in_dWorldTime,in_sLoadingMapName,in_nFailType,in_nLoadingFailCnt);
    IF 0 <> 0 THEN
        SET out_nRet=-1001;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    SET out_nRet=0;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_chkKeyMap
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_chkKeyMap`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_chkKeyMap`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = in_nCharNo ) THEN
        SELECT 1;
    ELSE
        SELECT 0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getAutoScroll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getAutoScroll`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getAutoScroll`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nAutoBuff,nAutoChargedbuff FROM tCharacterSingleOption WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getChangeIDList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getChangeIDList`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getChangeIDList`(
    IN in_nUserNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nRowNo,nUserNo,nCharNo,bInit FROM tCharacterChangeID WHERE nUserNo=in_nUserNo AND bFlag=0;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getGameOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getGameOption`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getGameOption`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = in_nCharNo ) THEN
        SELECT T1.nOptionNo,IFNULL(T2.nValue,T1.nDefaultValue) AS nValue FROM tGameOptionInit AS T1 LEFT OUTER JOIN (SELECT* FROM tCharacterGameOption WHERE nCharNo=in_nCharNo) AS T2 ON T1.nOptionNo=T2.nOptionNo WHERE T1.bIsActive=1;
    ELSE
        SELECT nOptionNo,nDefaultValue AS nValue FROM tGameOptionInit WHERE bIsActive=1;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getItemActionCoolTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getItemActionCoolTime`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getItemActionCoolTime`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = in_nCharNo AND dEndTime < NOW ( ) ) THEN
        START TRANSACTION;
        DELETE FROM tItemActionCoolTime WHERE nCharNo=in_nCharNo AND dEndTime<NOW();
        IF ROW_COUNT ( ) = 0 THEN
        SET in_nRet=-3294;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SELECT nCharNo,nActionGroup,dEndTime FROM tItemActionCoolTime WHERE nCharNo=in_nCharNo;
        COMMIT;
    ELSE
        SELECT nCharNo,nActionGroup,dEndTime FROM tItemActionCoolTime WHERE nCharNo=in_nCharNo;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getKeyMap
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getKeyMap`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getKeyMap`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = in_nCharNo ) THEN
        SELECT nFunctionNo,nExtendedKey,nASCIICode FROM tCharacterKeyMap WHERE nCharNo=in_nCharNo;
    ELSE
        SET in_nRet=-2725;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getShortCutList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getShortCutList`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getShortCutList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted = 0 ) THEN
        SELECT nSlotNo,nCodeNo,nValue FROM tCharacterShortCut WHERE nCharNo=in_nCharNo AND nCodeNo>0;
    ELSE
        SET in_nRet=-2721;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_getSingleOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_getSingleOption`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_getSingleOption`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        INSERT INTO tCharacterSingleOption(nCharNo) VALUES (in_nCharNo);
        IF 0 <> 0 THEN
        SET in_nRet=-3192;
        LEAVE sp_body;
    END IF;
    END IF;
    SELECT bMysteryBolt,bAutoPick,dRestExpLastExecuteTime,nChatBan FROM tCharacterSingleOption WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_initKeyMap
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_initKeyMap`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_initKeyMap`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = in_nCharNo ) THEN
        START TRANSACTION;
        DELETE FROM tCharacterKeyMap WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-2723;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        INSERT INTO tCharacterKeyMap(nCharNo,nFunctionNo,sys_create_date,nExtendedKey,nASCIICode)SELECT in_nCharNo,nFunctionNo,NOW(),nExtendedKey,nASCIICode FROM tKeyMapInit;
        IF 0 <> 0 THEN
        SET in_nRet=-2724;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        COMMIT;
    ELSE
        INSERT INTO tCharacterKeyMap(nCharNo,nFunctionNo,sys_create_date,nExtendedKey,nASCIICode)SELECT in_nCharNo,nFunctionNo,NOW(),nExtendedKey,nASCIICode FROM tKeyMapInit;
        IF 0 <> 0 THEN
        SET in_nRet=-2724;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setAutoPick
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setAutoPick`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setAutoPick`(
    IN in_nCharNo INT,
    IN in_bAutoPick TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        INSERT INTO tCharacterSingleOption(nCharNo,bAutoPick) VALUES (in_nCharNo,in_bAutoPick);
        IF 0 <> 0 THEN
        SET in_nRet=-3192;
        LEAVE sp_body;
    END IF;
    ELSE
        UPDATE tCharacterSingleOption SET bAutoPick=in_bAutoPick WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-3193;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setBuffScroll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setBuffScroll`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setBuffScroll`(
    IN in_nCharNo INT,
    IN in_nAutoBuff TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nAutoChargedBuff TINYINT;

    SET in_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tCharacterSingleOption SET nAutoBuff=in_nAutoBuff WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    ELSE
        SET v_nAutoChargedBuff=0;
        INSERT INTO tCharacterSingleOption(nCharNo,nAutoBuff,nAutoChargedBuff) VALUES (in_nCharNo,in_nAutoBuff,v_nAutoChargedBuff);
        IF 0 <> 0 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setChangeID
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setChangeID`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setChangeID`(
    IN in_nRowNo INT,
    IN in_nCharNo INT,
    IN in_sID VARCHAR(16),
    IN in_bInit TINYINT,
    IN in_nUserNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_sOldID VARCHAR(50);
    DECLARE v_nDBUserNo INT;
    DECLARE v_nGrade TINYINT;
    DECLARE v_nMyGuildNo INT;

    SET in_nRet=0;
    SELECT T1.sID,T1.nUserNo,nGrade INTO v_sOldID,v_nDBUserNo,v_nGrade FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo=T2.nCharNo WHERE T1.nCharNo=in_nCharNo AND T1.bDeleted=0 AND bDeleted=0;
    IF v_sOldID = in_sID THEN
        SET in_nRet=-3039;
        LEAVE sp_body;
    END IF;
    IF v_nDBUserNo IS NULL THEN
        SET in_nRet=-3023;
        LEAVE sp_body;
    ELSE
        IF v_nDBUserNo <> in_nUserNo THEN
        SET in_nRet=-3024;
        LEAVE sp_body;
    END IF;
        IF v_nGrade IS NOT NULL AND v_nGrade = 0 AND in_bInit = 1 THEN
        SET in_nRet=-3025;
        LEAVE sp_body;
    END IF;
        IF NOT EXISTS ( SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = in_nCharNo ) THEN
        SET in_nRet=-3024;
        LEAVE sp_body;
    END IF;
        START TRANSACTION;
        IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE sID = in_sID ) THEN
        SET in_nRet=-3039;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        UPDATE tCharacter SET sID=in_sID WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3026;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        IF in_bInit = 1 THEN
        IF EXISTS ( SELECT nCharNo FROM tFriend WHERE nCharNo = in_nCharNo ) THEN
        DELETE FROM tFriend WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 OR ROW_COUNT ( ) < 1 THEN
        SET in_nRet=-3027;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        IF EXISTS ( SELECT nCharNo FROM tFriend WHERE nFriendCharNo = in_nCharNo ) THEN
        DELETE FROM tFriend WHERE nFriendCharNo=in_nCharNo;
        IF 0 <> 0 OR ROW_COUNT ( ) < 1 THEN
        SET in_nRet=-3027;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nCharNo = in_nCharNo ) THEN
        DELETE FROM tMarriage WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 OR ROW_COUNT ( ) < 1 THEN
        SET in_nRet=-3028;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        IF EXISTS ( SELECT nCharNo FROM tMarriage WHERE nPartnerCharNo = in_nCharNo ) THEN
        DELETE FROM tMarriage WHERE nPartnerCharNo=in_nCharNo;
        IF 0 <> 0 OR ROW_COUNT ( ) < 1 THEN
        SET in_nRet=-3028;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        IF EXISTS ( SELECT nCharNo FROM tHolyPromise WHERE nCharNo = in_nCharNo ) THEN
        DELETE FROM tHolyPromise WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 OR ROW_COUNT ( ) < 1 THEN
        SET in_nRet=-3029;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setChangeIDItemUsed
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setChangeIDItemUsed`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setChangeIDItemUsed`(
    IN in_nCharNo INT,
    IN in_bInit TINYINT,
    IN in_nUserNo INT,
    IN in_nItemKey BIGINT,
    IN in_nItemOption SMALLINT,
    IN in_nItemLot SMALLINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nGrade TINYINT;

    SET in_nRet=0;
    CALL p_Item_IsValidOwner(in_nItemKey,in_nCharNo,in_nRet);
    IF in_nRet <> 0 THEN
        SET in_nRet=-3036;
        LEAVE sp_body;
    END IF;
    SELECT nGrade INTO v_nGrade FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo=T2.nCharNo WHERE T1.nCharNo=in_nCharNo AND T1.bDeleted=0 AND bDeleted=0;
    IF v_nGrade IS NOT NULL AND v_nGrade = 0 AND in_bInit = 1 THEN
        SET in_nRet=-3025;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nRowNo FROM tCharacterChangeID WHERE nCharNo = in_nCharNo AND bFlag = 0 ) THEN
        SET in_nRet=-3022;
        LEAVE sp_body;
    ELSE
        IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nUserNo = in_nUserNo AND bDeleted = 0 ) THEN
        START TRANSACTION;
        INSERT INTO tCharacterChangeID(nUserNo,nCharNo,bInit) VALUES (in_nUserNo,in_nCharNo,in_bInit);
        IF 0 <> 0 OR LAST_INSERT_ID ( ) = 0 THEN
        SET in_nRet=-3092;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        IF in_nItemLot = 0 THEN
        CALL p_Item_Delete(in_nItemKey,in_nRet);
        IF in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3037;
        LEAVE sp_body;
    END IF;
    ELSE
        CALL p_Item_SetOption(in_nItemKey,in_nItemOption,in_nItemLot,in_nRet);
        IF in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3038;
        LEAVE sp_body;
    END IF;
    END IF;
        COMMIT;
    ELSE
        SET in_nRet=-3023;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setChargedbuffScroll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setChargedbuffScroll`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setChargedbuffScroll`(
    IN in_nCharNo INT,
    IN in_nAutoChargedBuff TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nAutoBuff TINYINT;

    SET in_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tCharacterSingleOption SET nAutoChargedBuff=in_nAutoChargedBuff WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    ELSE
        SET v_nAutoBuff=0;
        INSERT INTO tCharacterSingleOption(nCharNo,nAutoBuff,nAutoChargedBuff) VALUES (in_nCharNo,v_nAutoBuff,in_nAutoChargedBuff);
        IF 0 <> 0 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setChatBan
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setChatBan`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setChatBan`(
    IN in_sID VARCHAR(40),
    IN in_nChatBan INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nCharNo INT;

    SET in_nRet=0;
    SELECT nCharNo INTO v_nCharNo FROM tCharacter WHERE sID=in_sID AND bDeleted=0;
    IF v_nCharNo IS NULL THEN
        SET in_nRet=-3121;
    ELSE
        IF NOT EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = v_nCharNo ) THEN
        INSERT INTO tCharacterSingleOption(nCharNo,nChatBan) VALUES (v_nCharNo,in_nChatBan);
        IF 0 <> 0 THEN
        SET in_nRet=-3192;
        LEAVE sp_body;
    END IF;
    ELSE
        UPDATE tCharacterSingleOption SET nChatBan=in_nChatBan WHERE nCharNo=v_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-3193;
        LEAVE sp_body;
    END IF;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setGameOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setGameOption`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setGameOption`(
    IN in_nCharNo INT,
    IN in_nOptionNo SMALLINT,
    IN in_nValue TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tCharacterGameOption WHERE nCharNo = in_nCharNo AND nOptionNo = in_nOptionNo ) THEN
        UPDATE tCharacterGameOption SET nValue=in_nValue WHERE nCharNo=in_nCharNo AND nOptionNo=in_nOptionNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2792;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tCharacterGameOption(nCharNo,nOptionNo,nValue) VALUES (in_nCharNo,in_nOptionNo,in_nValue);
        IF 0 <> 0 THEN
        SET in_nRet=-2792;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setGameOptionInit
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setGameOptionInit`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setGameOptionInit`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    DELETE FROM tCharacterGameOption WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setItemActionCoolTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setItemActionCoolTime`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setItemActionCoolTime`(
    IN in_nCharNo INT,
    IN in_nActionGroup SMALLINT,
    IN in_dEndTime DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = in_nCharNo AND nActionGroup = in_nActionGroup ) THEN
        START TRANSACTION;
        UPDATE tItemActionCoolTime SET dEndTime=in_dEndTime WHERE nCharNo=in_nCharNo AND nActionGroup=in_nActionGroup;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3291;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        COMMIT;
    ELSE
        INSERT INTO tItemActionCoolTime(nCharNo,nActionGroup,dEndTime) VALUES (in_nCharNo,in_nActionGroup,in_dEndTime);
        IF 0 <> 0 THEN
        SET in_nRet=-3291;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setKeyMap
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setKeyMap`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setKeyMap`(
    IN in_nCharNo INT,
    IN in_nFunctionNo SMALLINT,
    IN in_nExtendedKey TINYINT,
    IN in_nASCIICode TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = in_nCharNo AND nFunctionNo = in_nFunctionNo ) THEN
        UPDATE tCharacterKeyMap SET nExtendedKey=in_nExtendedKey,nASCIICode=in_nASCIICode WHERE nCharNo=in_nCharNo AND nFunctionNo=in_nFunctionNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2793;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tCharacterKeyMap(nCharNo,nFunctionNo,nExtendedKey,nASCIICode) VALUES (in_nCharNo,in_nFunctionNo,in_nExtendedKey,in_nASCIICode);
        IF 0 <> 0 THEN
        SET in_nRet=-2792;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setMysteryBolt
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setMysteryBolt`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setMysteryBolt`(
    IN in_nCharNo INT,
    IN in_bMysteryBolt TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        INSERT INTO tCharacterSingleOption(nCharNo,bMysteryBolt) VALUES (in_nCharNo,in_bMysteryBolt);
        IF 0 <> 0 THEN
        SET in_nRet=-3192;
        LEAVE sp_body;
    END IF;
    ELSE
        UPDATE tCharacterSingleOption SET bMysteryBolt=in_bMysteryBolt WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-3193;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setRestExpExecTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setRestExpExecTime`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setRestExpExecTime`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = in_nCharNo ) THEN
        INSERT INTO tCharacterSingleOption(nCharNo) VALUES (in_nCharNo);
        IF 0 <> 0 THEN
        SET in_nRet=-3192;
        LEAVE sp_body;
    END IF;
    ELSE
        UPDATE tCharacterSingleOption SET dRestExpLastExecuteTime=NOW() WHERE nCharNo=in_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-3193;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setShortCut
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setShortCut`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setShortCut`(
    IN in_nCharNo INT,
    IN in_nSlotNo TINYINT,
    IN in_nCodeNo SMALLINT,
    IN in_nValue INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nCodeType SMALLINT;
    DECLARE v_nDBCharNo INT;
    DECLARE v_nPrevCodeNo SMALLINT;
    DECLARE v_nPrevValue INT;

    SET in_nRet=0;
    SET v_nCodeType=0;
    SET v_nPrevCodeNo=0;
    SET v_nPrevValue=0;
    SELECT nCharNo,nCodeNo,nValue INTO v_nDBCharNo,v_nPrevCodeNo,v_nPrevValue FROM tCharacterShortCut WHERE nCharNo=in_nCharNo AND nSlotNo=in_nSlotNo;
    IF v_nDBCharNo IS NULL THEN
        INSERT INTO tCharacterShortCut(nCharNo,nSlotNo,nCodeType,nCodeNo,nValue) VALUES (in_nCharNo,in_nSlotNo,1,in_nCodeNo,in_nValue);
        IF 0 <> 0 THEN
        SET in_nRet=-2791;
    END IF;
    ELSE
        IF v_nPrevCodeNo = in_nCodeNo AND v_nPrevValue = in_nValue THEN
        SET in_nRet=-2722;
    ELSE
        UPDATE tCharacterShortCut SET nCodeNo=in_nCodeNo,nValue=in_nValue,sys_update_date=NOW() WHERE nCharNo=in_nCharNo AND nSlotNo=in_nSlotNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2793;
    END IF;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Character_setStatPoint
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Character_setStatPoint`;
DELIMITER $$
CREATE PROCEDURE `usp_Character_setStatPoint`(
    IN in_nCharNo INT,
    IN in_nMaxStatPoint INT,
    IN in_nStrength TINYINT,
    IN in_nConstitute TINYINT,
    IN in_nDexterity TINYINT,
    IN in_nIntelligence TINYINT,
    IN in_nMentalPower TINYINT,
    IN in_nItemID INT,
    IN in_nOptionType SMALLINT,
    IN in_nDeleteCount INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nTotalStat INT;
    DECLARE v_nDeleteCheck INT;
    DECLARE v_nStr INT;
    DECLARE v_nCon INT;
    DECLARE v_nDex INT;
    DECLARE v_nInt INT;
    DECLARE v_nMen INT;
    DECLARE v_nTotalItemCount INT;
    DECLARE v_nItemKey BIGINT;
    DECLARE v_nOptionData BIGINT;
    DECLARE v_nItemDeleteRemain INT;
    DECLARE v_nRemainCount INT;

    SET in_nRet=0;
    SET v_nTotalStat=0;
    SET v_nDeleteCheck=0;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1011;
        LEAVE sp_body;
    END IF;
    IF in_nStrength < 0 OR in_nConstitute < 0 OR in_nDexterity < 0 OR in_nIntelligence < 0 OR in_nMentalPower < 0 THEN
        SET in_nRet=-1012;
        LEAVE sp_body;
    END IF;
    SET v_nTotalStat=in_nStrength+in_nConstitute+in_nDexterity+in_nIntelligence+in_nMentalPower;
    IF in_nMaxStatPoint < v_nTotalStat THEN
        SET in_nRet=-1013;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted = 0 ) THEN
        SELECT nStrength,nConstitute,nDexterity,nIntelligence,nMentalPower INTO v_nStr,v_nCon,v_nDex,v_nInt,v_nMen FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStr > in_nStrength THEN
        SET v_nDeleteCheck=v_nDeleteCheck+(v_nStr-in_nStrength);
    END IF;
        IF v_nCon > in_nConstitute THEN
        SET v_nDeleteCheck=v_nDeleteCheck+(v_nCon-in_nConstitute);
    END IF;
        IF v_nDex > in_nDexterity THEN
        SET v_nDeleteCheck=v_nDeleteCheck+(v_nDex-in_nDexterity);
    END IF;
        IF v_nInt > in_nIntelligence THEN
        SET v_nDeleteCheck=v_nDeleteCheck+(v_nInt-in_nIntelligence);
    END IF;
        IF v_nMen > in_nMentalPower THEN
        SET v_nDeleteCheck=v_nDeleteCheck+(v_nMen-in_nMentalPower);
    END IF;
        IF in_nDeleteCount <> v_nDeleteCheck THEN
        SET in_nRet=-1021;
        LEAVE sp_body;
    END IF;
        START TRANSACTION;
        UPDATE tCharacter SET nStrength=in_nStrength,nConstitute=in_nConstitute,nDexterity=in_nDexterity,nIntelligence=in_nIntelligence,nMentalPower=in_nMentalPower WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1093;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        IF in_nDeleteCount > 0 THEN
        IF in_nDeleteCount = 1 THEN
        SELECT T1.nItemKey,T2.nOptionData INTO v_nItemKey,v_nOptionData FROM tItem T1,tItemOptions T2 WHERE T1.nItemKey=T2.nItemKey AND T1.nOwner=in_nCharNo AND T1.nStorageType=9 AND T1.nItemID=in_nItemID ORDER BY nStorage ASC LIMIT 1;
        IF v_nOptionData > 1 THEN
        SET v_nOptionData=v_nOptionData-1;
        CALL p_Item_SetOption(v_nItemKey,in_nOptionType,v_nOptionData,in_nRet);
        IF in_nRet <> 0 THEN
        SET in_nRet=-1022;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        CALL p_Item_Delete(v_nItemKey,n);
    END IF;
    END IF;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_Draw
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_Draw`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_Draw`(
    IN in_ROW_NO INT,
    IN in_USER_NO INT,
    IN in_RET INT
)
sp_body: BEGIN

    CALL usp_Charge_ItemDraw(in_ROW_NO,in_USER_NO,in_RET);

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_GetList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_GetList`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_GetList`(
    IN in_USER_NO INT
)
sp_body: BEGIN

    CALL usp_Charge_ItemSelectList(in_USER_NO);

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseSkillPowerInitDel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseSkillPowerInitDel`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseSkillPowerInitDel`(
    IN in_nCharNo INT,
    IN in_nDeleteItemKey BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_TempRet INT;

    START TRANSACTION;
    UPDATE tSkill SET nSkillPowerDemage=0,nSkillPowerSP=0,nSkillPowerKeepTime=0,nSkillPowerCoolTime=0 WHERE nCharNo=in_nCharNo;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nDeleteItemKey,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseSkillPowerInitLotChange
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseSkillPowerInitLotChange`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseSkillPowerInitLotChange`(
    IN in_nCharNo INT,
    IN in_nChangeItemKey BIGINT,
    IN in_nChangeItemOptionType SMALLINT,
    IN in_nChangeItemOptionData BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_TempRet INT;

    START TRANSACTION;
    UPDATE tSkill SET nSkillPowerDemage=0,nSkillPowerSP=0,nSkillPowerKeepTime=0,nSkillPowerCoolTime=0 WHERE nCharNo=in_nCharNo;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nChangeItemKey,in_nChangeItemOptionType,in_nChangeItemOptionData,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseStatPointDecDel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseStatPointDecDel`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseStatPointDecDel`(
    IN in_nCharNo INT,
    IN in_nStatPointType TINYINT,
    IN in_nDeleteItemKey BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nStatPoint TINYINT;
    DECLARE v_TempRet INT;

    START TRANSACTION;
    IF in_nStatPointType = 0 THEN
        SELECT nStrength INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nStrength=nStrength-1 WHERE nCharNo=in_nCharNo AND nStrength>0;
    END IF;
    IF in_nStatPointType = 1 THEN
        SELECT nConstitute INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nConstitute=nConstitute-1 WHERE nCharNo=in_nCharNo AND nConstitute>0;
    END IF;
    IF in_nStatPointType = 2 THEN
        SELECT nDexterity INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nDexterity=nDexterity-1 WHERE nCharNo=in_nCharNo AND nDexterity>0;
    END IF;
    IF in_nStatPointType = 3 THEN
        SELECT nIntelligence INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nIntelligence=nIntelligence-1 WHERE nCharNo=in_nCharNo AND nIntelligence>0;
    END IF;
    IF in_nStatPointType = 4 THEN
        SELECT nMentalPower INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMentalPower=nMentalPower-1 WHERE nCharNo=in_nCharNo AND nMentalPower>0;
    END IF;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nDeleteItemKey,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseStatPointDecLotChange
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseStatPointDecLotChange`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseStatPointDecLotChange`(
    IN in_nCharNo INT,
    IN in_nStatPointType TINYINT,
    IN in_nChangeItemKey BIGINT,
    IN in_nChangeItemOptionType SMALLINT,
    IN in_nChangeItemOptionData BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_nStatPoint TINYINT;
    DECLARE v_TempRet INT;

    START TRANSACTION;
    IF in_nStatPointType = 0 THEN
        SELECT nStrength INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nStrength=nStrength-1 WHERE nCharNo=in_nCharNo AND nStrength>0;
    END IF;
    IF in_nStatPointType = 1 THEN
        SELECT nConstitute INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nConstitute=nConstitute-1 WHERE nCharNo=in_nCharNo AND nConstitute>0;
    END IF;
    IF in_nStatPointType = 2 THEN
        SELECT nDexterity INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nDexterity=nDexterity-1 WHERE nCharNo=in_nCharNo AND nDexterity>0;
    END IF;
    IF in_nStatPointType = 3 THEN
        SELECT nIntelligence INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nIntelligence=nIntelligence-1 WHERE nCharNo=in_nCharNo AND nIntelligence>0;
    END IF;
    IF in_nStatPointType = 4 THEN
        SELECT nMentalPower INTO v_nStatPoint FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nStatPoint < 1 THEN
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMentalPower=nMentalPower-1 WHERE nCharNo=in_nCharNo AND nMentalPower>0;
    END IF;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nChangeItemKey,in_nChangeItemOptionType,in_nChangeItemOptionData,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseStatPointInitDel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseStatPointInitDel`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseStatPointInitDel`(
    IN in_nCharNo INT,
    IN in_nDeleteItemKey BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_TempRet INT;

    START TRANSACTION;
    UPDATE tCharacter SET nStrength=0,nConstitute=0,nDexterity=0,nIntelligence=0,nMentalPower=0 WHERE nCharNo=in_nCharNo;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nDeleteItemKey,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChargeItem_UseStatPointInitLotChange
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChargeItem_UseStatPointInitLotChange`;
DELIMITER $$
CREATE PROCEDURE `usp_ChargeItem_UseStatPointInitLotChange`(
    IN in_nCharNo INT,
    IN in_nChangeItemKey BIGINT,
    IN in_nChangeItemOptionType SMALLINT,
    IN in_nChangeItemOptionData BIGINT,
    OUT out_nRet TINYINT
)
sp_body: BEGIN

    DECLARE v_TempRet INT;

    START TRANSACTION;
    UPDATE tCharacter SET nStrength=0,nConstitute=0,nDexterity=0,nIntelligence=0,nMentalPower=0 WHERE nCharNo=in_nCharNo;
    IF 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nChangeItemKey,in_nChangeItemOptionType,in_nChangeItemOptionData,v_TempRet);
    IF v_TempRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=2;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatBlock_delBlock
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatBlock_delBlock`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatBlock_delBlock`(
    IN in_nUserNo INT,
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nUserNo , 0 ) <= 0 THEN
        SET in_nRet=-2311;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2312;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nUserNo FROM tChatBlock WHERE nUserNo = in_nUserNo AND nCharNo = in_nCharNo ) THEN
        SET in_nRet=-2323;
        LEAVE sp_body;
    ELSE
        DELETE FROM tChatBlock WHERE nUserNo=in_nUserNo AND nCharNo=in_nCharNo;
        IF 0 <> 0 THEN
        SET in_nRet=-2394;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatBlock_delBlockALL
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatBlock_delBlockALL`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatBlock_delBlockALL`(
    IN in_nUserNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nUserNo , 0 ) <= 0 THEN
        SET in_nRet=-2311;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nUserNo FROM tChatBlock WHERE nUserNo = in_nUserNo ) THEN
        SET in_nRet=-2321;
        LEAVE sp_body;
    ELSE
        DELETE FROM tChatBlock WHERE nUserNo=in_nUserNo;
        IF 0 <> 0 THEN
        SET in_nRet=-2394;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatBlock_getBlockList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatBlock_getBlockList`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatBlock_getBlockList`(
    IN in_nUserNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nUserNo , 0 ) <= 0 THEN
        SET in_nRet=-2311;
        LEAVE sp_body;
    END IF;
    SELECT T1.nCharNo,T2.sID,T2.nAdminLevel,T2.bDeleted FROM tChatBlock T1,tCharacter T2 WHERE T1.nUserNo=in_nUserNo AND T1.nCharNo=T2.nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatBlock_setBlock
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatBlock_setBlock`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatBlock_setBlock`(
    IN in_nUserNo INT,
    IN in_sID VARCHAR(40),
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_cnt INT;
    DECLARE v_maxCnt INT;
    DECLARE v_nCharNo INT;

    SET in_nRet=0;
    SET v_nCharNo=0;
    SET v_maxCnt=100;
    IF IFNULL ( in_nUserNo , 0 ) <= 0 THEN
        SET in_nRet=-2311;
        LEAVE sp_body;
    END IF;
    IF CHAR_LENGTH ( IFNULL ( in_sID , '' ) ) = 0 THEN
        SET in_nRet=-2312;
        LEAVE sp_body;
    END IF;
    SELECT nCharNo INTO v_nCharNo FROM tCharacter WHERE sID=in_sID AND bDeleted=0;
    IF v_nCharNo <= 0 THEN
        SET in_nRet=-2323;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nUserNo FROM tChatBlock WHERE nUserNo = in_nUserNo AND nCharNo = v_nCharNo ) THEN
        SET in_nRet=-2321;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(COUNT(*),0) INTO v_cnt FROM tChatBlock WHERE nUserNo=in_nUserNo;
    IF v_cnt > v_maxCnt THEN
        SET in_nRet=-2322;
        LEAVE sp_body;
    ELSE
        INSERT INTO tChatBlock(nUserNo,nCharNo) VALUES (in_nUserNo,v_nCharNo);
        IF 0 <> 0 THEN
        SET in_nRet=-2392;
        LEAVE sp_body;
    END IF;
        SELECT T1.nCharNo,T2.sID,T2.nAdminLevel,T2.bDeleted FROM tChatBlock T1,tCharacter T2 WHERE T1.nUserNo=in_nUserNo AND T1.nCharNo=v_nCharNo AND T1.nCharNo=T2.nCharNo;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatColor_get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatColor_get`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatColor_get`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tChatColor WHERE nCharNo = in_nCharNo ) THEN
        SELECT nFontColorID,nBalloonColorID FROM tChatColor WHERE nCharNo=in_nCharNo;
    ELSE
        SELECT 0 AS nFontColorID,0 AS nBalloonColorID;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatColor_setBalloon
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatColor_setBalloon`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatColor_setBalloon`(
    IN in_nCharNo INT,
    IN in_nBalloonColorID TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nFontColorID TINYINT;

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tChatColor WHERE nCharNo = in_nCharNo ) THEN
        START TRANSACTION;
        UPDATE tChatColor SET nBalloonColorID=in_nBalloonColorID WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-1;
    ELSE
        COMMIT;
    END IF;
    ELSE
        SET v_nFontColorID=0;
        INSERT INTO tChatColor(nCharNo,nFontColorID,nBalloonColorID) VALUES (in_nCharNo,v_nFontColorID,in_nBalloonColorID);
        IF 0 <> 0 THEN
        SET in_nRet=-2;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_ChatColor_setFont
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ChatColor_setFont`;
DELIMITER $$
CREATE PROCEDURE `usp_ChatColor_setFont`(
    IN in_nCharNo INT,
    IN in_nFontColorID TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nBalloonColorID TINYINT;

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tChatColor WHERE nCharNo = in_nCharNo ) THEN
        START TRANSACTION;
        UPDATE tChatColor SET nFontColorID=in_nFontColorID WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-1;
    ELSE
        COMMIT;
    END IF;
    ELSE
        SET v_nBalloonColorID=0;
        INSERT INTO tChatColor(nCharNo,nFontColorID,nBalloonColorID) VALUES (in_nCharNo,in_nFontColorID,v_nBalloonColorID);
        IF 0 <> 0 THEN
        SET in_nRet=-2;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_getCardList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_getCardList`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_getCardList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    SELECT nCardID,nSerialNo,nStarNumber,dDate FROM tCollectionCard WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_getFavoritesList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_getFavoritesList`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_getFavoritesList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    SELECT T1.nSlotNo,T1.nCardID,T2.nSerialNo,T2.nStarNumber FROM tCollectionFavorites T1,tCollectionCard T2 WHERE T1.nCharNo=T2.nCharNo AND T1.nCardID=T2.nCardID AND T1.nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_getRewardList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_getRewardList`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_getRewardList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    SELECT nReward FROM tCollectionReward WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_getSerialList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_getSerialList`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_getSerialList`(
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nCardID,nCount FROM tCollectionSerialCount;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_setCard
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_setCard`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_setCard`(
    IN in_nCharNo INT,
    IN in_nCardID INT,
    IN in_nSerialNo BIGINT,
    IN in_nStarNumber TINYINT,
    IN in_nItemKey BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCardID , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nSerialNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nStarNumber , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nItemKey ) THEN
        START TRANSACTION;
        IF EXISTS ( SELECT nCharNo FROM tCollectionCard WHERE nCharNo = in_nCharNo AND nCardID = in_nCardID ) THEN
        UPDATE tCollectionCard SET nSerialNo=in_nSerialNo,nStarNumber=in_nStarNumber WHERE nCharNo=in_nCharNo AND nCardID=in_nCardID;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tCollectionCard(nCharNo,nCardID,nSerialNo,nStarNumber) VALUES (in_nCharNo,in_nCardID,in_nSerialNo,in_nStarNumber);
        IF 0 <> 0 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
        DELETE FROM tItem WHERE nItemKey=in_nItemKey;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        COMMIT;
    ELSE
        SET in_nRet=-9999;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_setFavorites
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_setFavorites`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_setFavorites`(
    IN in_nCharNo INT,
    IN in_nSlotNo TINYINT,
    IN in_nCardID INT,
    IN in_bRegist TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nSlotNo , - 1 ) < 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF in_bRegist = 1 THEN
        IF IFNULL ( in_nCardID , 0 ) = 0 THEN
        SET in_nRet=-1111;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        IF EXISTS ( SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = in_nCharNo AND nSlotNo = in_nSlotNo ) THEN
        UPDATE tCollectionFavorites SET nCardID=in_nCardID WHERE nCharNo=in_nCharNo AND nSlotNo=in_nSlotNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tCollectionFavorites(nCharNo,nSlotNo,nCardID) VALUES (in_nCharNo,in_nSlotNo,in_nCardID);
        IF 0 <> 0 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    ELSE
        IF EXISTS ( SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = in_nCharNo AND nSlotNo = in_nSlotNo ) THEN
        DELETE FROM tCollectionFavorites WHERE nCharNo=in_nCharNo AND nSlotNo=in_nSlotNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-9999;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-5555;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_setReward
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_setReward`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_setReward`(
    IN in_nCharNo INT,
    IN in_nReward INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tCollectionReward WHERE nCharNo = in_nCharNo AND nReward = in_nReward ) THEN
        SET in_nRet=-2222;
        LEAVE sp_body;
    ELSE
        INSERT INTO tCollectionReward(nCharNo,nReward) VALUES (in_nCharNo,in_nReward);
        IF 0 <> 0 THEN
        SET in_nRet=-3333;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Collection_setSerial
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Collection_setSerial`;
DELIMITER $$
CREATE PROCEDURE `usp_Collection_setSerial`(
    IN in_nCardID INT,
    IN in_nCount BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCardID , 0 ) = 0 THEN
        SET in_nRet=-1111;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCardID FROM tCollectionSerialCount WHERE nCardID = in_nCardID ) THEN
        START TRANSACTION;
        UPDATE tCollectionSerialCount SET sys_update_date=NOW(),nCount=in_nCount WHERE nCardID=in_nCardID;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-4444;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        COMMIT;
    ELSE
        IF in_nCount <> 1 THEN
        SET in_nRet=-5555;
        LEAVE sp_body;
    END IF;
        INSERT INTO tCollectionSerialCount(nCardID,nCount) VALUES (in_nCardID,in_nCount);
        IF 0 <> 0 THEN
        SET in_nRet=-6666;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_chgMachineToMoneycard
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_chgMachineToMoneycard`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_chgMachineToMoneycard`(
    IN in_nItemKey BIGINT,
    IN in_nItemIDOld INT,
    IN in_nItemIDNew INT,
    IN in_nOptionType SMALLINT,
    IN in_nTotalDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nIsPlace SMALLINT;
    DECLARE v_nIsPlaceValue BIGINT;

    SET in_nRet=0;
    SET v_nIsPlace=904;
    SET v_nIsPlaceValue=0;
    IF IFNULL ( in_nItemKey , 0 ) <= 0 THEN
        SET in_nRet=-2412;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemIDOld , 0 ) <= 0 THEN
        SET in_nRet=-2417;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemIDNew , 0 ) <= 0 THEN
        SET in_nRet=-2418;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nOptionType , 0 ) <= 0 THEN
        SET in_nRet=-2419;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nTotalDividend , 0 ) <= 0 THEN
        SET in_nRet=-2416;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nItemKey AND nItemID = in_nItemIDOld ) THEN
        START TRANSACTION;
        UPDATE tItem SET nItemID=in_nItemIDNew WHERE nItemKey=in_nItemKey;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2493;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tItemOptions SET nOptionData=v_nIsPlaceValue WHERE nItemKey=in_nItemKey AND nOptionType=v_nIsPlace;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2493;
        ROLLBACK;
    END IF;
        UPDATE tItemOptions SET nOptionData=in_nTotalDividend WHERE nItemKey=in_nItemKey AND nOptionType=in_nOptionType;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2493;
        ROLLBACK;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;
    ELSE
        SET in_nRet=-2422;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_getDiceExpect
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_getDiceExpect`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_getDiceExpect`(
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2411;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemKey , 0 ) <= 0 THEN
        SET in_nRet=-2412;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(COUNT(DISTINCT(A.nGameNo)),0) AS nCount,IFNULL(SUM(nTotalDividend),0) AS nTotalDividend FROM tDice AS A,tDiceResult AS B WHERE A.nGameNo=B.nGameNo AND A.nCharNo=in_nCharNo AND A.nItemKey=in_nItemKey;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_getDiceRank
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_getDiceRank`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_getDiceRank`(
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2411;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemKey , 0 ) <= 0 THEN
        SET in_nRet=-2412;
        LEAVE sp_body;
    END IF;
    SELECT C.sID,COUNT(B.nGameNo) AS nWinCount FROM tDice AS A,tDiceLog AS B,tCharacter AS C WHERE A.nGameNo=B.nGameNo AND B.nCustomCharNo=C.nCharNo AND A.nCharNo=in_nCharNo AND A.nItemKey=in_nItemKey AND B.bStatus=1 AND C.bDeleted=0 GROUP BY C.sID ORDER BY nWinCount DESC LIMIT 20;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_setDice
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_setDice`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_setDice`(
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_sDice VARCHAR(3),
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nGameNo INT;

    SET in_nRet=0;
    SET v_nGameNo=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2411;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemKey , 0 ) <= 0 THEN
        SET in_nRet=-2412;
        LEAVE sp_body;
    END IF;
    IF CHAR_LENGTH ( IFNULL ( in_sDice , '' ) ) <= 0 THEN
        SET in_nRet=-2413;
        LEAVE sp_body;
    END IF;
    INSERT INTO tDice(nCharNo,nItemKey,sDice) VALUES (in_nCharNo,in_nItemKey,in_sDice);
    IF 0 <> 0 OR LAST_INSERT_ID ( ) = 0 THEN
        SET in_nRet=-2492;
        LEAVE sp_body;
    ELSE
        SET v_nGameNo=LAST_INSERT_ID();
        SELECT v_nGameNo;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_setDiceLog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_setDiceLog`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_setDiceLog`(
    IN in_nGameNo INT,
    IN in_nCustomCharNo INT,
    IN in_nSlotNo INT,
    IN in_nBattingMoney BIGINT,
    IN in_nBattingRate INT,
    IN in_bStatus TINYINT,
    IN in_nTaxRate INT,
    IN in_nDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    INSERT INTO tDiceLog(nGameNo,nCustomCharNo,nSlotNo,nBattingMoney,nBattingRate,bStatus,nTaxRate,nDividend) VALUES (in_nGameNo,in_nCustomCharNo,in_nSlotNo,in_nBattingMoney,in_nBattingRate,in_bStatus,in_nTaxRate,in_nDividend);
    IF 0 <> 0 THEN
        SET in_nRet=-2492;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_setDicePay
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_setDicePay`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_setDicePay`(
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_nTotalDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2411;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nItemKey , 0 ) <= 0 THEN
        SET in_nRet=-2412;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nTotalDividend , 0 ) <= 0 THEN
        SET in_nRet=-2416;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nItemKey AND nOwner = in_nCharNo ) THEN
        START TRANSACTION;
        DELETE FROM tItem WHERE nItemKey=in_nItemKey AND nOwner=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2494;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMoney=nMoney+in_nTotalDividend WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2493;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        INSERT INTO tDicePayLog(nCharNo,nItemKey,nTotalDividend) VALUES (in_nCharNo,in_nItemKey,in_nTotalDividend);
        IF 0 <> 0 THEN
        SET in_nRet=-2492;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-2423;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Dice_setDiceResult
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Dice_setDiceResult`;
DELIMITER $$
CREATE PROCEDURE `usp_Dice_setDiceResult`(
    IN in_nGameNo INT,
    IN in_nCustomCharNo INT,
    IN in_nTotalDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_money BIGINT;
    DECLARE v_resultMoney BIGINT;

    SET in_nRet=0;
    SET v_money=0;
    IF IFNULL ( in_nGameNo , 0 ) <= 0 THEN
        SET in_nRet=-2414;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCustomCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2415;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SELECT IFNULL(nMoney,0) INTO v_money FROM tCharacter WHERE nCharNo=in_nCustomCharNo;
    IF ( v_money + in_nTotalDividend ) > 0 THEN
        SET v_resultMoney=v_money+in_nTotalDividend;
    ELSE
        SET v_resultMoney=0;
    END IF;
    UPDATE tCharacter SET nMoney=v_resultMoney WHERE nCharNo=in_nCustomCharNo;
    IF 0 <> 0 OR ROW_COUNT ( ) > 1 THEN
        SET in_nRet=-2493;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    INSERT INTO tDiceResult(nGameNo,nCustomCharNo,nBeforeMoney,nTotalDividend,nAfterMoney) VALUES (in_nGameNo,in_nCustomCharNo,v_money,in_nTotalDividend,v_resultMoney);
    IF 0 <> 0 THEN
        SET in_nRet=-2492;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;
    SELECT nCharNo AS nCustomCharNo,nMoney FROM tCharacter WHERE nCharNo=in_nCustomCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_addCoin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_addCoin`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_addCoin`(
    IN in_nCharNo INT,
    IN in_nItemKey BIGINT,
    IN in_nCoin BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_status TINYINT;
    DECLARE v_coin BIGINT;

    SET v_status=2;
    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-3311;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCoin , 0 ) <= 0 THEN
        SET in_nRet=-3313;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nItemKey AND nOwner = in_nCharNo ) THEN
        START TRANSACTION;
        DELETE FROM tItem WHERE nItemKey=in_nItemKey AND nOwner=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3332;
        LEAVE sp_body;
    END IF;
        SELECT nCoin INTO v_coin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
        IF v_coin IS NULL THEN
        INSERT INTO tGambleCoin(nCharNo,nCoin) VALUES (in_nCharNo,in_nCoin);
        IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3329;
        LEAVE sp_body;
    END IF;
        SET v_coin=0;
    ELSE
        UPDATE tGambleCoin SET sys_update_date=NOW(),nCoin=nCoin+in_nCoin WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3323;
        LEAVE sp_body;
    END IF;
    END IF;
        INSERT INTO tGambleCoinExchangeLog(nCharNo,nExchangedCoin,nStatus,nBeforeCoin,nAfterCoin) VALUES (in_nCharNo,in_nCoin,v_status,v_coin,v_coin+in_nCoin);
        IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3326;
        LEAVE sp_body;
    ELSE
        COMMIT;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-3331;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_getCoin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_getCoin`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_getCoin`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tGambleCoin WHERE nCharNo = in_nCharNo ) THEN
        SELECT nCoin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
    ELSE
        SELECT 0 AS nCoin;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_getDiceRank
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_getDiceRank`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_getDiceRank`(
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT C.sID,COUNT(B.nGameNo) AS nWinCount FROM tGambleDice AS A,tGambleDiceLog AS B,tCharacter AS C WHERE A.nGameNo=B.nGameNo AND B.nCharNo=C.nCharNo AND B.bStatus=1 AND C.bDeleted=0 GROUP BY C.sID ORDER BY nWinCount DESC LIMIT 20;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_getExchangedCoin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_getExchangedCoin`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_getExchangedCoin`(
    IN in_nCharNo INT,
    IN in_dDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT IFNULL(SUM(nExchangedCoin),0) AS nExchangedCoin FROM tGambleCoinExchangeLog WHERE nCharNo=in_nCharNo AND dDate>=in_dDate AND nStatus=1;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_getJackpotCoin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_getJackpotCoin`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_getJackpotCoin`(
    IN in_sGameName VARCHAR(20),
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nBattingUnit,nJackpotCoin FROM tGambleJackpot WHERE sGameName=in_sGameName;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_getJackpotRanking
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_getJackpotRanking`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_getJackpotRanking`(
    IN in_sGameName VARCHAR(20),
    IN in_nBattingUnit BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT sID,nJackpotCoin FROM tGambleJackpotLog WHERE sGameName=in_sGameName AND nBattingUnit=in_nBattingUnit LIMIT 50;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setCoinToMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setCoinToMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setCoinToMoney`(
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    IN in_nCoin BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_status TINYINT;
    DECLARE v_money BIGINT;
    DECLARE v_coin BIGINT;

    SET v_status=0;
    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-3311;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nMoney , 0 ) <= 0 THEN
        SET in_nRet=-3312;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCoin , 0 ) <= 0 THEN
        SET in_nRet=-3313;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SELECT nCoin INTO v_coin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
    IF v_coin IS NULL THEN
        SET in_nRet=-3321;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF v_coin < in_nCoin THEN
        SET in_nRet=-3322;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tGambleCoin SET sys_update_date=NOW(),nCoin=nCoin-in_nCoin WHERE nCharNo=in_nCharNo AND nCoin>=v_coin;
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3323;
        LEAVE sp_body;
    END IF;
    SELECT nMoney INTO v_money FROM tCharacter WHERE nCharNo=in_nCharNo AND bDeleted=0;
    IF v_money IS NULL THEN
        ROLLBACK;
        SET in_nRet=-3324;
        LEAVE sp_body;
    ELSE
        UPDATE tCharacter SET nMoney=nMoney+in_nMoney WHERE nCharNo=in_nCharNo AND bDeleted=0;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3325;
        LEAVE sp_body;
    END IF;
    END IF;
    INSERT INTO tGambleCoinExchangeLog(nCharNo,nExchangedCoin,nStatus,nBeforeCoin,nAfterCoin,nBeforeMoney,nAfterMoney) VALUES (in_nCharNo,in_nCoin,v_status,v_coin,v_coin-in_nCoin,v_money,v_money+in_nMoney);
    IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3326;
        LEAVE sp_body;
    ELSE
        COMMIT;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setDice
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setDice`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setDice`(
    IN in_sDice VARCHAR(3),
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nGameNo INT;

    SET in_nRet=0;
    SET v_nGameNo=0;
    IF CHAR_LENGTH ( IFNULL ( in_sDice , '' ) ) <= 0 THEN
        SET in_nRet=-3314;
        LEAVE sp_body;
    END IF;
    INSERT INTO tGambleDice(sDice) VALUES (in_sDice);
    IF 0 <> 0 OR LAST_INSERT_ID ( ) = 0 THEN
        SET in_nRet=-3392;
        LEAVE sp_body;
    ELSE
        SET v_nGameNo=LAST_INSERT_ID();
        SELECT v_nGameNo;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setDiceLog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setDiceLog`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setDiceLog`(
    IN in_nGameNo INT,
    IN in_nCharNo INT,
    IN in_nSlotNo INT,
    IN in_nBattingMoney BIGINT,
    IN in_nBattingRate INT,
    IN in_bStatus TINYINT,
    IN in_nTaxRate INT,
    IN in_nDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    INSERT INTO tGambleDiceLog(nGameNo,nCharNo,nSlotNo,nBattingMoney,nBattingRate,bStatus,nTaxRate,nDividend) VALUES (in_nGameNo,in_nCharNo,in_nSlotNo,in_nBattingMoney,in_nBattingRate,in_bStatus,in_nTaxRate,in_nDividend);
    IF 0 <> 0 THEN
        SET in_nRet=-3392;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setDiceResult
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setDiceResult`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setDiceResult`(
    IN in_nGameNo INT,
    IN in_nCharNo INT,
    IN in_nTotalDividend BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_coin BIGINT;
    DECLARE v_resultCoin BIGINT;

    SET in_nRet=0;
    SET v_coin=0;
    SET v_resultCoin=0;
    IF IFNULL ( in_nGameNo , 0 ) <= 0 THEN
        SET in_nRet=-3315;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-3311;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SELECT IFNULL(nCoin,0) INTO v_coin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
    IF ( v_coin + in_nTotalDividend ) > 0 THEN
        SET v_resultCoin=v_coin+in_nTotalDividend;
    ELSE
        SET v_resultCoin=0;
    END IF;
    UPDATE tGambleCoin SET nCoin=v_resultCoin,sys_update_date=NOW() WHERE nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3323;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    INSERT INTO tGambleDiceResult(nGameNo,nCharNo,nBeforeMoney,nTotalDividend,nAfterMoney) VALUES (in_nGameNo,in_nCharNo,v_coin,in_nTotalDividend,v_resultCoin);
    IF 0 <> 0 THEN
        SET in_nRet=-3330;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;
    SELECT nCharNo,nCoin FROM tGambleCoin WHERE nCharNo=in_nCharNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setMoneyToCoin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setMoneyToCoin`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setMoneyToCoin`(
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    IN in_nCoin BIGINT,
    IN in_nLimit BIGINT,
    IN in_dDate DATETIME,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_status TINYINT;
    DECLARE v_money BIGINT;
    DECLARE v_coin BIGINT;
    DECLARE v_exchangedCoin BIGINT;

    SET v_status=1;
    SET in_nRet=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-3311;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nMoney , 0 ) <= 0 THEN
        SET in_nRet=-3312;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCoin , 0 ) <= 0 THEN
        SET in_nRet=-3313;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(SUM(nExchangedCoin),0) INTO v_exchangedCoin FROM tGambleCoinExchangeLog WHERE nCharNo=in_nCharNo AND dDate>=in_dDate AND nStatus=v_status;
    IF in_nLimit < v_exchangedCoin + in_nCoin THEN
        SET in_nRet=-3327;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SELECT IFNULL(nMoney,0) INTO v_money FROM tCharacter WHERE nCharNo=in_nCharNo AND bDeleted=0;
    IF v_money IS NULL THEN
        SET in_nRet=-3324;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF v_money < in_nMoney THEN
        SET in_nRet=-3328;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tCharacter SET nMoney=nMoney-in_nMoney WHERE nCharNo=in_nCharNo AND bDeleted=0 AND nMoney>=in_nMoney;
    IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3325;
        LEAVE sp_body;
    END IF;
    SELECT nCoin INTO v_coin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
    IF v_coin IS NULL THEN
        INSERT INTO tGambleCoin(nCharNo,nCoin) VALUES (in_nCharNo,in_nCoin);
        IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3329;
        LEAVE sp_body;
    END IF;
        SET v_coin=0;
    ELSE
        UPDATE tGambleCoin SET sys_update_date=NOW(),nCoin=nCoin+in_nCoin WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        ROLLBACK;
        SET in_nRet=-3323;
        LEAVE sp_body;
    END IF;
    END IF;
    INSERT INTO tGambleCoinExchangeLog(nCharNo,nExchangedCoin,nStatus,nBeforeCoin,nAfterCoin,nBeforeMoney,nAfterMoney) VALUES (in_nCharNo,in_nCoin,v_status,v_coin,v_coin+in_nCoin,v_money,v_money-in_nMoney);
    IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-3326;
        LEAVE sp_body;
    ELSE
        COMMIT;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Gamble_setSlotMachineResult
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Gamble_setSlotMachineResult`;
DELIMITER $$
CREATE PROCEDURE `usp_Gamble_setSlotMachineResult`(
    IN in_nCharNo INT,
    IN in_nBattingUnit BIGINT,
    IN in_nBatting TINYINT,
    IN in_sGameResult VARCHAR(10),
    IN in_bIsWin TINYINT,
    IN in_nPrizeCoin BIGINT,
    IN in_bIsJackpot TINYINT,
    IN in_nJackpotRate INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_gameName VARCHAR(20);
    DECLARE v_charID VARCHAR(80);
    DECLARE v_gameNo BIGINT;
    DECLARE v_battingCoin BIGINT;
    DECLARE v_coin BIGINT;
    DECLARE v_resultCoin BIGINT;
    DECLARE v_jackpotCoin BIGINT;

    SET in_nRet=0;
    SET v_gameName='SlotMachine';
    SET v_battingCoin=0;
    SET v_resultCoin=0;
    SET v_battingCoin=in_nBattingUnit*in_nBatting;
    IF IFNULL ( in_nCharNo , 0 ) = 0 THEN
        SET in_nRet=-3311;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nBattingUnit , 0 ) = 0 THEN
        SET in_nRet=-3314;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nBatting , 0 ) = 0 THEN
        SET in_nRet=-3315;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_sGameResult , '' ) = '' THEN
        SET in_nRet=-3316;
        LEAVE sp_body;
    END IF;
    SELECT sID INTO v_charID FROM tCharacter WHERE nCharNo=in_nCharNo AND bDeleted=0;
    IF v_charID IS NOT NULL THEN
        START TRANSACTION;
        INSERT INTO tGambleSlotMachineLog(nCharNo,nBattingUnit,nBatting,sGameResult,bIsWin,nPrizeCoin) VALUES (in_nCharNo,in_nBattingUnit,in_nBatting,in_sGameResult,in_bIsWin,in_nPrizeCoin);
        IF 0 <> 0 OR LAST_INSERT_ID ( ) = 0 THEN
        SET in_nRet=-3333;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SET v_gameNo=LAST_INSERT_ID();
        SELECT IFNULL(nCoin,0) INTO v_coin FROM tGambleCoin WHERE nCharNo=in_nCharNo;
        IF ( v_coin - v_battingCoin + in_nPrizeCoin ) > 0 THEN
        SET v_resultCoin=v_coin-v_battingCoin+in_nPrizeCoin;
    ELSE
        SET v_resultCoin=0;
    END IF;
        UPDATE tGambleCoin SET nCoin=v_resultCoin,sys_update_date=NOW() WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3334;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        IF in_bIsJackpot = 1 THEN
        SELECT IFNULL(nJackpotCoin,0) INTO v_jackpotCoin FROM tGambleJackpot WHERE sGameName=v_gameName AND nBattingUnit=in_nBattingUnit;
        IF v_jackpotCoin <> 0 THEN
        SET v_jackpotCoin=v_jackpotCoin*in_nJackpotRate/1000;
    END IF;
        UPDATE tGambleCoin SET nCoin=nCoin+v_jackpotCoin,sys_update_date=NOW() WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3335;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tGambleJackpot SET nJackpotCoin=0,sys_update_date=NOW() WHERE sGameName=v_gameName AND nBattingUnit=in_nBattingUnit;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3336;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        INSERT INTO tGambleJackpotLog(sGameName,nBattingUnit,sID,nJackpotCoin,nGameNo) VALUES (v_gameName,in_nBattingUnit,v_charID,v_jackpotCoin,v_gameNo);
        IF 0 <> 0 THEN
        SET in_nRet=-3337;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        IF EXISTS ( SELECT sGameName FROM tGambleJackpot WHERE sGameName = v_gameName AND nBattingUnit = in_nBattingUnit ) THEN
        UPDATE tGambleJackpot SET nJackpotCoin=nJackpotCoin+v_battingCoin,sys_update_date=NOW() WHERE sGameName=v_gameName AND nBattingUnit=in_nBattingUnit;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-3338;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tGambleJackpot(sGameName,nBattingUnit,nJackpotCoin) VALUES (v_gameName,in_nBattingUnit,v_battingCoin);
        IF 0 <> 0 THEN
        SET in_nRet=-3338;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    END IF;
        COMMIT;
    ELSE
        SET in_nRet=-3339;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_delMember
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_delMember`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_delMember`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    DELETE FROM tGuildAcademyMember WHERE nNo=in_nNo AND nCharNo=in_nCharNo;
    IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1592;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_delRewardAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_delRewardAll`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_delRewardAll`(
    IN in_nNo INT,
    IN in_nTypeID TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nTypeID = 1 THEN
        DELETE FROM tGuildAcademyReward WHERE nNo=in_nNo AND nTypeID=in_nTypeID;
        IF 0 <> 0 THEN
        SET in_nRet=-1596;
        LEAVE sp_body;
    END IF;
    ELSE
        DELETE FROM tGuildAcademyReward WHERE nNo=in_nNo AND nTypeID>1;
        IF 0 <> 0 THEN
        SET in_nRet=-1596;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_depositMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_depositMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_depositMoney`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_money BIGINT;

    SET in_nRet=0;
    SET v_money=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    IF in_nMoney <= 0 THEN
        SET in_nRet=-1514;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted = 0 ) THEN
        SET in_nRet=-1524;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nNo FROM tGuildAcademy WHERE nNo = in_nNo ) THEN
        SET in_nRet=-1525;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    SELECT nMoney INTO v_money FROM tCharacter WHERE nCharNo=in_nCharNo;
    IF v_money < in_nMoney THEN
        SET in_nRet=-1526;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        UPDATE tCharacter SET nMoney=nMoney-in_nMoney WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tGuildAcademy SET nMoney=nMoney+in_nMoney WHERE nNo=in_nNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getGraduateCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getGraduateCount`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getGraduateCount`(
    IN in_nNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(COUNT(nCharNo),0) AS cnt FROM tGuildAcademyGraduate WHERE nNo=in_nNo;
    IF 0 <> 0 THEN
        SET in_nRet=-1591;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getGraduateList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getGraduateList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getGraduateList`(
    IN in_nNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    SELECT sID,nClass,nJoinLevel,dJoinDate,nGraduateLevel,dGraduateDate FROM tGuildAcademyGraduate WHERE nNo=in_nNo;
    IF 0 <> 0 THEN
        SET in_nRet=-1591;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getGuildAcademyList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getGuildAcademyList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getGuildAcademyList`(
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT GA.nNo,nMasterCharNo,C1.sID,nMaxMember,GA.nMoney,nPoint,T2.nRanking,IFNULL(dBuffEndTime,NOW()),IFNULL(dNotifyDate,NOW()),IFNULL(sNotify,''),IFNULL(sIntro,''),T1.nCharNo,T1.sID,IFNULL(T3.nGraduateCount,0) AS nGraduateCount FROM tGuildAcademy AS GA INNER JOIN tCharacter AS C1 ON GA.nMasterCharNo=C1.nCharNo INNER JOIN (SELECT G.nNo,GM.nCharNo,C.sID FROM tGuild AS G,tGuildMember AS GM,tCharacter AS C WHERE G.nNo=GM.nNo AND GM.nCharNo=C.nCharNo AND GM.nGrade=0 AND C.bDeleted=0) AS T1 ON GA.nNo=T1.nNo INNER JOIN (SELECT T1.nNo,COUNT(T2.nPoint)+1 AS nRanking FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint<T2.nPoint GROUP BY T1.nNo,T1.nPoint) AS T2 ON GA.nNo=T2.nNo LEFT OUTER JOIN (SELECT nNo,IFNULL(COUNT(nCharNo),0) AS nGraduateCount FROM tGuildAcademyGraduate GROUP BY nNo) AS T3 ON GA.nNo=T3.nNo ORDER BY T2.nRanking;
    IF 0 <> 0 THEN
        SET in_nRet=-1594;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getGuildNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getGuildNo`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getGuildNo`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nNo FROM tGuildAcademyMember WHERE nCharNo = in_nCharNo ) THEN
        SELECT nNo FROM tGuildAcademyMember WHERE nCharNo=in_nCharNo;
    ELSE
        SET in_nRet=0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getHistoryList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getHistoryList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getHistoryList`(
    IN in_nAcademyID INT,
    IN in_nAmount INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT TOP(in_nAmount)sID,nAction,dDate FROM tGuildAcademyHistory WHERE nAcademyID=in_nAcademyID;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getMemberList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getMemberList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getMemberList`(
    IN in_nNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    SELECT G.nCharNo,G.nJoinLevel,G.dJoinDate,G.nGrade,G.sIntro,G.nChatBanTime,C.sID,C.nLevel,C.sLoginZone,C.dLastLoginDate,S.nClass FROM tGuildAcademyMember AS G INNER JOIN tCharacter AS C ON C.nCharNo=G.nCharNo INNER JOIN tCharacterShape AS S ON S.nCharNo=G.nCharNo WHERE G.nNo=in_nNo AND C.bDeleted=0;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getMoney`(
    IN in_nNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nNo FROM tGuildAcademy WHERE nNo = in_nNo ) THEN
        SELECT nMoney FROM tGuildAcademy WHERE nNo=in_nNo;
    ELSE
        SET in_nRet=-1521;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getRanking
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getRanking`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getRanking`(
    IN in_nCount INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_MAX_COUNT INT;

    SET in_nRet=0;
    SELECT IFNULL(MAX(nCount),0) INTO v_MAX_COUNT FROM tGuildAcademyRanking;
    SELECT dDate,nRanking,nGuildNo,sGuildName,nMasterNo,sMasterID,nPoint FROM tGuildAcademyRanking WHERE nCount=v_MAX_COUNT-(in_nCount-1);

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getReward
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getReward`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getReward`(
    IN in_nNo INT,
    IN in_nLevel TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nLevel <= 0 THEN
        SET in_nRet=-1513;
        LEAVE sp_body;
    END IF;
    SELECT nTypeID,nValue,nItemID FROM tGuildAcademyReward WHERE nNo=in_nNo AND nLevel=in_nLevel;
    IF 0 <> 0 THEN
        SET in_nRet=-1591;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_getRewardAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_getRewardAll`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_getRewardAll`(
    IN in_nNo INT,
    IN in_nTypeID INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nTypeID = 1 THEN
        SELECT nLevel,nTypeID,nValue,nItemID FROM tGuildAcademyReward WHERE nNo=in_nNo AND nTypeID=1;
        IF 0 <> 0 THEN
        SET in_nRet=-1591;
    END IF;
    ELSE
        SELECT nLevel,nTypeID,nValue,nItemID FROM tGuildAcademyReward WHERE nNo=in_nNo AND nTypeID>1;
        IF 0 <> 0 THEN
        SET in_nRet=-1591;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_raiseMember
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_raiseMember`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_raiseMember`(
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
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setChatBanTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setChatBanTime`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setChatBanTime`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nChatBanTime INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nNo FROM tGuildAcademyMember WHERE nNo = in_nNo AND nCharNo = in_nCharNo ) THEN
        UPDATE tGuildAcademyMember SET nChatBanTime=in_nChatBanTime WHERE nNo=in_nNo AND nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setGraduate
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setGraduate`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setGraduate`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_sID VARCHAR(40),
    IN in_nClass TINYINT,
    IN in_nGraduateLevel SMALLINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nJoinLevel SMALLINT;
    DECLARE v_dJoinDate DATETIME;

    SET in_nRet=0;
    SET v_nJoinLevel=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    IF CHAR_LENGTH ( IFNULL ( in_sID , '' ) ) = 0 THEN
        SET in_nRet=-1515;
        LEAVE sp_body;
    END IF;
    IF in_nClass <= 0 THEN
        SET in_nRet=-1516;
        LEAVE sp_body;
    END IF;
    IF in_nGraduateLevel <= 0 THEN
        SET in_nRet=-1517;
        LEAVE sp_body;
    END IF;
    SELECT nJoinLevel,dJoinDate INTO v_nJoinLevel,v_dJoinDate FROM tGuildAcademyMember WHERE nNo=in_nNo AND nCharNo=in_nCharNo;
    IF v_nJoinLevel <= 0 THEN
        SET in_nRet=-1528;
        LEAVE sp_body;
    ELSE
        INSERT INTO tGuildAcademyGraduate(nNo,nCharNo,sID,nClass,nJoinLevel,dJoinDate,nGraduateLevel) VALUES (in_nNo,in_nCharNo,in_sID,in_nClass,v_nJoinLevel,v_dJoinDate,in_nGraduateLevel);
        IF 0 <> 0 THEN
        SET in_nRet=-1592;
        LEAVE sp_body;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setHistory
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setHistory`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setHistory`(
    IN in_nAcademyID INT,
    IN in_nCharNo INT,
    IN in_sID VARCHAR(50),
    IN in_nAction INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    INSERT INTO tGuildAcademyHistory(nAcademyID,nCharNo,sID,dDate,nAction) VALUES (in_nAcademyID,in_nCharNo,in_sID,CURRENT_TIMESTAMP,in_nAction);
    IF 0 <> 0 THEN
        SET out_nRet=-2;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setMaster
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setMaster`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setMaster`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nNo FROM tGuildMember WHERE nNo = in_nNo AND nCharNo = in_nCharNo ) THEN
        UPDATE tGuildAcademy SET nMasterCharNo=in_nCharNo WHERE nNo=in_nNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setMember
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setMember`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setMember`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_dJoinDate DATETIME,
    IN in_nGrade TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nLevel INT;

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nNo FROM tGuildMember WHERE nNo = in_nNo AND nCharNo = in_nCharNo ) THEN
        SET in_nRet=-1522;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted = 0 ) THEN
        SELECT nLevel INTO v_nLevel FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nLevel <= 0 THEN
        SET in_nRet=-1523;
        LEAVE sp_body;
    END IF;
    ELSE
        SET in_nRet=-1524;
        LEAVE sp_body;
    END IF;
    INSERT tGuildAcademyMember(nNo,nCharNo,nJoinLevel,dJoinDate,nGrade) VALUES (in_nNo,in_nCharNo,v_nLevel,in_dJoinDate,in_nGrade);
    IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1591;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setMemberLevelup
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setMemberLevelup`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setMemberLevelup`(
    IN in_nNo INT,
    IN in_nPoint INT,
    IN in_nBuffTime INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_BUFF_END_TIME DATETIME;

    SET in_nRet=0;
    IF EXISTS ( SELECT nNo FROM tGuildAcademy WHERE nNo = in_nNo ) THEN
        START TRANSACTION;
        SELECT IFNULL(dBuffEndTime,NOW()) INTO v_BUFF_END_TIME FROM tGuildAcademy WHERE nNo=in_nNo;
        IF v_BUFF_END_TIME > NOW ( ) THEN
        SET v_BUFF_END_TIME=DATE_ADD(v_BUFF_END_TIME,INTERVAL in_nBuffTime MINUTE);
    ELSE
        SET v_BUFF_END_TIME=DATE_ADD(NOW(),INTERVAL in_nBuffTime MINUTE);
    END IF;
        UPDATE tGuildAcademy SET nPoint=nPoint+in_nPoint,dBuffEndTime=v_BUFF_END_TIME WHERE nNo=in_nNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        COMMIT;
    END IF;
    ELSE
        SET in_nRet=-1521;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setNotify
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setNotify`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setNotify`(
    IN in_nNo INT,
    IN in_sNotify VARCHAR(1024),
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    UPDATE tGuildAcademy SET dNotifyDate=NOW(),sNotify=in_sNotify WHERE nNo=in_nNo;
    IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
    ELSE
        SET in_nRet=0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setRanking
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setRanking`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setRanking`(
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_MAX_COUNT INT;

    SET in_nRet=0;
    START TRANSACTION;
    SELECT IFNULL(MAX(nCount),0)+1 INTO v_MAX_COUNT FROM tGuildAcademyRanking;
    INSERT INTO tGuildAcademyRanking(dDate,nRanking,nGuildNo,sGuildName,nMasterNo,sMasterID,nPoint,nCount)SELECT NOW() AS dDate,R.nRanking,R.nNo,G.sName,GA.nMasterCharNo,C.sID,R.nPoint,v_MAX_COUNT FROM tGuild AS G,tGuildAcademy AS GA,tCharacter AS C,(SELECT COUNT(T2.nPoint)+1 AS nRanking,T1.nNo,T1.nPoint FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint<T2.nPoint GROUP BY T1.nNo,T1.nPoint) AS R WHERE G.nNo=R.nNo AND C.nCharNo=GA.nMasterCharNo AND GA.nNo=R.nNo;
    IF 0 <> 0 THEN
        SET in_nRet=-1592;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tGuildAcademy SET nPoint=0;
    IF ROW_COUNT ( ) = 0 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_setReward
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_setReward`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_setReward`(
    IN in_nNo INT,
    IN in_nLevel SMALLINT,
    IN in_nTypeID TINYINT,
    IN in_nValue BIGINT,
    IN in_nItemID INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nLevel <= 0 THEN
        SET in_nRet=-1513;
        LEAVE sp_body;
    END IF;
    INSERT INTO tGuildAcademyReward(nNo,nLevel,nTypeID,nValue,nItemID) VALUES (in_nNo,in_nLevel,in_nTypeID,in_nValue,in_nItemID);
    IF 0 <> 0 THEN
        SET in_nRet=-1595;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildAcademy_withdrawMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildAcademy_withdrawMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildAcademy_withdrawMoney`(
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_money BIGINT;

    SET in_nRet=0;
    SET v_money=0;
    IF in_nNo <= 0 THEN
        SET in_nRet=-1511;
        LEAVE sp_body;
    END IF;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1512;
        LEAVE sp_body;
    END IF;
    IF in_nMoney <= 0 THEN
        SET in_nRet=-1514;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    IF NOT EXISTS ( SELECT nCharNo FROM tCharacter WHERE nCharNo = in_nCharNo AND bDeleted = 0 ) THEN
        SET in_nRet=-1524;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nNo FROM tGuildAcademy WHERE nNo = in_nNo ) THEN
        SET in_nRet=-1525;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(nMoney,0) INTO v_money FROM tGuildAcademy WHERE nNo=in_nNo;
    IF v_money < in_nMoney THEN
        SET in_nRet=-1527;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        UPDATE tGuildAcademy SET nMoney=nMoney-in_nMoney WHERE nNo=in_nNo AND nMoney>=in_nMoney;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        UPDATE tCharacter SET nMoney=nMoney+in_nMoney WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1593;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildEmblem_GetInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildEmblem_GetInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildEmblem_GetInfo`(
    IN in_nGuildNo INT
)
sp_body: BEGIN

    SELECT nColor,nIcon,dCreated,dUpdated,nUnk FROM tGuildEmblem WHERE nGuildNo=in_nGuildNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildEmblem_GetInfoAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildEmblem_GetInfoAll`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildEmblem_GetInfoAll`()
sp_body: BEGIN

    SELECT nGuildNo,nColor,nIcon,dCreated,dUpdated,nUnk FROM tGuildEmblem END0;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildEmblem_SetInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildEmblem_SetInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildEmblem_SetInfo`(
    IN in_nGuildNo INT,
    IN in_nColor INT,
    IN in_nIcon INT,
    IN in_dCreated DATETIME,
    IN in_dUpdated DATETIME,
    IN in_nUnk INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    IF EXISTS ( SELECT nGuildNo FROM tGuildEmblem WHERE nGuildNo = in_nGuildNo ) THEN
        UPDATE tGuildEmblem SET nColor=in_nColor,nIcon=in_nIcon,dUpdated=in_dUpdated,nUnk=in_nUnk WHERE nGuildNo=in_nGuildNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tGuildEmblem(nGuildNo,nColor,nIcon,dCreated,dUpdated,nUnk) VALUES (in_nGuildNo,in_nColor,in_nIcon,in_dCreated,in_dUpdated,in_nUnk);
        IF 0 <> 0 THEN
        SET out_nRet=-2;
        LEAVE sp_body;
    END IF;
    END IF;
    IF 0 <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    LEAVE sp_body;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildTournament_getChoiceList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildTournament_getChoiceList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildTournament_getChoiceList`(
    IN in_nGTNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nGTNo , 0 ) <= 0 THEN
        SET in_nRet=-1411;
        LEAVE sp_body;
    END IF;
    SELECT nNo FROM tGuildTournamentRegist WHERE nGTNo=in_nGTNo AND bChoice=1;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildTournament_getMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildTournament_getMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildTournament_getMoney`(
    IN in_nGTNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nGTNo , 0 ) <= 0 THEN
        SET in_nRet=-1411;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(SUM(nMoney),0) AS nTotalMoney FROM tGuildTournamentRegist WHERE nGTNo=in_nGTNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildTournament_getRegistList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildTournament_getRegistList`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildTournament_getRegistList`(
    IN in_nGTNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nGTNo , 0 ) <= 0 THEN
        SET in_nRet=-1411;
        LEAVE sp_body;
    END IF;
    SELECT nNo FROM tGuildTournamentRegist WHERE nGTNo=in_nGTNo;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildTournament_setChoice
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildTournament_setChoice`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildTournament_setChoice`(
    IN in_nGTNo INT,
    IN in_nLimitLevel INT,
    IN in_nLimitMember INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF IFNULL ( in_nGTNo , 0 ) <= 0 THEN
        SET in_nRet=-1411;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    UPDATE tGuildTournamentRegist SET bChoice=1 WHERE nGTNo=in_nGTNo AND nNo IN(SELECT T4.nNo FROM tGuildTournamentRegist AS T4 INNER JOIN (SELECT T1.nNo,IFNULL(COUNT(T2.nCharNo),0) AS nMember FROM tGuildMember AS T1,tCharacter AS T2 WHERE T1.nCharNo=T2.nCharNo AND T2.nLevel>=in_nLimitLevel GROUP BY T1.nNo HAVING COUNT(T2.nCharNo)>=in_nLimitMember) AS T3 ON T4.nNo=T3.nNo WHERE nGTNo=in_nGTNo ORDER BY nSeed DESC) LIMIT 16;
    IF 0 <> 0 THEN
        SET in_nRet=-1493;
        ROLLBACK;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_GuildTournament_setRegist
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_GuildTournament_setRegist`;
DELIMITER $$
CREATE PROCEDURE `usp_GuildTournament_setRegist`(
    IN in_nGTNo INT,
    IN in_nNo INT,
    IN in_nCharNo INT,
    IN in_nMoney BIGINT,
    IN in_nLimitLevel INT,
    IN in_nLimitMember INT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_prevGTNo INT;
    DECLARE v_prevWinnerNo INT;
    DECLARE v_cnt INT;

    SET in_nRet=0;
    SET v_cnt=0;
    IF IFNULL ( in_nGTNo , 0 ) <= 0 THEN
        SET in_nRet=-1411;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nNo , 0 ) <= 0 THEN
        SET in_nRet=-1412;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-1413;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nGTNo FROM tGuildTournamentRegist WHERE nGTNo = in_nGTNo AND nNo = in_nNo ) THEN
        SET in_nRet=-1421;
        LEAVE sp_body;
    END IF;
    SELECT IFNULL(COUNT(T1.nCharNo),0) INTO v_cnt FROM tGuildMember AS T1,tCharacter AS T2 WHERE T1.nCharNo=T2.nCharNo AND T1.nNo=in_nNo AND T2.nLevel>=in_nLimitLevel;
    IF v_cnt < in_nLimitMember THEN
        SET in_nRet=-1422;
        LEAVE sp_body;
    END IF;
    SELECT nGTNo INTO v_prevGTNo FROM tGuildTournament WHERE nGTNo<in_nGTNo ORDER BY nGTNo DESC LIMIT 1;
    IF v_prevGTNo IS NOT NULL THEN
        SELECT nGuildNo INTO v_prevWinnerNo FROM tGuildTournament WHERE nGTNo=v_prevGTNo AND nType>=3028 AND nType<=3029 AND nData=1;
        IF v_prevWinnerNo = in_nNo THEN
        SET in_nRet=-1424;
        LEAVE sp_body;
    END IF;
    END IF;
    START TRANSACTION;
    CALL p_Char_MoneySub(in_nCharNo,in_nMoney,in_nRet);
    IF in_nRet <> 1 THEN
        SET in_nRet=-1423;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    INSERT INTO tGuildTournamentRegist(nGTNo,nNo,nMoney,nSeed) VALUES (in_nGTNo,in_nNo,in_nMoney,CAST(RAND()*100000 AS SIGNED));
    IF 0 <> 0 THEN
        SET in_nRet=-1492;
        ROLLBACK;
    ELSE
        SET in_nRet=0;
        COMMIT;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Guild_GetToken
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Guild_GetToken`;
DELIMITER $$
CREATE PROCEDURE `usp_Guild_GetToken`(
    IN in_nNo INT,
    OUT out_nToken BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    SELECT IFNULL(MAX(nToken),0) INTO out_nToken FROM tGuildToken WHERE nNo=in_nNo;
    IF 0 <> 0 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;

    SELECT out_nToken AS nToken, out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Guild_GetTokenAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Guild_GetTokenAll`;
DELIMITER $$
CREATE PROCEDURE `usp_Guild_GetTokenAll`()
sp_body: BEGIN

    SELECT nNo,nToken FROM tGuildToken;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Guild_SetToken
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Guild_SetToken`;
DELIMITER $$
CREATE PROCEDURE `usp_Guild_SetToken`(
    IN in_nNo INT,
    IN in_nToken BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    IF EXISTS ( SELECT nNo FROM tGuildToken WHERE nNo = in_nNo ) THEN
        UPDATE tGuildToken SET nToken=in_nToken WHERE nNo=in_nNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tGuildToken(nNo,nToken) VALUES (in_nNo,in_nToken);
        IF 0 <> 0 THEN
        SET out_nRet=-2;
        LEAVE sp_body;
    END IF;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_HolyPromise_getRemainMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_HolyPromise_getRemainMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_HolyPromise_getRemainMoney`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nRemainMoney BIGINT;

    SET v_nRemainMoney=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET out_nRet=-2211;
        LEAVE sp_body;
    END IF;
    SELECT nRemainMoney INTO v_nRemainMoney FROM tHolyPromiseReward WHERE nCharNo=in_nCharNo ORDER BY nSeqNo DESC LIMIT 1;
    IF v_nRemainMoney > 0 THEN
        SELECT v_nRemainMoney;
        SET out_nRet=0;
        LEAVE sp_body;
    ELSE
        IF EXISTS ( SELECT nUpCharNo FROM tHolyPromise WHERE nUpCharNo = in_nCharNo ) THEN
        SELECT 0;
        SET out_nRet=0;
        LEAVE sp_body;
    ELSE
        SET out_nRet=-2224;
        LEAVE sp_body;
    END IF;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_HolyPromise_setMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_HolyPromise_setMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_HolyPromise_setMoney`(
    IN in_nCharNo INT,
    IN in_nRewardMoney BIGINT,
    IN in_nLevelGap INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nUpMemberCharNo INT;
    DECLARE v_nUpMemberLevel INT;
    DECLARE v_nMyLevel INT;
    DECLARE v_nType TINYINT;
    DECLARE v_nSeqNo INT;
    DECLARE v_nRemainMoney BIGINT;

    SET v_nUpMemberCharNo=0;
    SET v_nUpMemberLevel=0;
    SET v_nMyLevel=0;
    SET v_nType=1;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET out_nRet=-2211;
        LEAVE sp_body;
    END IF;
    IF IFNULL ( in_nRewardMoney , 0 ) <= 0 THEN
        SET out_nRet=-2212;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nUpCharNo FROM tHolyPromise WHERE nCharNo = in_nCharNo ) THEN
        SELECT P.nUpCharNo,C.nLevel INTO v_nUpMemberCharNo,v_nUpMemberLevel FROM tHolyPromise AS P INNER JOIN tCharacter AS C ON C.nCharNo=P.nUpCharNo WHERE P.nCharNo=in_nCharNo AND C.bDeleted=0;
        IF v_nUpMemberCharNo > 0 and v_nUpMemberLevel > 0 THEN
        SELECT nLevel INTO v_nMyLevel FROM tCharacter WHERE nCharNo=in_nCharNo;
        IF v_nUpMemberLevel - v_nMyLevel < in_nLevelGap THEN
        SET out_nRet=-2223;
        LEAVE sp_body;
    END IF;
        SET v_nSeqNo=0;
        SET v_nRemainMoney=0;
        START TRANSACTION;
        SELECT nSeqNo,nRemainMoney INTO v_nSeqNo,v_nRemainMoney FROM tHolyPromiseReward WHERE nCharNo=v_nUpMemberCharNo ORDER BY nSeqNo DESC LIMIT 1;
        SET v_nSeqNo=v_nSeqNo+1;
        SET v_nRemainMoney=v_nRemainMoney+in_nRewardMoney;
        INSERT INTO tHolyPromiseReward(nCharNo,nSeqNo,nType,nMoney,nRemainMoney,nReceiptCharNo) VALUES (v_nUpMemberCharNo,v_nSeqNo,v_nType,in_nRewardMoney,v_nRemainMoney,in_nCharNo);
        IF 0 <> 0 THEN
        SET out_nRet=-2292;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        SET out_nRet=0;
        COMMIT;
        LEAVE sp_body;
    END IF;
    ELSE
        SET out_nRet=-2222;
        LEAVE sp_body;
    END IF;
    ELSE
        SET out_nRet=-2221;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_HolyPromise_withDrawRemainMoney
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_HolyPromise_withDrawRemainMoney`;
DELIMITER $$
CREATE PROCEDURE `usp_HolyPromise_withDrawRemainMoney`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_nType TINYINT;
    DECLARE v_nRewardMoney BIGINT;
    DECLARE v_nSeqNo INT;

    SET v_nType=2;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET out_nRet=-2211;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nMoney FROM tHolyPromiseReward WHERE nCharNo = in_nCharNo ) THEN
        START TRANSACTION;
        SELECT nSeqNo,nRemainMoney INTO v_nSeqNo,v_nRewardMoney FROM tHolyPromiseReward WHERE nCharNo=in_nCharNo ORDER BY nSeqNo DESC LIMIT 1;
        IF v_nRewardMoney <= 0 THEN
        SET out_nRet=-2226;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        CALL p_Char_MoneyAdd(in_nCharNo,v_nRewardMoney,out_nRet);
        IF out_nRet <> 1 THEN
        SET out_nRet=-2227;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SET v_nSeqNo=v_nSeqNo+1;
        INSERT INTO tHolyPromiseReward(nCharNo,nSeqNo,nType,nMoney,nRemainMoney,nReceiptCharNo) VALUES (in_nCharNo,v_nSeqNo,v_nType,v_nRewardMoney,0,0);
        IF 0 <> 0 THEN
        SET out_nRet=-2292;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
        SET out_nRet=0;
        SELECT v_nRewardMoney;
        COMMIT;
    ELSE
        SET out_nRet=-2225;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_InGameGuide_getViewList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_InGameGuide_getViewList`;
DELIMITER $$
CREATE PROCEDURE `usp_InGameGuide_getViewList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tInGameGuide WHERE nCharNo = in_nCharNo ) THEN
        SELECT nGuideNo FROM tInGameGuide WHERE nCharNo=in_nCharNo;
    ELSE
        SET in_nRet=-1111;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_InGameGuide_setGuideNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_InGameGuide_setGuideNo`;
DELIMITER $$
CREATE PROCEDURE `usp_InGameGuide_setGuideNo`(
    IN in_nCharNo INT,
    IN in_nGuideNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tInGameGuide WHERE nCharNo = in_nCharNo AND nGuideNo = in_nGuideNo ) THEN
        SET in_nRet=-1111;
    ELSE
        INSERT INTO tInGameGuide(nCharNo,nGuideNo) VALUES (in_nCharNo,in_nGuideNo);
        IF 0 <> 0 THEN
        SET in_nRet=-2222;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_Add_EnchantGem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_Add_EnchantGem`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_Add_EnchantGem`(
    IN in_nCharNo INT,
    IN in_nEquipmentItemKey BIGINT,
    IN in_nGemItemKey BIGINT,
    IN in_nGemSocketIndex TINYINT,
    IN in_nGemID INT,
    IN in_nRestGemCount TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    START TRANSACTION;
    INSERT tItemEnchant(nItemKey,nSlot,nGemItemID,nRestCount) VALUES (in_nEquipmentItemKey,in_nGemSocketIndex,in_nGemID,in_nRestGemCount);
    IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-101;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Delete(in_nGemItemKey,in_nRet);
    IF 0 <> 0 OR in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-202;
        LEAVE sp_body;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_Add_EnchantSocket
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_Add_EnchantSocket`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_Add_EnchantSocket`(
    IN in_nCharNo INT,
    IN in_nEquipmentItemKey BIGINT,
    IN in_nDrillItemKey BIGINT,
    OUT out_nAllSocketCount TINYINT,
    IN in_nDrillLotType SMALLINT,
    OUT out_nDrillLotValue BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_nCreatedSocketCount TINYINT;

    SET in_nRet=0;
    START TRANSACTION;
    CALL p_Item_Delete(in_nDrillItemKey,in_nRet);
    IF 0 <> 0 OR in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-101;
        LEAVE sp_body;
    END IF;
    SET out_nDrillLotValue=0;
    CALL p_Item_GetOption(in_nEquipmentItemKey,523,out_nAllSocketCount,in_nRet);
    IF 0 <> 0 OR in_nRet = 0 THEN
        ROLLBACK;
        SET in_nRet=-201;
        LEAVE sp_body;
    END IF;
    CALL p_Item_GetOption(in_nEquipmentItemKey,524,v_nCreatedSocketCount,in_nRet);
    IF 0 <> 0 OR in_nRet = 0 THEN
        ROLLBACK;
        SET in_nRet=-202;
        LEAVE sp_body;
    END IF;
    SET out_nAllSocketCount=out_nAllSocketCount+1;
    SET v_nCreatedSocketCount=v_nCreatedSocketCount+1;
    CALL p_Item_SetOption(in_nEquipmentItemKey,523,out_nAllSocketCount,in_nRet);
    IF 0 <> 0 OR in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-301;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nEquipmentItemKey,524,v_nCreatedSocketCount,in_nRet);
    IF 0 <> 0 OR in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-302;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nAllSocketCount AS nAllSocketCount, out_nDrillLotValue AS nDrillLotValue;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_Del_EnchantGem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_Del_EnchantGem`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_Del_EnchantGem`(
    IN in_nCharNo INT,
    IN in_nEquipmentItemKey BIGINT,
    IN in_nRemoverItemKey BIGINT,
    IN in_nResetSocketIndex TINYINT,
    IN in_bNeedItem TINYINT,
    IN in_nRemoverLotType SMALLINT,
    IN in_nRemoverLotValue BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SET in_nRemoverLotValue=0;
    START TRANSACTION;
    IF in_bNeedItem = 1 THEN
        CALL p_Item_Delete(in_nRemoverItemKey,in_nRet);
        IF 0 <> 0 OR in_nRet <> 0 THEN
        ROLLBACK;
        SET in_nRet=-101;
        LEAVE sp_body;
    END IF;
    END IF;
    DELETE FROM tItemEnchant WHERE nItemKey=in_nEquipmentItemKey AND nSlot=in_nResetSocketIndex;
    IF ROW_COUNT ( ) = 0 THEN
        ROLLBACK;
        SET in_nRet=-201;
        LEAVE sp_body;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_Get_EnchantAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_Get_EnchantAll`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_Get_EnchantAll`(
    IN in_nEquipmentItemKey BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nSlot,nGemItemID,nRestCount FROM tItemEnchant WHERE nItemKey=in_nEquipmentItemKey;
    IF 0 <> 0 THEN
        SET in_nRet=-101;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_GuildToken_BuyAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_GuildToken_BuyAll`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_GuildToken_BuyAll`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    IN in_nItemID INT,
    IN in_nFlags INT,
    IN in_nPriceGuildToken BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    DECLARE v_CreateItemRet INT;

    START TRANSACTION;
    CALL usp_Guild_TokenSub(in_nOwner,in_nPriceGuildToken,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Create(in_nItemKey,in_nOwner,in_nStorageType,in_nStorage,in_nItemID,in_nFlags,v_CreateItemRet);
    IF v_CreateItemRet <> 0 OR 0 <> 0 THEN
        ROLLBACK;
        SET out_nRet=v_CreateItemRet;
        LEAVE sp_body;
    END IF;
    COMMIT;
    SET out_nRet=0;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_GuildToken_BuyLot
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_GuildToken_BuyLot`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_GuildToken_BuyLot`(
    IN in_nItemKey BIGINT,
    IN in_nOwner INT,
    IN in_nOptionType SMALLINT,
    IN in_nOptionData BIGINT,
    IN in_nPriceGuildToken BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    START TRANSACTION;
    CALL usp_Guild_TokenSub(in_nOwner,in_nPriceGuildToken,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        SET out_nRet=1;
        LEAVE sp_body;
    END IF;
    CALL p_Item_SetOption(in_nItemKey,in_nOptionType,in_nOptionData,out_nRet);
    IF out_nRet <> 0 THEN
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_Set_EnchantGemLot
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_Set_EnchantGemLot`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_Set_EnchantGemLot`(
    IN in_nCharNo INT,
    IN in_nEquipmentItemKey BIGINT,
    IN in_nSlot TINYINT,
    IN in_nRestGemCount TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    UPDATE tItemEnchant SET nRestCount=in_nRestGemCount WHERE nItemKey=in_nEquipmentItemKey AND nSlot=in_nSlot;
    IF ROW_COUNT ( ) <> 1 OR 0 <> 0 THEN
        SET in_nRet=-101;
        LEAVE sp_body;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_getAutoMonList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getAutoMonList`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getAutoMonList`(
    IN in_nCharNo INT,
    IN in_nType INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    SELECT nSlot,nItemID FROM tAutoMon WHERE nCharNo=in_nCharNo AND nType=in_nType;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_setAutoMon
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setAutoMon`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setAutoMon`(
    IN in_nCharNo INT,
    IN in_nType INT,
    IN in_nSlot INT,
    IN in_nItemID INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    START TRANSACTION;
    IF EXISTS ( SELECT nCharNo , nType , nSlot FROM tAutoMon WHERE nCharNo = in_nCharNo AND nType = in_nType AND nSlot = in_nSlot ) THEN
        UPDATE tAutoMon SET nItemID=in_nItemID WHERE nCharNo=in_nCharNo AND nType=in_nType AND nSlot=in_nSlot;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tAutoMon(nCharNo,nType,nSlot,nItemID) VALUES (in_nCharNo,in_nType,in_nSlot,in_nItemID);
        IF 0 <> 0 THEN
        SET in_nRet=1;
        IF 1 > 0 THEN
        ROLLBACK;
    END IF;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_setPutOnBelonged
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setPutOnBelonged`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setPutOnBelonged`(
    IN in_nItemKey BIGINT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    IN in_nOptionType SMALLINT,
    IN in_nOptionData BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    UPDATE tItem SET nStorageType=in_nStorageType,nStorage=in_nStorage WHERE nItemKey=in_nItemKey;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItemOptions WHERE nItemKey = in_nItemKey AND nOptionType = in_nOptionType ) THEN
        UPDATE tItemOptions SET nOptionType=in_nOptionType,nOptionData=in_nOptionData WHERE nItemKey=in_nItemKey AND nOptionType=in_nOptionType;
        IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tItemOptions(nItemKey,nOptionType,nOptionData) VALUES (in_nItemKey,in_nOptionType,in_nOptionData);
        IF 0 <> 0 THEN
        SET out_nRet=-2592;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_setPutOnBelongedCancel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setPutOnBelongedCancel`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setPutOnBelongedCancel`(
    IN in_nCharNo INT,
    IN in_nCancelItemKey BIGINT,
    IN in_nCancelOptionType SMALLINT,
    IN in_nCancelOptionData BIGINT,
    IN in_nTargetItemKey BIGINT,
    IN in_nTargetOptionType SMALLINT,
    IN in_nTargetOptionData BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_cancelOptionData BIGINT;
    DECLARE v_targetOptionData BIGINT;

    SET in_nRet=0;
    SET v_cancelOptionData=0;
    SET v_targetOptionData=0;
    IF IFNULL ( in_nCharNo , 0 ) <= 0 THEN
        SET in_nRet=-2511;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nCancelItemKey AND nOwner = in_nCharNo ) THEN
        SET in_nRet=-2512;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nItemKey FROM tItem WHERE nItemKey = in_nTargetItemKey AND nOwner = in_nCharNo ) THEN
        SET in_nRet=-2513;
        LEAVE sp_body;
    END IF;
    START TRANSACTION;
    SELECT nOptionData INTO v_cancelOptionData FROM tItemOptions WHERE nItemKey=in_nCancelItemKey AND nOptionType=in_nCancelOptionType;
    IF in_nCancelOptionData <> v_cancelOptionData THEN
        SET in_nRet=-2521;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF v_cancelOptionData > 1 THEN
        SET v_cancelOptionData=v_cancelOptionData-1;
        CALL p_Item_SetOption(in_nCancelItemKey,in_nCancelOptionType,v_cancelOptionData,in_nRet);
        IF in_nRet <> 0 THEN
        SET in_nRet=-2522;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        CALL p_Item_Delete(in_nCancelItemKey,in_nRet);
        IF in_nRet <> 0 THEN
        SET in_nRet=-2523;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    SELECT nOptionData INTO v_targetOptionData FROM tItemOptions WHERE nItemKey=in_nTargetItemKey AND nOptionType=in_nTargetOptionType;
    IF v_targetOptionData <> 1 THEN
        SET in_nRet=-2524;
        ROLLBACK;
        LEAVE sp_body;
    ELSE
        CALL p_Item_SetOption(in_nTargetItemKey,in_nTargetOptionType,in_nTargetOptionData,in_nRet);
        IF in_nRet <> 0 THEN
        SET in_nRet=-2525;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Item_setPutOnBelongedExchange
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setPutOnBelongedExchange`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setPutOnBelongedExchange`(
    IN in_nItemKey1 BIGINT,
    IN in_nStorageType1 TINYINT,
    IN in_nStorage1 SMALLINT,
    IN in_nItemKey2 BIGINT,
    IN in_nStorageType2 TINYINT,
    IN in_nStorage2 SMALLINT,
    IN in_nOptionType1 SMALLINT,
    IN in_nOptionData1 BIGINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SET out_nRet=0;
    START TRANSACTION;
    UPDATE tItem SET nStorageType=11 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType2,nStorage=in_nStorage2 WHERE nItemKey=in_nItemKey2;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    UPDATE tItem SET nStorageType=in_nStorageType1,nStorage=in_nStorage1 WHERE nItemKey=in_nItemKey1;
    IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    IF EXISTS ( SELECT nItemKey FROM tItemOptions WHERE nItemKey = in_nItemKey1 AND nOptionType = in_nOptionType1 ) THEN
        UPDATE tItemOptions SET nOptionType=in_nOptionType1,nOptionData=in_nOptionData1 WHERE nItemKey=in_nItemKey1 AND nOptionType=in_nOptionType1;
        IF ROW_COUNT ( ) = 0 THEN
        SET out_nRet=-2593;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    ELSE
        INSERT INTO tItemOptions(nItemKey,nOptionType,nOptionData) VALUES (in_nItemKey1,in_nOptionType1,in_nOptionData1);
        IF 0 <> 0 THEN
        SET out_nRet=-2592;
        ROLLBACK;
        LEAVE sp_body;
    END IF;
    END IF;
    COMMIT;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Minihouse_delGate
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Minihouse_delGate`;
DELIMITER $$
CREATE PROCEDURE `usp_Minihouse_delGate`(
    IN in_nCharNo INT,
    IN in_nItemID INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = in_nCharNo AND nItemID = in_nItemID ) THEN
        DELETE FROM tMinihouseGate WHERE nCharNo=in_nCharNo AND nItemID=in_nItemID;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1694;
    END IF;
    ELSE
        SET in_nRet=-1622;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Minihouse_getGateList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Minihouse_getGateList`;
DELIMITER $$
CREATE PROCEDURE `usp_Minihouse_getGateList`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF in_nCharNo <= 0 THEN
        SET in_nRet=-1611;
        LEAVE sp_body;
    END IF;
    IF NOT EXISTS ( SELECT nCharNo FROM tMinihouse WHERE nCharNo = in_nCharNo ) THEN
        INSERT INTO tMinihouse(nCharNo) VALUES (in_nCharNo);
        IF 0 <> 0 THEN
        SET in_nRet=-1692;
        LEAVE sp_body;
    END IF;
    END IF;
    SELECT nCharNo,nItemID,bType,IFNULL(dEndDate,0) AS dEndDate FROM tMinihouseGate WHERE nCharNo=in_nCharNo AND (dEndDate>NOW() OR dEndDate IS NULL);

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Minihouse_setGate
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Minihouse_setGate`;
DELIMITER $$
CREATE PROCEDURE `usp_Minihouse_setGate`(
    IN in_nCharNo INT,
    IN in_nItemID INT,
    IN in_bType TINYINT,
    IN in_dEndDate DATETIME,
    IN in_nItemKey BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = in_nCharNo AND nItemID = in_nItemID ) THEN
        DELETE FROM tMinihouseGate WHERE nCharNo=in_nCharNo AND nItemID=in_nItemID AND dEndDate<NOW() AND dEndDate IS NOT NULL;
    END IF;
    IF EXISTS ( SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = in_nCharNo AND nItemID = in_nItemID ) THEN
        SET in_nRet=-1621;
    ELSE
        IF in_bType = 0 THEN
        SET in_dEndDate=NULL;
    END IF;
        START TRANSACTION;
        CALL p_Item_Delete(in_nItemKey,in_nRet);
        IF in_nRet = 0 THEN
        INSERT INTO tMinihouseGate(nCharNo,nItemID,bType,dEndDate) VALUES (in_nCharNo,in_nItemID,in_bType,in_dEndDate);
        IF 0 <> 0 THEN
        SET in_nRet=-1692;
        ROLLBACK;
    ELSE
        COMMIT;
    END IF;
    ELSE
        SET in_nRet=-1622;
        ROLLBACK;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Minihouse_setVisitorCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Minihouse_setVisitorCount`;
DELIMITER $$
CREATE PROCEDURE `usp_Minihouse_setVisitorCount`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tMinihouse WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tMinihouse SET nVisitor=nVisitor+1 WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1693;
    END IF;
    ELSE
        INSERT INTO tMinihouse(nCharNo) VALUES (in_nCharNo);
        IF 0 <> 0 THEN
        SET in_nRet=-1692;
        LEAVE sp_body;
    END IF;
        UPDATE tMinihouse SET nVisitor=nVisitor+1 WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-1693;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Pet_getInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pet_getInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_Pet_getInfo`(
    IN in_nPetNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN

    SELECT nOwnerNo,sPetName,nPetMind,nPetStress FROM tPetNew WHERE nPetNo=in_nPetNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Pet_setCreate
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pet_setCreate`;
DELIMITER $$
CREATE PROCEDURE `usp_Pet_setCreate`(
    IN in_nOwerNo INT,
    IN in_sPetName VARCHAR(16),
    IN in_nPetMind TINYINT,
    IN in_nPetStress TINYINT,
    OUT out_nRet INT
)
sp_body: BEGIN

    INSERT INTO tPetNew(nOwnerNo,sPetName,nPetMind,nPetStress) VALUES (in_nOwerNo,in_sPetName,in_nPetMind,in_nPetStress);
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;
    SELECT LAST_INSERT_ID();

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Pet_setName
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pet_setName`;
DELIMITER $$
CREATE PROCEDURE `usp_Pet_setName`(
    IN in_nPetNo INT,
    IN in_sPetName VARCHAR(16),
    OUT out_nRet INT
)
sp_body: BEGIN

    UPDATE tPetNew SET sPetName=in_sPetName WHERE nPetNo=in_nPetNo;
    SET out_nRet=ROW_COUNT();
    IF ROW_COUNT ( ) <> 1 THEN
        SET out_nRet=-1;
        LEAVE sp_body;
    END IF;

    SELECT out_nRet AS nRet;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SafeTrade_setAll
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SafeTrade_setAll`;
DELIMITER $$
CREATE PROCEDURE `usp_SafeTrade_setAll`(
    IN in_nBuyer INT,
    IN in_nSeller INT,
    IN in_nMoveMoney BIGINT,
    IN in_nItemKey BIGINT,
    IN in_nStorageType TINYINT,
    IN in_nStorage SMALLINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_ret INT;
    DECLARE v_nOwner INT;
    DECLARE v_nOwnerStorageType TINYINT;

    SET v_ret=0;
    SET v_nOwner=0;
    SET v_nOwnerStorageType=0;
    SELECT nOwner,nStorageType INTO v_nOwner,v_nOwnerStorageType FROM tItem WHERE nItemKey=in_nItemKey;
    IF v_nOwner <> in_nSeller THEN
        SET in_nRet=-2621;
        LEAVE sp_body;
    END IF;
    IF v_nOwnerStorageType <> 9 THEN
        SET in_nRet=-2622;
        LEAVE sp_body;
    ELSE
        START TRANSACTION;
        CALL p_Char_MoneyMove(in_nBuyer,in_nSeller,in_nMoveMoney,v_ret);
        IF v_ret = 0 THEN
        ROLLBACK;
        SET in_nRet=-2623;
        LEAVE sp_body;
    END IF;
        CALL p_Item_SetStorageAndOwner(in_nItemKey,in_nBuyer,in_nStorageType,in_nStorage,v_ret);
        IF v_ret = 0 THEN
        ROLLBACK;
        SET in_nRet=-2624;
        LEAVE sp_body;
    ELSE
        COMMIT;
        SET in_nRet=0;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SafeTrade_setLot
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SafeTrade_setLot`;
DELIMITER $$
CREATE PROCEDURE `usp_SafeTrade_setLot`(
    IN in_nBuyer INT,
    IN in_nBuyerItemKey BIGINT,
    IN in_nBuyerStorageType TINYINT,
    IN in_nBuyerStorage SMALLINT,
    IN in_nBuyerItemID INT,
    IN in_nBuyerFlags INT,
    IN in_nBuyerOptionType SMALLINT,
    IN in_nSeller INT,
    IN in_nSellerItemKey BIGINT,
    IN in_nSellerOptionType SMALLINT,
    IN in_nMoveMoney BIGINT,
    IN in_nSplitLot BIGINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_ret INT;
    DECLARE v_nBuyerOptionData BIGINT;
    DECLARE v_nSellerOptionData BIGINT;
    DECLARE v_bIsDelete TINYINT;

    SET v_ret=0;
    SET v_nBuyerOptionData=0;
    SET v_nSellerOptionData=0;
    SET v_bIsDelete=0;
    START TRANSACTION;
    INSERT INTO tItem(nItemKey,nOwner,nStorageType,nStorage,nItemID,nFlags) VALUES (in_nBuyerItemKey,in_nBuyer,in_nBuyerStorageType,in_nBuyerStorage,in_nBuyerItemID,in_nBuyerFlags);
    IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-2628;
        LEAVE sp_body;
    ELSE
        INSERT INTO tItemOptions(nItemKey,nOptionType,nOptionData) VALUES (in_nBuyerItemKey,in_nBuyerOptionType,v_nBuyerOptionData);
        IF 0 <> 0 THEN
        ROLLBACK;
        SET in_nRet=-2629;
        LEAVE sp_body;
    END IF;
    END IF;
    CALL usp_SafeTrade_setMerge(in_nBuyer,in_nBuyerItemKey,in_nBuyerOptionType,in_nSeller,in_nSellerItemKey,in_nSellerOptionType,in_nMoveMoney,in_nSplitLot,v_bIsDelete,v_ret);
    IF v_ret <> 0 THEN
        ROLLBACK;
        SET in_nRet=v_ret;
        LEAVE sp_body;
    ELSE
        COMMIT;
        SET in_nRet=0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SafeTrade_setMerge
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SafeTrade_setMerge`;
DELIMITER $$
CREATE PROCEDURE `usp_SafeTrade_setMerge`(
    IN in_nBuyer INT,
    IN in_nBuyerItemKey BIGINT,
    IN in_nBuyerOptionType SMALLINT,
    IN in_nSeller INT,
    IN in_nSellerItemKey BIGINT,
    IN in_nSellerOptionType SMALLINT,
    IN in_nMoveMoney BIGINT,
    IN in_nSplitLot BIGINT,
    IN in_bIsDelete TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    DECLARE v_ret INT;
    DECLARE v_nBuyerOptionData BIGINT;
    DECLARE v_nSellerOptionData BIGINT;
    DECLARE v_nSellerLot TINYINT;

    SET v_ret=0;
    SET v_nBuyerOptionData=0;
    SET v_nSellerOptionData=0;
    SET v_nSellerLot=0;
    START TRANSACTION;
    CALL p_Char_MoneyMove(in_nBuyer,in_nSeller,in_nMoveMoney,v_ret);
    IF v_ret = 0 THEN
        ROLLBACK;
        SET in_nRet=-2623;
        LEAVE sp_body;
    END IF;
    CALL p_Item_Split_N_Merge_New(in_nSellerItemKey,in_nSeller,in_nSellerOptionType,v_nSellerOptionData,in_nBuyerItemKey,in_nBuyer,in_nBuyerOptionType,v_nBuyerOptionData,in_nSplitLot,v_ret);
    IF v_ret <> 0 THEN
        ROLLBACK;
        SET in_nRet=-2624;
        LEAVE sp_body;
    END IF;
    IF v_nSellerOptionData > 0 THEN
        SET v_nSellerLot=1;
    END IF;
    IF in_bIsDelete = v_nSellerLot THEN
        IF in_bIsDelete = 0 THEN
        ROLLBACK;
        SET in_nRet=-2625;
    ELSE
        ROLLBACK;
        SET in_nRet=-2626;
    END IF;
    ELSE
        IF in_bIsDelete = 1 THEN
        CALL p_Item_Delete(in_nSellerItemKey,v_ret);
        IF v_ret <> 0 THEN
        ROLLBACK;
        SET in_nRet=-2627;
    END IF;
    END IF;
        COMMIT;
        SET in_nRet=0;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SpamBlock_getTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SpamBlock_getTime`;
DELIMITER $$
CREATE PROCEDURE `usp_SpamBlock_getTime`(
    IN in_nCharNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = in_nCharNo ) THEN
        SELECT nCharNo,nBlockTime,nStandbyTime FROM tSpamBlockTime WHERE nCharNo=in_nCharNo;
    ELSE
        SET in_nRet=-2822;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SpamBlock_setTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SpamBlock_setTime`;
DELIMITER $$
CREATE PROCEDURE `usp_SpamBlock_setTime`(
    IN in_nCharNo INT,
    IN in_nBlockTime INT,
    IN in_nStandbyTime INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = in_nCharNo ) THEN
        UPDATE tSpamBlockTime SET sys_update_date=NOW(),nBlockTime=in_nBlockTime,nStandbyTime=in_nStandbyTime WHERE nCharNo=in_nCharNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2893;
    END IF;
    ELSE
        INSERT INTO tSpamBlockTime(nCharNo,nBlockTime,nStandbyTime) VALUES (in_nCharNo,in_nBlockTime,in_nStandbyTime);
        IF 0 <> 0 THEN
        SET in_nRet=-2892;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SpamFilter_get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SpamFilter_get`;
DELIMITER $$
CREATE PROCEDURE `usp_SpamFilter_get`(
    IN in_nUserNo INT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nUserNo FROM tSpamFilter WHERE nUserNo = in_nUserNo ) THEN
        SELECT nUserNo,bIsSpamSpread,bIsSpamBlock,nSpamFilterRate FROM tSpamFilter WHERE nUserNo=in_nUserNo;
    ELSE
        SET in_nRet=-2821;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_SpamFilter_set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_SpamFilter_set`;
DELIMITER $$
CREATE PROCEDURE `usp_SpamFilter_set`(
    IN in_nUserNo INT,
    IN in_bIsSpamSpread TINYINT,
    IN in_bIsSpamBlock TINYINT,
    IN in_nSpamFilterRate TINYINT,
    IN in_nRet INT
)
sp_body: BEGIN

    SET in_nRet=0;
    IF EXISTS ( SELECT nUserNo FROM tSpamFilter WHERE nUserNo = in_nUserNo ) THEN
        UPDATE tSpamFilter SET sys_update_date=NOW(),bIsSpamSpread=in_bIsSpamSpread,bIsSpamBlock=in_bIsSpamBlock,nSpamFilterRate=in_nSpamFilterRate WHERE nUserNo=in_nUserNo;
        IF ROW_COUNT ( ) <> 1 THEN
        SET in_nRet=-2893;
    END IF;
    ELSE
        INSERT INTO tSpamFilter(nUserNo,bIsSpamSpread,bIsSpamBlock,nSpamFilterRate) VALUES (in_nUserNo,in_bIsSpamSpread,in_bIsSpamBlock,in_nSpamFilterRate);
        IF 0 <> 0 THEN
        SET in_nRet=-2892;
    END IF;
    END IF;

END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Tutorial_getCharInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Tutorial_getCharInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_Tutorial_getCharInfo`(
    IN in_nCharNo INT
)
sp_body: BEGIN
    SELECT nState, nStep FROM tTutorial WHERE nCharNo = in_nCharNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Tutorial_getIsSkip
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Tutorial_getIsSkip`;
DELIMITER $$
CREATE PROCEDURE `usp_Tutorial_getIsSkip`(
    IN in_nUserNo INT,
    OUT out_nRet INT
)
sp_body: BEGIN
    SET out_nRet = 0;

    IF EXISTS(SELECT 1 FROM tTutorial WHERE nUserNo = in_nUserNo AND nState >= 1) THEN
        SET out_nRet = 1;
    ELSE
        SET out_nRet = 0;
    END IF;

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- Procedure: usp_Tutorial_setCharInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Tutorial_setCharInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_Tutorial_setCharInfo`(
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
        INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep) VALUES (in_nUserNo, in_nCharNo, in_nState, in_nStep);
        IF ROW_COUNT() <> 1 THEN SET out_nRet = -1; END IF;
    END IF;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;
