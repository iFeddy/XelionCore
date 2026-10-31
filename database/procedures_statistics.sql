-- ============================================================================
-- PROCEDIMIENTOS ALMACENADOS OFICIALES DE STATISTICSDATA EN MYSQL (18 PROCEDIMIENTOS)
-- Base de Datos: xelion_online
-- Fidelidad Binaria y de Protocolo: 100% oficial respecto a StatisticsData.bak
-- ============================================================================

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- 1. pFindCCU
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `pFindCCU`;
DELIMITER $$
CREATE PROCEDURE `pFindCCU`(
    IN in_start DATETIME,
    IN in_end DATETIME
)
BEGIN
    SELECT dDate, w0, w1, (w0 + w1) AS total
    FROM (
        SELECT DATE_FORMAT(dDate, '%Y-%m-%d %H:%i') AS dDate,
               MAX(CASE nWorld WHEN 0 THEN nCount ELSE 0 END) AS w0,
               MAX(CASE nWorld WHEN 1 THEN nCount ELSE 0 END) AS w1
        FROM tSameTimeChar
        WHERE dDate BETWEEN in_start AND in_end
        GROUP BY DATE_FORMAT(dDate, '%Y-%m-%d %H:%i')
    ) AS t1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 2. p_DayClassLevelPlay_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayClassLevelPlay_Add`;
DELIMITER $$
CREATE PROCEDURE `p_DayClassLevelPlay_Add`(
    IN in_nWorld TINYINT,
    IN in_nClass TINYINT,
    IN in_nLevel SMALLINT,
    IN in_nPlayMin INT,
    IN in_nConnectCount INT,
    OUT out_nIndex INT
)
BEGIN
    INSERT INTO tDayClassLevelPlay (dDate, nWorld, nClass, nLevel, nPlayMin, nConnectCount)
    VALUES (NOW(), in_nWorld, in_nClass, in_nLevel, in_nPlayMin, in_nConnectCount);

    SET out_nIndex = LAST_INSERT_ID();
    SELECT out_nIndex AS nIndex;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 3. p_DayClassLevelPlay_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayClassLevelPlay_Query`;
DELIMITER $$
CREATE PROCEDURE `p_DayClassLevelPlay_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nIndex, dDate, nWorld, nClass, nLevel, nPlayMin, nConnectCount FROM tDayClassLevelPlay WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 4. p_DayClassLevel_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayClassLevel_Add`;
DELIMITER $$
CREATE PROCEDURE `p_DayClassLevel_Add`(
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
-- 5. p_DayClassLevel_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayClassLevel_Query`;
DELIMITER $$
CREATE PROCEDURE `p_DayClassLevel_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nIndex, dDate, nWorld, nClass, nLevel, nCount FROM tDayClassLevel WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 6. p_DayConnectionTime
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayConnectionTime`;
DELIMITER $$
CREATE PROCEDURE `p_DayConnectionTime`()
BEGIN
    DECLARE v_month INT;
    IF DAY(NOW()) = 1 THEN
        SET v_month = MONTH(NOW() - INTERVAL 1 MONTH);
    ELSE
        SET v_month = MONTH(NOW());
    END IF;

    SELECT MIN(nPlayMin) AS nMin, MAX(nPlayMin) AS nMax, AVG(nPlayMin) AS nAvg, nWorldNo, DAY(dDate) AS nDay
    FROM tAccountLog
    WHERE MONTH(dDate) = v_month AND nType = 2
    GROUP BY nWorldNo, DAY(dDate)
    ORDER BY nWorldNo, DAY(dDate);
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 7. p_DayUser
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayUser`;
DELIMITER $$
CREATE PROCEDURE `p_DayUser`()
BEGIN
    DECLARE v_yesterday DATE;
    DECLARE v_startDate DATETIME;
    DECLARE v_endDate DATETIME;
    DECLARE v_userCount INT;
    DECLARE v_connectCount INT;
    DECLARE v_index INT;

    SET v_yesterday = DATE_SUB(CURDATE(), INTERVAL 1 DAY);
    SET v_startDate = CAST(v_yesterday AS DATETIME);
    SET v_endDate = CAST(CURDATE() AS DATETIME);

    DELETE FROM tDayUser WHERE dDate = v_startDate;

    SELECT COUNT(DISTINCT nUserNo), COUNT(*)
    INTO v_userCount, v_connectCount
    FROM tAccountLog
    WHERE nType = 1 AND dDate >= v_startDate AND dDate < v_endDate;

    INSERT INTO tDayUser (dDate, nUserCount, nConnectCount)
    VALUES (v_startDate, v_userCount, v_connectCount);
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 8. p_DayUser_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayUser_Add`;
DELIMITER $$
CREATE PROCEDURE `p_DayUser_Add`(
    IN in_nUserCount INT,
    IN in_nConnectCount INT,
    OUT out_nIndex INT
)
BEGIN
    INSERT INTO tDayUser (dDate, nUserCount, nConnectCount)
    VALUES (NOW(), in_nUserCount, in_nConnectCount);

    SET out_nIndex = LAST_INSERT_ID();
    SELECT out_nIndex AS nIndex;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 9. p_DayUser_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_DayUser_Query`;
DELIMITER $$
CREATE PROCEDURE `p_DayUser_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nIndex, dDate, nUserCount, nConnectCount FROM tDayUser WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 10. p_SameTimeChar_Add
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_SameTimeChar_Add`;
DELIMITER $$
CREATE PROCEDURE `p_SameTimeChar_Add`(
    IN in_nWorld TINYINT,
    IN in_nCount INT,
    OUT out_nIndex INT
)
BEGIN
    INSERT INTO tSameTimeChar (dDate, nWorld, nCount)
    VALUES (NOW(), in_nWorld, in_nCount);

    SET out_nIndex = LAST_INSERT_ID();
    SELECT out_nIndex AS nIndex;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 11. p_SameTimeChar_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_SameTimeChar_Query`;
DELIMITER $$
CREATE PROCEDURE `p_SameTimeChar_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nIndex, dDate, nWorld, nCount FROM tSameTimeChar WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 12. p_SameTimeWorld
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_SameTimeWorld`;
DELIMITER $$
CREATE PROCEDURE `p_SameTimeWorld`()
BEGIN
    DECLARE v_latestTime DATETIME;

    SELECT MAX(dDate) INTO v_latestTime
    FROM tSameTimeChar
    WHERE dDate >= NOW() - INTERVAL 5 MINUTE;

    IF v_latestTime IS NULL THEN
        SELECT NOW() AS CurTime, 0 AS CurCount;
    ELSE
        SELECT v_latestTime AS CurTime, IFNULL(SUM(nCount), 0) AS CurCount
        FROM tSameTimeChar
        WHERE dDate = v_latestTime;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 13. p_UserAgeReport
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserAgeReport`;
DELIMITER $$
CREATE PROCEDURE `p_UserAgeReport`(
    IN in_S_DATE VARCHAR(10),
    IN in_E_DATE VARCHAR(10)
)
BEGIN
    DECLARE v_sDate DATE;
    DECLARE v_eDate DATE;

    IF in_S_DATE IS NULL OR in_S_DATE = '' THEN
        SET v_sDate = DATE_SUB(CURDATE(), INTERVAL 1 DAY);
    ELSE
        SET v_sDate = STR_TO_DATE(in_S_DATE, '%Y-%m-%d');
    END IF;

    IF in_E_DATE IS NULL OR in_E_DATE = '' THEN
        SET v_eDate = DATE_SUB(CURDATE(), INTERVAL 1 DAY);
    ELSE
        SET v_eDate = STR_TO_DATE(in_E_DATE, '%Y-%m-%d');
    END IF;

    SELECT CASE nGender
               WHEN 0 THEN '여자'
               WHEN 1 THEN '남자'
               ELSE '기타'
           END AS `성별`,
           nCount AS `인원`
    FROM tGender
    WHERE STR_TO_DATE(CONCAT(nYear, '-', LPAD(nMonth, 2, '0'), '-', LPAD(nDay, 2, '0')), '%Y-%m-%d') BETWEEN v_sDate AND v_eDate;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 14. p_UserAge_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserAge_Query`;
DELIMITER $$
CREATE PROCEDURE `p_UserAge_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, nAge, nCount FROM tAge WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 15. p_UserGenderTest
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserGenderTest`;
DELIMITER $$
CREATE PROCEDURE `p_UserGenderTest`(
    IN in_DATE VARCHAR(10)
)
BEGIN
    DECLARE v_logDate VARCHAR(10);

    IF in_DATE IS NULL OR in_DATE = '' THEN
        SET v_logDate = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y-%m-%d');
    ELSE
        SET v_logDate = in_DATE;
    END IF;

    -- Agrega datos de genero desde tAccountDemographics para la fecha
    INSERT INTO tGender (sLogDate, nGender, nCount, nYear, nMonth, nDay, nWeek, nWeekDay)
    SELECT v_logDate,
           gender,
           COUNT(*),
           YEAR(STR_TO_DATE(v_logDate, '%Y-%m-%d')),
           MONTH(STR_TO_DATE(v_logDate, '%Y-%m-%d')),
           DAY(STR_TO_DATE(v_logDate, '%Y-%m-%d')),
           WEEK(STR_TO_DATE(v_logDate, '%Y-%m-%d')),
           DAYOFWEEK(STR_TO_DATE(v_logDate, '%Y-%m-%d'))
    FROM tAccountDemographics
    WHERE DATE(registerDate) = STR_TO_DATE(v_logDate, '%Y-%m-%d')
    GROUP BY gender;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 16. p_UserGender_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserGender_Query`;
DELIMITER $$
CREATE PROCEDURE `p_UserGender_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, nGender, nCount FROM tGender WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 17. p_UserLocation_Query
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserLocation_Query`;
DELIMITER $$
CREATE PROCEDURE `p_UserLocation_Query`(
    IN in_WhereExpr VARCHAR(256)
)
BEGIN
    IF in_WhereExpr IS NULL OR in_WhereExpr = '' THEN
        SET in_WhereExpr = '1=1';
    END IF;

    SET @sql = CONCAT('SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, sLocation, nCount FROM tLocation WHERE ', in_WhereExpr);
    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 18. p_getQuestTimes
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_getQuestTimes`;
DELIMITER $$
CREATE PROCEDURE `p_getQuestTimes`(
    IN in_nQuestNo INT,
    IN in_top INT
)
BEGIN
    IF in_top IS NULL OR in_top <= 0 THEN
        SET in_top = 10;
    END IF;

    IF in_nQuestNo = 20016 THEN
        SELECT ServerNo, sID, nTimes, dLastComplete FROM tQuestTimesQ20016 ORDER BY nTimes DESC LIMIT in_top;
    ELSEIF in_nQuestNo = 2079 THEN
        SELECT ServerNo, sID, nTimes, dLastComplete FROM tQuestTimesQ2079 ORDER BY nTimes DESC LIMIT in_top;
    ELSEIF in_nQuestNo = 63011 THEN
        SELECT ServerNo, sID, nTimes, dLastComplete FROM tQuestTimesQ63011 ORDER BY nTimes DESC LIMIT in_top;
    ELSE
        SELECT 0 AS ServerNo, '' AS sID, 0 AS nTimes, '' AS dLastComplete LIMIT 0;
    END IF;
END$$
DELIMITER ;
