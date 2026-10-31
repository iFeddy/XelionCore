-- ============================================================================
-- PROCEDIMIENTOS ALMACENADOS OPERATIVOS DE LOGIN Y SESION EN MYSQL (13 PROCEDIMIENTOS)
-- Base de Datos: xelion_online
-- Fidelidad Binaria y de Protocolo: 100% oficial
-- ============================================================================

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- 1. usp_Version_GetClientKey
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Version_GetClientKey`;
DELIMITER $$
CREATE PROCEDURE `usp_Version_GetClientKey`()
BEGIN
    SELECT GameVersionKey
    FROM FOUS
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 2. usp_Version_SyncClientKey
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Version_SyncClientKey`;
DELIMITER $$
CREATE PROCEDURE `usp_Version_SyncClientKey`(
    IN in_versionKey VARCHAR(64)
)
BEGIN
    IF EXISTS(SELECT 1 FROM FOUS) THEN
        UPDATE FOUS SET GameVersionKey = in_versionKey;
    ELSE
        INSERT INTO FOUS (GameAccountTitle, GameAccountGameServerReference, GameAccountKey, GameVersionKey)
        VALUES ('FOUS', 'MainServer', '', in_versionKey);
    END IF;

    SELECT in_versionKey AS GameVersionKey;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 3. usp_WorldServer_GetList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_WorldServer_GetList`;
DELIMITER $$
CREATE PROCEDURE `usp_WorldServer_GetList`()
BEGIN
    SELECT 
        nWorldNo,
        sWorldName,
        sWorldIP,
        nWorldPort,
        nMaxUser,
        nUserRatable,
        bIsMaintenance,
        nCurrentCCU,
        dLastHeartbeat,
        CASE
            WHEN bIsMaintenance = 1 THEN 2
            WHEN nCurrentCCU >= FLOOR((nMaxUser * nUserRatable) / 100) THEN 1
            ELSE 0
        END AS nWorldStatus
    FROM tWorldServer
    ORDER BY nWorldNo ASC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 4. usp_WorldServer_GetByNo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_WorldServer_GetByNo`;
DELIMITER $$
CREATE PROCEDURE `usp_WorldServer_GetByNo`(
    IN in_worldNo TINYINT
)
BEGIN
    SELECT 
        nWorldNo,
        sWorldName,
        sWorldIP,
        nWorldPort,
        nMaxUser,
        nUserRatable,
        bIsMaintenance,
        nCurrentCCU,
        dLastHeartbeat,
        CASE
            WHEN bIsMaintenance = 1 THEN 2
            WHEN nCurrentCCU >= FLOOR((nMaxUser * nUserRatable) / 100) THEN 1
            ELSE 0
        END AS nWorldStatus
    FROM tWorldServer
    WHERE nWorldNo = in_worldNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 5. usp_WorldServer_UpdateCCU
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_WorldServer_UpdateCCU`;
DELIMITER $$
CREATE PROCEDURE `usp_WorldServer_UpdateCCU`(
    IN in_worldNo TINYINT,
    IN in_ccu INT
)
BEGIN
    UPDATE tWorldServer
    SET nCurrentCCU = in_ccu,
        dLastHeartbeat = NOW()
    WHERE nWorldNo = in_worldNo;

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 6. usp_Session_RegisterLogin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_RegisterLogin`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_RegisterLogin`(
    IN in_userNo INT,
    IN in_userID VARCHAR(60),
    IN in_ip VARCHAR(30)
)
BEGIN
    DECLARE v_isPreviousActive TINYINT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM tActiveAccountSession WHERE nUserNo = in_userNo AND nState = 1) THEN
        SET v_isPreviousActive = 1;
    END IF;

    INSERT INTO tActiveAccountSession (nUserNo, sUserID, sClientIP, nState, dLoginTime, dUpdatedAt)
    VALUES (in_userNo, in_userID, in_ip, 1, NOW(), NOW())
    ON DUPLICATE KEY UPDATE 
        sUserID = VALUES(sUserID),
        sClientIP = VALUES(sClientIP),
        nState = 1,
        dLoginTime = NOW(),
        dUpdatedAt = NOW();

    SELECT v_isPreviousActive AS bIsPreviousActive;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 7. usp_Session_IssueWorldToken
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_IssueWorldToken`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_IssueWorldToken`(
    IN in_userNo INT,
    IN in_worldNo TINYINT,
    IN in_token64 BINARY(64)
)
BEGIN
    UPDATE tActiveAccountSession
    SET nWorldNo = in_worldNo,
        abValidateToken = in_token64,
        dTokenExpire = NOW() + INTERVAL 60 SECOND,
        dUpdatedAt = NOW()
    WHERE nUserNo = in_userNo;

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 8. usp_Session_ConsumeWorldToken
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_ConsumeWorldToken`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_ConsumeWorldToken`(
    IN in_worldNo TINYINT,
    IN in_token64 BINARY(64),
    IN in_ip VARCHAR(30)
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_userID VARCHAR(60) DEFAULT '';
    DECLARE v_ret INT DEFAULT -1;

    SELECT nUserNo, sUserID INTO v_userNo, v_userID
    FROM tActiveAccountSession
    WHERE nWorldNo = in_worldNo
      AND abValidateToken = in_token64
      AND dTokenExpire > NOW()
    LIMIT 1;

    IF v_userNo > 0 THEN
        UPDATE tActiveAccountSession
        SET abValidateToken = NULL,
            dTokenExpire = NULL,
            dUpdatedAt = NOW()
        WHERE nUserNo = v_userNo;

        SET v_ret = 0;
    ELSE
        SET v_userNo = 0;
        SET v_userID = '';
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet, v_userNo AS nUserNo, v_userID AS sUserID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 9. usp_Session_CreateOTP
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_CreateOTP`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_CreateOTP`(
    IN in_userNo INT,
    IN in_otp32 CHAR(32)
)
BEGIN
    UPDATE tActiveAccountSession
    SET sOTP = in_otp32,
        dOTPExpire = NOW() + INTERVAL 300 SECOND,
        dUpdatedAt = NOW()
    WHERE nUserNo = in_userNo;

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 10. usp_Session_ConsumeOTP
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_ConsumeOTP`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_ConsumeOTP`(
    IN in_otp32 CHAR(32),
    IN in_ip VARCHAR(30)
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_userID VARCHAR(60) DEFAULT '';
    DECLARE v_ret INT DEFAULT -1;

    SELECT nUserNo, sUserID INTO v_userNo, v_userID
    FROM tActiveAccountSession
    WHERE sOTP = in_otp32
      AND dOTPExpire > NOW()
    LIMIT 1;

    IF v_userNo > 0 THEN
        UPDATE tActiveAccountSession
        SET sOTP = NULL,
            dOTPExpire = NULL,
            dUpdatedAt = NOW()
        WHERE nUserNo = v_userNo;

        SET v_ret = 0;
    ELSE
        SET v_userNo = 0;
        SET v_userID = '';
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet, v_userNo AS nUserNo, v_userID AS sUserID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 11. usp_Session_CloseAndGetPlayMin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Session_CloseAndGetPlayMin`;
DELIMITER $$
CREATE PROCEDURE `usp_Session_CloseAndGetPlayMin`(
    IN in_userNo INT
)
BEGIN
    DECLARE v_loginTime DATETIME;
    DECLARE v_playMin INT DEFAULT 0;

    SELECT dLoginTime INTO v_loginTime
    FROM tActiveAccountSession
    WHERE nUserNo = in_userNo
    LIMIT 1;

    IF v_loginTime IS NOT NULL THEN
        SET v_playMin = GREATEST(0, TIMESTAMPDIFF(MINUTE, v_loginTime, NOW()));
        DELETE FROM tActiveAccountSession WHERE nUserNo = in_userNo;
    ELSE
        SET v_playMin = 0;
    END IF;

    SELECT v_playMin AS nPlayMin;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 12. p_HackScan_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_HackScan_Add`;
DELIMITER $$
CREATE PROCEDURE `p_HackScan_Add`(
    IN in_nUserNo INT,
    IN in_nCharNo INT,
    IN in_nWorldNo TINYINT,
    IN in_nHackType SMALLINT,
    IN in_sDetail VARCHAR(64)
)
BEGIN
    INSERT INTO tHackScanLog (nUserNo, nCharNo, nWorldNo, nHackType, sDetail, dDate)
    VALUES (in_nUserNo, in_nCharNo, in_nWorldNo, in_nHackType, in_sDetail, NOW());

    SELECT LAST_INSERT_ID() AS nLogNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 13. p_Syslog_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Syslog_Add`;
DELIMITER $$
CREATE PROCEDURE `p_Syslog_Add`(
    IN in_nOpCode SMALLINT UNSIGNED,
    IN in_nUserNo INT,
    IN in_sUserID VARCHAR(60),
    IN in_sIP VARCHAR(30),
    IN in_nWorldNo TINYINT,
    IN in_nCCU INT,
    IN in_nErrorCode SMALLINT
)
BEGIN
    INSERT INTO tSyslog (nOpCode, nUserNo, sUserID, sIP, nWorldNo, nCCU, nErrorCode, dDate)
    VALUES (in_nOpCode, in_nUserNo, in_sUserID, in_sIP, in_nWorldNo, in_nCCU, in_nErrorCode, NOW());

    SELECT LAST_INSERT_ID() AS nLogNo;
END$$
DELIMITER ;
