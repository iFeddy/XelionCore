-- ============================================================================
-- PROCEDIMIENTOS ALMACENADOS OFICIALES DE ACCOUNTLOG EN MYSQL (28 PROCEDIMIENTOS)
-- Base de Datos: xelion_online
-- Fidelidad Binaria y de Protocolo: 100% oficial respecto a AccountLog.bak y AccountLog.pdb (11 ODBC + 17 no-OPT)
-- ============================================================================

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- 1. p_AccountLog_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_AccountLog_Add`;
DELIMITER $$
CREATE PROCEDURE `p_AccountLog_Add`(
    IN in_nType TINYINT,
    IN in_nUserNo INT,
    IN in_nWorldNo TINYINT,
    IN in_nPlayMin SMALLINT,
    IN in_nErrorCode SMALLINT,
    IN in_sIP VARCHAR(16),
    IN in_sUserID VARCHAR(60),
    IN in_sUserPassword VARCHAR(32),
    OUT out_nLogNo BIGINT
)
BEGIN
    SET out_nLogNo = 0;

    -- Si es Login Exitoso (nType = 1): actualizar unicamente sUserIP en tUser (sin tocar dDate)
    -- y registrar primer ingreso diario en tUserLog si no existe
    IF in_nType = 1 THEN
        UPDATE tUser SET sUserIP = in_sIP WHERE nUserNo = in_nUserNo;

        IF NOT EXISTS(SELECT 1 FROM tUserLog WHERE userNo = in_nUserNo AND joinDate = DATE_FORMAT(NOW(), '%Y-%m-%d')) THEN
            INSERT INTO tUserLog (joinDate, userNo) VALUES (DATE_FORMAT(NOW(), '%Y-%m-%d'), in_nUserNo);
        END IF;
    END IF;

    INSERT INTO tAccountLog (nType, nUserNo, nWorldNo, nPlayMin, nErrorCode, sIP, sUserID, sUserPassword, dDate)
    VALUES (in_nType, in_nUserNo, in_nWorldNo, in_nPlayMin, in_nErrorCode, in_sIP, in_sUserID, in_sUserPassword, NOW());

    SET out_nLogNo = LAST_INSERT_ID();
    SELECT out_nLogNo AS nLogNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 2. p_IPBlock_IsBlock
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_IPBlock_IsBlock`;
DELIMITER $$
CREATE PROCEDURE `p_IPBlock_IsBlock`(
    IN in_nUserNo INT,
    IN in_sIP VARCHAR(20),
    OUT out_bIsBlock TINYINT
)
BEGIN
    DECLARE v_ipNum BIGINT UNSIGNED;
    SET out_bIsBlock = 0;
    SET v_ipNum = INET_ATON(in_sIP);

    IF EXISTS(
        SELECT 1 FROM tIPBlock
        WHERE (nIPStartNum IS NOT NULL AND nIPEndNum IS NOT NULL AND v_ipNum BETWEEN nIPStartNum AND nIPEndNum)
           OR (v_ipNum BETWEEN INET_ATON(sIPStart) AND INET_ATON(sIPEnd))
    ) THEN
        -- Whitelist de cuentas veteranas: si ya tiene al menos un login exitoso (nType = 1), permitir acceso
        IF NOT EXISTS(SELECT 1 FROM tAccountLog WHERE nUserNo = in_nUserNo AND nType = 1 LIMIT 1) THEN
            SET out_bIsBlock = 1;
        END IF;
    END IF;

    SELECT out_bIsBlock AS bIsBlock;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 3. usp_Log_setLogout
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Log_setLogout`;
DELIMITER $$
CREATE PROCEDURE `usp_Log_setLogout`(
    IN in_WORLD_NO TINYINT,
    IN in_USER_NO INT,
    IN in_CHAR_NO INT,
    IN in_CLASS TINYINT,
    IN in_LEVEL SMALLINT,
    IN in_EXPERIENCE BIGINT,
    IN in_PLAYTIME SMALLINT,
    OUT out_RET INT
)
BEGIN
    INSERT INTO tLogout (logDate, worldNo, userNo, charNo, class, level, experience, playTime)
    VALUES (NOW(), in_WORLD_NO, in_USER_NO, in_CHAR_NO, in_CLASS, in_LEVEL, in_EXPERIENCE, in_PLAYTIME);

    SET out_RET = 1;
    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 4. p_Report_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Report_Add`;
DELIMITER $$
CREATE PROCEDURE `p_Report_Add`(
    IN in_nReporter INT,
    IN in_sReporterCharID VARCHAR(16),
    IN in_sHarmer VARCHAR(16),
    IN in_nWorldNo TINYINT,
    IN in_sReportType VARCHAR(16),
    IN in_sReason VARCHAR(512),
    IN in_sChatLog TEXT,
    OUT out_nReportRegNo INT
)
BEGIN
    SET out_nReportRegNo = 0;

    -- Cooldown 1h por reportero en el mismo mundo
    IF EXISTS(SELECT 1 FROM tReport WHERE nReporter = in_nReporter AND nWorldNo = in_nWorldNo AND dReportDate > NOW() - INTERVAL 1 HOUR) THEN
        SET out_nReportRegNo = -1;
    -- Cooldown 24h para denunciar al mismo personaje
    ELSEIF EXISTS(SELECT 1 FROM tReport WHERE nReporter = in_nReporter AND nWorldNo = in_nWorldNo AND sHarmer = in_sHarmer AND dReportDate > NOW() - INTERVAL 1 DAY) THEN
        SET out_nReportRegNo = -2;
    ELSE
        INSERT INTO tReport (sReportType, nReporter, sReporterCharID, sHarmer, nWorldNo, sReason, sChatLog, dReportDate)
        VALUES (in_sReportType, in_nReporter, in_sReporterCharID, in_sHarmer, in_nWorldNo, in_sReason, in_sChatLog, NOW());

        SET out_nReportRegNo = LAST_INSERT_ID();
    END IF;

    SELECT out_nReportRegNo AS nReportRegNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 5. p_Report_Del
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Report_Del`;
DELIMITER $$
CREATE PROCEDURE `p_Report_Del`(
    IN in_nRegNumber INT,
    OUT out_nRet SMALLINT
)
BEGIN
    DECLARE v_reportDate DATETIME;
    DECLARE v_result TINYINT;

    SET out_nRet = 0;

    SELECT dReportDate, nResult INTO v_reportDate, v_result
    FROM tReport
    WHERE nRegNumber = in_nRegNumber
    LIMIT 1;

    IF v_reportDate IS NULL THEN
        SET out_nRet = 1; -- Inexistente
    ELSEIF v_reportDate < NOW() - INTERVAL 1 HOUR THEN
        SET out_nRet = 2; -- Excedio 1 hora
    ELSEIF v_result <> 0 THEN
        SET out_nRet = 3; -- Ya resuelto
    ELSE
        DELETE FROM tReport WHERE nRegNumber = in_nRegNumber;
        SET out_nRet = 0;
    END IF;

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 6. p_Report_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Report_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Report_Get`(
    IN in_nReporter INT,
    IN in_nWorldNo TINYINT
)
BEGIN
    SELECT nRegNumber, nResult, sRemark
    FROM tReport
    WHERE nReporter = in_nReporter
      AND nWorldNo = in_nWorldNo
      AND dReportDate > NOW() - INTERVAL 1 HOUR;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 7. p_SameTimeCharClass_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_SameTimeCharClass_Add`;
DELIMITER $$
CREATE PROCEDURE `p_SameTimeCharClass_Add`(
    IN in_nWorld TINYINT,
    IN in_nClass TINYINT,
    IN in_nLevel SMALLINT,
    IN in_nCount INT,
    OUT out_nIndex INT
)
BEGIN
    INSERT INTO tDayClassLevel (dDate, nWorld, nClass, nLevel, nCount)
    VALUES (NOW(), in_nWorld, in_nClass, in_nLevel, in_nCount);

    SET out_nIndex = LAST_INSERT_ID();
    SELECT out_nIndex AS nIndex;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 9. p_Event_L20_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_L20_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Event_L20_Get`(
    IN in_nUserNo INT,
    OUT out_bRewarded TINYINT,
    OUT out_nRet INT
)
BEGIN
    SET out_bRewarded = 0;
    SET out_nRet = -1;

    IF EXISTS(SELECT 1 FROM tEvent_L20 WHERE nUserNo = in_nUserNo) THEN
        SELECT bRewarded INTO out_bRewarded
        FROM tEvent_L20
        WHERE nUserNo = in_nUserNo
        LIMIT 1;
        SET out_nRet = 0;
    ELSE
        SET out_bRewarded = 0;
        SET out_nRet = -1;
    END IF;

    SELECT out_nRet AS nRet, out_bRewarded AS bRewarded;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 10. p_Event_L20_Set
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_L20_Set`;
DELIMITER $$
CREATE PROCEDURE `p_Event_L20_Set`(
    IN in_nUserNo INT,
    IN in_sUserID VARCHAR(60),
    IN in_nWorldNo TINYINT,
    IN in_nCharNo INT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRet = 0;

    IF EXISTS(SELECT 1 FROM tEvent_L20 WHERE nUserNo = in_nUserNo) THEN
        SET out_nRet = -1;
    ELSE
        INSERT INTO tEvent_L20 (nUserNo, sUserID, nWorldNo, nCharNo, bRewarded, dDate)
        VALUES (in_nUserNo, in_sUserID, in_nWorldNo, in_nCharNo, 1, NOW());

        SET out_nRet = 0;
    END IF;

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 11. usp_AccountLog_GetDateDiffLastlogout
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_AccountLog_GetDateDiffLastlogout`;
DELIMITER $$
CREATE PROCEDURE `usp_AccountLog_GetDateDiffLastlogout`(
    IN in_nUserNo INT,
    OUT out_nDiffSec INT
)
BEGIN
    DECLARE v_lastDate DATETIME;
    SET out_nDiffSec = -1;

    SELECT MAX(logDate) INTO v_lastDate
    FROM tLogout
    WHERE userNo = in_nUserNo;

    IF v_lastDate IS NULL THEN
        SELECT MAX(dDate) INTO v_lastDate
        FROM tAccountLog
        WHERE nUserNo = in_nUserNo AND nType IN (1, 2);
    END IF;

    IF v_lastDate IS NOT NULL THEN
        SET out_nDiffSec = TIMESTAMPDIFF(SECOND, v_lastDate, NOW());
    ELSE
        SET out_nDiffSec = -1;
    END IF;

    SELECT out_nDiffSec AS nDiffSec;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 12. usp_EventiPad_setTicket
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_EventiPad_setTicket`;
DELIMITER $$
CREATE PROCEDURE `usp_EventiPad_setTicket`(
    IN in_nUserNo INT,
    IN in_nPlayTime INT,
    OUT out_nRet INT
)
BEGIN
    DECLARE v_tickets INT DEFAULT 0;
    SET out_nRet = 0;

    IF in_nPlayTime IS NULL OR in_nPlayTime < 60 THEN
        SET out_nRet = -1;
    ELSE
        SET v_tickets = FLOOR(in_nPlayTime / 60);

        INSERT INTO tEventiPad (nUserNo, nTicket, sys_create_date, sys_update_date)
        VALUES (in_nUserNo, v_tickets, NOW(), NOW())
        ON DUPLICATE KEY UPDATE nTicket = nTicket + VALUES(nTicket), sys_update_date = NOW();

        SET out_nRet = 0;
    END IF;

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 13. p_AccountLog_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_AccountLog_Get`;
DELIMITER $$
CREATE PROCEDURE `p_AccountLog_Get`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nLogNo, dDate, nType, nUserNo, nWorldNo, nPlayMin, nErrorCode, sIP, sUserID, sUserPassword FROM tAccountLog WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 14. p_AccountLog_GetCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_AccountLog_GetCount`;
DELIMITER $$
CREATE PROCEDURE `p_AccountLog_GetCount`(
    IN in_WhereExpr VARCHAR(256),
    OUT out_Count INT
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT COUNT(*) INTO @cnt FROM tAccountLog WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SET out_Count = IFNULL(@cnt, 0);
    SELECT out_Count AS `Count`;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 15. p_AccountLog_GetLog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_AccountLog_GetLog`;
DELIMITER $$
CREATE PROCEDURE `p_AccountLog_GetLog`(
    IN in_QueryExpr VARCHAR(256)
)
BEGIN
    IF in_QueryExpr IS NULL OR in_QueryExpr = '' THEN
        SET in_QueryExpr = '';
    END IF;

    SET @sql = CONCAT('SELECT a.nLogNo, a.dDate, a.nUserNo, b.sUserID, a.nType, a.sIP, a.nPlayMin, a.sUserPassword, a.nErrorCode, a.nWorldNo FROM tAccountLog a INNER JOIN tUser b ON a.nUserNo = b.nUserNo ', in_QueryExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 16. p_AccountLog_GetLogCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_AccountLog_GetLogCount`;
DELIMITER $$
CREATE PROCEDURE `p_AccountLog_GetLogCount`(
    IN in_QueryExpr VARCHAR(256),
    OUT out_Count INT
)
BEGIN
    IF in_QueryExpr IS NULL OR in_QueryExpr = '' THEN
        SET in_QueryExpr = '';
    END IF;

    SET @sql = CONCAT('SELECT COUNT(*) INTO @cnt FROM tAccountLog a INNER JOIN tUser b ON a.nUserNo = b.nUserNo ', in_QueryExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    SET out_Count = IFNULL(@cnt, 0);
    SELECT out_Count AS `Count`;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 17. p_IPBlock_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_IPBlock_Add`;
DELIMITER $$
CREATE PROCEDURE `p_IPBlock_Add`(
    IN in_sIPStart VARCHAR(20),
    IN in_sIPEnd VARCHAR(20),
    IN in_sDesc VARCHAR(30),
    OUT out_nRet INT
)
BEGIN
    INSERT INTO tIPBlock (sIPStart, sIPEnd, sDesc, dDate)
    VALUES (in_sIPStart, in_sIPEnd, in_sDesc, NOW());

    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 18. p_IPBlock_Delete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_IPBlock_Delete`;
DELIMITER $$
CREATE PROCEDURE `p_IPBlock_Delete`(
    IN in_sIP VARCHAR(20),
    OUT out_nRet INT
)
BEGIN
    DECLARE v_ipNum BIGINT UNSIGNED;
    SET v_ipNum = INET_ATON(in_sIP);

    DELETE FROM tIPBlock
    WHERE (nIPStartNum IS NOT NULL AND nIPEndNum IS NOT NULL AND v_ipNum BETWEEN nIPStartNum AND nIPEndNum)
       OR (v_ipNum BETWEEN INET_ATON(sIPStart) AND INET_ATON(sIPEnd));

    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 19. p_IPBlock_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_IPBlock_Get`;
DELIMITER $$
CREATE PROCEDURE `p_IPBlock_Get`(
    IN in_sIP VARCHAR(20)
)
BEGIN
    DECLARE v_ipNum BIGINT UNSIGNED;
    SET v_ipNum = INET_ATON(in_sIP);

    SELECT sIPStart, sIPEnd, sDesc, dDate
    FROM tIPBlock
    WHERE (nIPStartNum IS NOT NULL AND nIPEndNum IS NOT NULL AND v_ipNum BETWEEN nIPStartNum AND nIPEndNum)
       OR (v_ipNum BETWEEN INET_ATON(sIPStart) AND INET_ATON(sIPEnd));
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 20. p_IPBlock_GetList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_IPBlock_GetList`;
DELIMITER $$
CREATE PROCEDURE `p_IPBlock_GetList`()
BEGIN
    SELECT sIPStart, sIPEnd, sDesc, dDate
    FROM tIPBlock;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 21. p_Event_Friend_DelRequest
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_DelRequest`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_DelRequest`(
    IN in_nEventUserNo INT,
    OUT out_err SMALLINT
)
BEGIN
    DELETE FROM tEvent_Friend WHERE nEventUserNo = in_nEventUserNo;

    IF ROW_COUNT() = 0 THEN
        SET out_err = 1;
    ELSE
        SET out_err = 0;
    END IF;

    SELECT out_err AS err;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 22. p_Event_Friend_EventComplete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_EventComplete`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_EventComplete`(
    IN in_nEventUserNo INT,
    OUT out_nFriendCharNo INT,
    OUT out_err SMALLINT
)
BEGIN
    SET out_nFriendCharNo = 0;

    SELECT nFriendCharNo INTO out_nFriendCharNo
    FROM tEvent_Friend
    WHERE nEventUserNo = in_nEventUserNo
    LIMIT 1;

    IF out_nFriendCharNo = 0 OR out_nFriendCharNo IS NULL THEN
        SET out_err = 1;
    ELSE
        UPDATE tEvent_Friend
        SET bEventComplete = 1, dCompleteDate = NOW()
        WHERE nEventUserNo = in_nEventUserNo;

        SET out_err = 0;
    END IF;

    SELECT out_err AS err, out_nFriendCharNo AS nFriendCharNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 23. p_Event_Friend_Get
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_Get`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_Get`(
    IN in_nFriendUserNo INT
)
BEGIN
    SELECT nEventUserNo, nEventCharNo, nFriendCharNo, nWorldNo, bEventComplete, dDate
    FROM tEvent_Friend
    WHERE nFriendUserNo = in_nFriendUserNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 24. p_Event_Friend_GetCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_GetCount`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_GetCount`(
    IN in_nUserNo INT,
    OUT out_nCount INT
)
BEGIN
    SELECT COUNT(*) INTO out_nCount
    FROM tEvent_Friend
    WHERE nFriendUserNo = in_nUserNo;

    SELECT out_nCount AS nCount;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 25. p_Event_Friend_Giveup
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_Giveup`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_Giveup`(
    IN in_nEventUserNo INT,
    OUT out_nFriendCharNo INT,
    OUT out_err SMALLINT
)
BEGIN
    SET out_nFriendCharNo = 0;

    SELECT nFriendCharNo INTO out_nFriendCharNo
    FROM tEvent_Friend
    WHERE nEventUserNo = in_nEventUserNo
    LIMIT 1;

    IF out_nFriendCharNo = 0 OR out_nFriendCharNo IS NULL THEN
        SET out_err = 1;
    ELSE
        DELETE FROM tEvent_Friend WHERE nEventUserNo = in_nEventUserNo;
        SET out_err = 0;
    END IF;

    SELECT out_err AS err, out_nFriendCharNo AS nFriendCharNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 26. p_Event_Friend_SetHideMsg
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_Friend_SetHideMsg`;
DELIMITER $$
CREATE PROCEDURE `p_Event_Friend_SetHideMsg`(
    IN in_nUserNo INT,
    OUT out_err SMALLINT
)
BEGIN
    IF EXISTS(SELECT 1 FROM tEvent_Friend_HideMsg WHERE nUserNo = in_nUserNo) THEN
        SET out_err = 1;
    ELSE
        INSERT INTO tEvent_Friend_HideMsg (nUserNo) VALUES (in_nUserNo);
        SET out_err = 0;
    END IF;

    SELECT out_err AS err;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 27. p_Event_WorldCup_Enter
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_Event_WorldCup_Enter`;
DELIMITER $$
CREATE PROCEDURE `p_Event_WorldCup_Enter`(
    IN in_nCharNo INT,
    IN in_nUserNo INT,
    IN in_nLevel TINYINT,
    IN in_nWorldNo TINYINT,
    OUT out_err SMALLINT
)
BEGIN
    IF EXISTS(SELECT 1 FROM tEventWorldCup WHERE nEventCharNo = in_nCharNo AND nLevel = in_nLevel AND nWorldNo = in_nWorldNo) THEN
        SET out_err = 1;
    ELSE
        INSERT INTO tEventWorldCup (nEventCharNo, nEventUserNo, nLevel, nWorldNo)
        VALUES (in_nCharNo, in_nUserNo, in_nLevel, in_nWorldNo);
        SET out_err = 0;
    END IF;

    SELECT out_err AS err;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 28. usp_EventiPad_getList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_EventiPad_getList`;
DELIMITER $$
CREATE PROCEDURE `usp_EventiPad_getList`()
BEGIN
    SELECT nUserNo, nTicket, sys_create_date AS dCreateDate, sys_update_date AS dLastUpdate
    FROM tEventiPad
    ORDER BY nTicket DESC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 29. usp_EventiPad_getTicket
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_EventiPad_getTicket`;
DELIMITER $$
CREATE PROCEDURE `usp_EventiPad_getTicket`(
    IN in_nUserNo INT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRet = 0;

    IF in_nUserNo IS NULL OR in_nUserNo <= 0 THEN
        SET out_nRet = -1;
    ELSE
        SELECT nUserNo, nTicket, sys_create_date AS dCreateDate, sys_update_date AS dLastUpdate
        FROM tEventiPad
        WHERE nUserNo = in_nUserNo;
        SET out_nRet = 0;
    END IF;
END$$
DELIMITER ;
