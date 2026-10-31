-- ============================================================================
-- PROCEDIMIENTOS ALMACENADOS OFICIALES DE ACCOUNT EN MYSQL (75 PROCEDIMIENTOS)
-- Base de Datos: xelion_online
-- Fidelidad Binaria y de Protocolo: 100% oficial respecto a Account.bak y Account.pdb
-- ============================================================================

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- 1. p_UserGetAuthList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserGetAuthList`;
DELIMITER $$
CREATE PROCEDURE `p_UserGetAuthList`()
BEGIN
    SELECT nAuthID, sAuthName, bIsLoginAble FROM tUserAuth;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 2. p_UserKey_Get_FromAccount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserKey_Get_FromAccount`;
DELIMITER $$
CREATE PROCEDURE `p_UserKey_Get_FromAccount`(
    IN in_nUserNo INT,
    OUT out_sWebKey VARCHAR(64),
    OUT out_dExpire DATETIME,
    OUT out_nRet INT
)
BEGIN
    SET out_sWebKey = NULL;
    SET out_dExpire = NULL;
    SET out_nRet = -1;

    SELECT sWebKey, dExpireDate
    INTO out_sWebKey, out_dExpire
    FROM tUserKey
    WHERE nUserNo = in_nUserNo AND dExpireDate > NOW()
    ORDER BY dExpireDate DESC
    LIMIT 1;

    IF out_sWebKey IS NOT NULL THEN
        SET out_nRet = 0;
    ELSE
        SET out_sWebKey = '';
        SET out_dExpire = NULL;
        SET out_nRet = -1;
    END IF;

    SELECT out_nRet AS nRet, out_sWebKey AS sWebKey, out_dExpire AS dExpire;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 3. p_UserPlayRemainTime_Get_FromAccount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserPlayRemainTime_Get_FromAccount`;
DELIMITER $$
CREATE PROCEDURE `p_UserPlayRemainTime_Get_FromAccount`(
    IN in_nUserNo INT,
    OUT out_nRemainSec INT,
    OUT out_bIsTeenager TINYINT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRemainSec = NULL;
    SET out_bIsTeenager = 0;
    SET out_nRet = -1;

    SELECT nRemainSec, bIsTeenager
    INTO out_nRemainSec, out_bIsTeenager
    FROM tUserPlayRemainTime
    WHERE nUserNo = in_nUserNo
    LIMIT 1;

    IF out_nRemainSec IS NOT NULL THEN
        SET out_nRet = 0;
    ELSE
        SET out_nRemainSec = 0;
        SET out_bIsTeenager = 0;
        SET out_nRet = -1;
    END IF;

    SELECT out_nRet AS nRet, out_nRemainSec AS nRemainSec, out_bIsTeenager AS bIsTeenager;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 4. p_UserPlayRemainTime_Set_FromAccount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserPlayRemainTime_Set_FromAccount`;
DELIMITER $$
CREATE PROCEDURE `p_UserPlayRemainTime_Set_FromAccount`(
    IN in_nUserNo INT,
    IN in_nUsedSec INT,
    OUT out_nRemainSec INT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRemainSec = 0;
    SET out_nRet = -1;

    UPDATE tUserPlayRemainTime
    SET nRemainSec = GREATEST(0, nRemainSec - in_nUsedSec),
        dLastUpdate = NOW()
    WHERE nUserNo = in_nUserNo;

    IF ROW_COUNT() > 0 THEN
        SELECT nRemainSec INTO out_nRemainSec FROM tUserPlayRemainTime WHERE nUserNo = in_nUserNo LIMIT 1;
        SET out_nRet = 0;
    ELSE
        SET out_nRemainSec = 0;
        SET out_nRet = -1;
    END IF;

    SELECT out_nRet AS nRet, out_nRemainSec AS nRemainSec;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 5. p_UserSetAuthLoginable
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `p_UserSetAuthLoginable`;
DELIMITER $$
CREATE PROCEDURE `p_UserSetAuthLoginable`(
    IN in_userNo INT,
    IN in_AuthID TINYINT
)
BEGIN
    DECLARE v_ret INT DEFAULT -1;

    IF EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        UPDATE tUser SET nAuthID = in_AuthID WHERE nUserNo = in_userNo;
        IF ROW_COUNT() = 1 THEN
            SET v_ret = 0;
        ELSE
            SET v_ret = -9;
        END IF;
    ELSE
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 6. usp_Cash_getCash
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Cash_getCash`;
DELIMITER $$
CREATE PROCEDURE `usp_Cash_getCash`(
    IN in_userNo INT,
    IN in_cashType TINYINT,
    OUT out_CASH INT
)
BEGIN
    SET out_CASH = 0;
    SELECT IFNULL(cash, 0) INTO out_CASH
    FROM tCash
    WHERE userNo = in_userNo AND cashType = in_cashType AND status = 0
    LIMIT 1;

    IF out_CASH IS NULL THEN
        SET out_CASH = 0;
    END IF;

    SELECT out_CASH AS CASH;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 7. usp_Cash_setCash
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Cash_setCash`;
DELIMITER $$
CREATE PROCEDURE `usp_Cash_setCash`(
    IN in_userNo INT,
    IN in_cashType TINYINT,
    IN in_cash INT,
    IN in_reference VARCHAR(20)
)
BEGIN
    START TRANSACTION;

    INSERT INTO tCash (userNo, cashType, cash, status)
    VALUES (in_userNo, in_cashType, in_cash, 0)
    ON DUPLICATE KEY UPDATE cash = cash + VALUES(cash), status = 0;

    INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
    VALUES (in_userNo, 1, in_cashType, in_cash, in_reference);

    COMMIT;
    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 8. usp_Cash_setCash_EVENT_L20
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Cash_setCash_EVENT_L20`;
DELIMITER $$
CREATE PROCEDURE `usp_Cash_setCash_EVENT_L20`(
    IN in_userNo INT,
    IN in_cashType TINYINT,
    IN in_cash INT,
    IN in_reference VARCHAR(20)
)
BEGIN
    IF EXISTS(SELECT 1 FROM tCashLog WHERE userNo = in_userNo AND cashType = in_cashType AND reference = '200708_EVENT_L20') THEN
        SELECT -1 AS nRet;
    ELSE
        START TRANSACTION;

        INSERT INTO tCash (userNo, cashType, cash, status)
        VALUES (in_userNo, in_cashType, in_cash, 0)
        ON DUPLICATE KEY UPDATE cash = cash + VALUES(cash), status = 0;

        INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
        VALUES (in_userNo, 1, in_cashType, in_cash, in_reference);

        COMMIT;
        SELECT 0 AS nRet;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 9. usp_Cash_useCash
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Cash_useCash`;
DELIMITER $$
CREATE PROCEDURE `usp_Cash_useCash`(
    IN in_userNo INT,
    IN in_cashType TINYINT,
    IN in_cash INT,
    IN in_reference VARCHAR(20)
)
BEGIN
    DECLARE v_remainCash INT DEFAULT 0;

    START TRANSACTION;

    SELECT IFNULL(cash, 0) INTO v_remainCash
    FROM tCash
    WHERE userNo = in_userNo AND cashType = in_cashType AND status = 0
    FOR UPDATE;

    IF v_remainCash >= in_cash THEN
        UPDATE tCash
        SET cash = v_remainCash - in_cash
        WHERE userNo = in_userNo AND cashType = in_cashType AND status = 0;

        INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
        VALUES (in_userNo, 2, in_cashType, in_cash, in_reference);

        COMMIT;
        SELECT 0 AS nRet;
    ELSE
        ROLLBACK;
        SELECT -1 AS nRet;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 10. usp_Charge_ItemDel
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Charge_ItemDel`;
DELIMITER $$
CREATE PROCEDURE `usp_Charge_ItemDel`(
    IN in_rowNo INT,
    IN in_userNo INT,
    OUT out_RET INT
)
BEGIN
    DECLARE v_isDraw TINYINT DEFAULT NULL;

    SET out_RET = 0;

    SELECT isDraw INTO v_isDraw
    FROM tChargeItem
    WHERE rowNo = in_rowNo AND userNo = in_userNo
    LIMIT 1;

    IF v_isDraw IS NULL THEN
        SET out_RET = -2; -- mismatch rowNo / userNo
    ELSEIF v_isDraw <> 0 THEN
        SET out_RET = -1; -- already used / drawn
    ELSE
        START TRANSACTION;

        INSERT INTO tChargeItemDelLog (rowNo, userNo, orderNo, goodsNo, amount, registerDate, deleteDate)
        SELECT rowNo, userNo, orderNo, goodsNo, amount, registerDate, NOW()
        FROM tChargeItem
        WHERE rowNo = in_rowNo AND userNo = in_userNo;

        DELETE FROM tChargeItem WHERE rowNo = in_rowNo AND userNo = in_userNo;

        IF ROW_COUNT() = 1 THEN
            SET out_RET = 1;
            COMMIT;
        ELSE
            SET out_RET = 0;
            ROLLBACK;
        END IF;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 11. usp_Charge_ItemDraw
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Charge_ItemDraw`;
DELIMITER $$
CREATE PROCEDURE `usp_Charge_ItemDraw`(
    IN in_rowNo INT,
    IN in_userNo INT,
    OUT out_RET INT
)
BEGIN
    DECLARE v_isDraw TINYINT DEFAULT NULL;

    SET out_RET = 0;

    SELECT isDraw INTO v_isDraw
    FROM tChargeItem
    WHERE rowNo = in_rowNo AND userNo = in_userNo
    LIMIT 1;

    IF v_isDraw IS NULL THEN
        SET out_RET = -2;
    ELSEIF v_isDraw <> 0 THEN
        SET out_RET = -1;
    ELSE
        UPDATE tChargeItem
        SET isDraw = 1, drawDate = NOW()
        WHERE rowNo = in_rowNo AND userNo = in_userNo;

        IF ROW_COUNT() = 1 THEN
            SET out_RET = 1;
        ELSE
            SET out_RET = 0;
        END IF;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 12. usp_Charge_ItemInsert
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Charge_ItemInsert`;
DELIMITER $$
CREATE PROCEDURE `usp_Charge_ItemInsert`(
    IN in_userNo INT,
    IN in_orderNo BIGINT,
    IN in_goodsNo INT,
    IN in_amount INT,
    OUT out_RET INT
)
BEGIN
    SET out_RET = 0;

    IF in_orderNo <> 0 AND EXISTS(SELECT 1 FROM tChargeItem WHERE orderNo = in_orderNo) THEN
        SET out_RET = -1;
    ELSE
        INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount, isDraw, registerDate)
        VALUES (in_userNo, in_orderNo, in_goodsNo, in_amount, 0, NOW());

        IF ROW_COUNT() = 1 THEN
            SET out_RET = 1;
        ELSE
            SET out_RET = 0;
        END IF;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 13. usp_Charge_ItemSelectList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Charge_ItemSelectList`;
DELIMITER $$
CREATE PROCEDURE `usp_Charge_ItemSelectList`(
    IN in_userNo INT
)
BEGIN
    SELECT rowNo, goodsNo, amount, registerDate
    FROM tChargeItem
    WHERE userNo = in_userNo AND isDraw = 0
    ORDER BY rowNo ASC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 14. usp_Event_addCash
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Event_addCash`;
DELIMITER $$
CREATE PROCEDURE `usp_Event_addCash`(
    IN in_nUserNo INT,
    IN in_nCash INT,
    IN in_sRef VARCHAR(50),
    OUT out_nRet INT
)
BEGIN
    START TRANSACTION;

    INSERT INTO tCash (userNo, cashType, cash, status)
    VALUES (in_nUserNo, 1, in_nCash, 0)
    ON DUPLICATE KEY UPDATE cash = cash + VALUES(cash), status = 0;

    INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
    VALUES (in_nUserNo, 1, 1, in_nCash, in_sRef);

    COMMIT;
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 15. usp_Item_getBuyCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getBuyCount`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getBuyCount`(
    IN in_userNo INT,
    OUT out_CNT INT
)
BEGIN
    SET out_CNT = 0;
    SELECT IFNULL(COUNT(orderNo), 0) INTO out_CNT
    FROM tOrder
    WHERE userNo = in_userNo AND isPay = 1;

    SELECT out_CNT AS CNT;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 16. usp_Item_getBuyList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getBuyList`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getBuyList`(
    IN in_userNo INT,
    IN in_pageNo INT,
    IN in_pageLimit INT
)
BEGIN
    DECLARE v_offset INT DEFAULT 0;
    SET v_offset = (GREATEST(1, in_pageNo) - 1) * in_pageLimit;

    SELECT DATE_FORMAT(o.orderDate, '%Y-%m-%d') AS orderDate,
           o.orderType,
           o.goodsNo,
           i.name,
           i.imageName,
           i.unit,
           i.unitDescription,
           o.orderPrice
    FROM tOrder o
    JOIN tAccountItem i ON o.goodsNo = i.goodsNo
    WHERE o.userNo = in_userNo AND o.isPay = 1
    ORDER BY o.orderNo DESC
    LIMIT in_pageLimit OFFSET v_offset;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 17. usp_Item_getCategoryName
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getCategoryName`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getCategoryName`(
    IN in_category VARCHAR(10),
    OUT out_CATEGORY_NAME VARCHAR(20)
)
BEGIN
    SET out_CATEGORY_NAME = '';
    SELECT name INTO out_CATEGORY_NAME
    FROM tItemCategory
    WHERE category = in_category
    LIMIT 1;

    IF out_CATEGORY_NAME IS NULL THEN
        SET out_CATEGORY_NAME = '';
    END IF;

    SELECT out_CATEGORY_NAME AS CATEGORY_NAME;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 18. usp_Item_getEventList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getEventList`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getEventList`(
    IN in_eventType TINYINT,
    IN in_pageNo INT,
    IN in_pageLimit INT
)
BEGIN
    DECLARE v_offset INT DEFAULT 0;
    SET v_offset = (GREATEST(1, in_pageNo) - 1) * in_pageLimit;

    SELECT goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount,
           CASE
               WHEN discountStartDate <= NOW() AND discountEndDate >= NOW() THEN price - ((price * discount) / 100)
               ELSE price
           END AS sellPrice
    FROM VT_ITEM_SELL
    WHERE (in_eventType = 1 AND isNew = 1)
       OR (in_eventType = 2 AND isPopular = 1)
       OR (in_eventType = 3 AND isRecommand = 1)
    ORDER BY goodsNo DESC
    LIMIT in_pageLimit OFFSET v_offset;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 19. usp_Item_getGroup
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getGroup`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getGroup`(
    IN in_groupNo INT
)
BEGIN
    SELECT category, name, imageName, itemType, description
    FROM VT_ITEM_GROUP
    WHERE groupNo = in_groupNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 20. usp_Item_getGroupList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getGroupList`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getGroupList`(
    IN in_category VARCHAR(10),
    IN in_pageNo INT,
    IN in_pageLimit INT
)
BEGIN
    DECLARE v_offset INT DEFAULT 0;
    SET v_offset = (GREATEST(1, in_pageNo) - 1) * in_pageLimit;

    IF in_category IS NOT NULL AND in_category <> '' THEN
        SELECT goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount,
               CASE
                   WHEN discountStartDate <= NOW() AND discountEndDate >= NOW() THEN price - ((price * discount) / 100)
                   ELSE price
               END AS sellPrice
        FROM VT_ITEM_GROUP
        WHERE category LIKE CONCAT(in_category, '%')
        ORDER BY goodsNo DESC
        LIMIT in_pageLimit OFFSET v_offset;
    ELSE
        SELECT goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount,
               CASE
                   WHEN discountStartDate <= NOW() AND discountEndDate >= NOW() THEN price - ((price * discount) / 100)
                   ELSE price
               END AS sellPrice
        FROM VT_ITEM_GROUP
        ORDER BY goodsNo DESC
        LIMIT in_pageLimit OFFSET v_offset;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 21. usp_Item_getGroupPrice
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getGroupPrice`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getGroupPrice`(
    IN in_groupNo INT
)
BEGIN
    SELECT goodsNo, billingNo, unit, unitDescription, price, discount,
           CASE
               WHEN discountStartDate <= NOW() AND discountEndDate >= NOW() THEN price - ((price * discount) / 100)
               ELSE price
           END AS sellPrice
    FROM VT_ITEM_SELL
    WHERE groupNo = in_groupNo
    ORDER BY goodsNo ASC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 22. usp_Item_getItem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getItem`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getItem`(
    IN in_goodsNo INT
)
BEGIN
    SELECT goodsNo, billingNo, category, groupNo, name, imageName, itemType, unit, unitDescription, price, discount,
           CASE
               WHEN discountStartDate <= NOW() AND discountEndDate >= NOW() THEN price - ((price * discount) / 100)
               ELSE price
           END AS sellPrice,
           isSell, description
    FROM tAccountItem
    WHERE goodsNo = in_goodsNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 23. usp_Item_getOrder
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getOrder`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getOrder`(
    IN in_orderNo INT
)
BEGIN
    SELECT orderNo, orderType, userNo, toUserNo, goodsNo, amount, orderPrice, cashType, isPay, orderDate
    FROM tOrder
    WHERE orderNo = in_orderNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 24. usp_Item_getOrderList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getOrderList`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getOrderList`(
    IN in_userNo INT
)
BEGIN
    SELECT orderNo, orderType, userNo, toUserNo, goodsNo, amount, orderPrice, cashType, isPay, orderDate
    FROM tOrder
    WHERE userNo = in_userNo
    ORDER BY orderNo DESC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 25. usp_Item_getSubCategory
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getSubCategory`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getSubCategory`(
    IN in_parent VARCHAR(10)
)
BEGIN
    SELECT category, name, parent, status, registerDate
    FROM tItemCategory
    WHERE parent = in_parent
    ORDER BY category ASC;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 26. usp_Item_getTotalCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getTotalCount`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getTotalCount`(
    IN in_eventType INT,
    IN in_category VARCHAR(10),
    OUT out_TOTAL INT
)
BEGIN
    SET out_TOTAL = 0;

    IF in_eventType > 0 THEN
        IF in_eventType = 1 THEN
            SELECT COUNT(goodsNo) INTO out_TOTAL FROM VT_ITEM_SELL WHERE isNew = 1;
        ELSEIF in_eventType = 2 THEN
            SELECT COUNT(goodsNo) INTO out_TOTAL FROM VT_ITEM_SELL WHERE isPopular = 1;
        ELSEIF in_eventType = 3 THEN
            SELECT COUNT(goodsNo) INTO out_TOTAL FROM VT_ITEM_SELL WHERE isRecommand = 1;
        END IF;
    ELSE
        IF in_category IS NOT NULL AND in_category <> '' THEN
            SELECT COUNT(goodsNo) INTO out_TOTAL FROM VT_ITEM_GROUP WHERE category LIKE CONCAT(in_category, '%');
        ELSE
            SELECT COUNT(goodsNo) INTO out_TOTAL FROM VT_ITEM_GROUP;
        END IF;
    END IF;

    IF out_TOTAL IS NULL THEN
        SET out_TOTAL = 0;
    END IF;

    SELECT out_TOTAL AS TOTAL;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 27. usp_Item_getUser
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_getUser`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_getUser`(
    IN in_userNo INT,
    IN in_serialKey BIGINT,
    OUT out_USER_ID VARCHAR(20),
    OUT out_AUTH_ID INT
)
BEGIN
    SET out_USER_ID = '';
    SET out_AUTH_ID = 0;

    IF in_serialKey <> 0 AND EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        SELECT u.sUserID, u.nAuthID
        INTO out_USER_ID, out_AUTH_ID
        FROM tUser u
        JOIN tSerialKey s ON u.nUserNo = s.userNo
        WHERE s.userNo = in_userNo AND s.serialKey = in_serialKey
        LIMIT 1;

        IF out_USER_ID IS NOT NULL AND out_USER_ID <> '' THEN
            DELETE FROM tSerialKey WHERE userNo = in_userNo;
        ELSE
            SET out_USER_ID = '';
            SET out_AUTH_ID = 0;
        END IF;
    END IF;

    SELECT out_USER_ID AS USER_ID, out_AUTH_ID AS AUTH_ID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 28. usp_Item_setBuyOrder
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setBuyOrder`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setBuyOrder`(
    IN in_orderNo INT,
    IN in_userNo INT,
    IN in_targetUser INT,
    OUT out_RET INT
)
BEGIN
    DECLARE v_orderType INT;
    DECLARE v_orderUser INT;
    DECLARE v_goodsNo INT;
    DECLARE v_amount INT;
    DECLARE v_orderPrice INT;
    DECLARE v_cashType TINYINT;
    DECLARE v_isPay INT;
    DECLARE v_owner INT;

    SET out_RET = 0;

    IF NOT EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        SET out_RET = -2;
    ELSE
        SELECT orderType, userNo, goodsNo, amount, orderPrice, cashType, isPay
        INTO v_orderType, v_orderUser, v_goodsNo, v_amount, v_orderPrice, v_cashType, v_isPay
        FROM tOrder
        WHERE orderNo = in_orderNo
        LIMIT 1;

        IF v_orderUser IS NULL THEN
            SET out_RET = -7;
        ELSEIF v_orderType = 1 AND (in_targetUser IS NULL OR in_targetUser = 0) THEN
            SET out_RET = -4;
        ELSEIF v_orderUser <> in_userNo THEN
            SET out_RET = -5;
        ELSEIF v_isPay = 1 THEN
            SET out_RET = -6;
        ELSE
            IF v_orderType = 0 THEN
                SET v_owner = in_userNo;
            ELSE
                SET v_owner = in_targetUser;
            END IF;

            START TRANSACTION;

            UPDATE tOrder
            SET toUserNo = v_owner, isPay = 1
            WHERE orderNo = in_orderNo;

            IF ROW_COUNT() <> 1 THEN
                SET out_RET = -9;
                ROLLBACK;
            ELSE
                -- Call useCash inside transaction
                IF (SELECT IFNULL(cash, 0) FROM tCash WHERE userNo = in_userNo AND cashType = v_cashType AND status = 0 FOR UPDATE) >= v_orderPrice THEN
                    UPDATE tCash
                    SET cash = cash - v_orderPrice
                    WHERE userNo = in_userNo AND cashType = v_cashType AND status = 0;

                    INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
                    VALUES (in_userNo, 2, v_cashType, v_orderPrice, CAST(in_orderNo AS CHAR));

                    INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount, isDraw, registerDate)
                    VALUES (v_owner, in_orderNo, v_goodsNo, v_amount, 0, NOW());

                    UPDATE tAccountItem SET sellCount = sellCount + 1 WHERE goodsNo = v_goodsNo;

                    SET out_RET = 0;
                    COMMIT;
                ELSE
                    SET out_RET = -1; -- not enough cash
                    ROLLBACK;
                END IF;
            END IF;
        END IF;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 29. usp_Item_setCategory
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setCategory`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setCategory`(
    IN in_category VARCHAR(10),
    IN in_name VARCHAR(20),
    IN in_parent VARCHAR(10),
    IN in_status TINYINT(1),
    OUT out_RET INT
)
BEGIN
    SET out_RET = 0;

    IF in_parent <> '' AND NOT EXISTS(SELECT 1 FROM tItemCategory WHERE category = in_parent) THEN
        SET out_RET = -3;
    ELSEIF EXISTS(SELECT 1 FROM tItemCategory WHERE category = in_category) THEN
        UPDATE tItemCategory
        SET name = in_name, parent = in_parent, status = in_status
        WHERE category = in_category;
        SET out_RET = 1;
    ELSE
        INSERT INTO tItemCategory (category, name, parent, status)
        VALUES (in_category, in_name, in_parent, in_status);
        SET out_RET = 1;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 30. usp_Item_setItem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setItem`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setItem`(
    IN in_goodsNo INT,
    IN in_category VARCHAR(10),
    IN in_groupNo INT,
    IN in_name VARCHAR(100),
    IN in_imageName VARCHAR(100),
    IN in_itemType VARCHAR(20),
    IN in_unit INT,
    IN in_price INT,
    IN in_discount INT,
    IN in_isRecommand TINYINT(1),
    IN in_isPopular TINYINT(1),
    IN in_isNew TINYINT(1),
    IN in_isSell TINYINT(1),
    IN in_sellStartDate DATETIME,
    IN in_sellEndDate DATETIME,
    IN in_discountStartDate DATETIME,
    IN in_discountEndDate DATETIME,
    OUT out_RET INT
)
BEGIN
    SET out_RET = 0;

    IF EXISTS(SELECT 1 FROM tAccountItem WHERE goodsNo = in_goodsNo) THEN
        UPDATE tAccountItem SET
            category = in_category,
            groupNo = in_groupNo,
            name = in_name,
            imageName = in_imageName,
            itemType = in_itemType,
            unit = in_unit,
            price = in_price,
            discount = in_discount,
            isRecommand = in_isRecommand,
            isPopular = in_isPopular,
            isNew = in_isNew,
            isSell = in_isSell,
            sellStartDate = in_sellStartDate,
            sellEndDate = in_sellEndDate,
            discountStartDate = in_discountStartDate,
            discountEndDate = in_discountEndDate
        WHERE goodsNo = in_goodsNo;

        IF ROW_COUNT() = 1 THEN
            SET out_RET = 1;
        ELSE
            SET out_RET = -1;
        END IF;
    ELSE
        SET out_RET = -2;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 31. usp_Item_setOrder
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setOrder`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setOrder`(
    IN in_orderType TINYINT,
    IN in_userNo INT,
    IN in_goodsNo INT,
    IN in_goodsPrice INT,
    IN in_amount INT,
    IN in_discount INT,
    IN in_orderPrice INT,
    IN in_cashType INT,
    IN in_userIP VARCHAR(30),
    OUT out_ORDER_NO INT
)
BEGIN
    SET out_ORDER_NO = 0;

    START TRANSACTION;

    INSERT INTO tOrder (orderType, userNo, goodsNo, goodsPrice, amount, discount, orderPrice, cashType, userIP, isPay, orderDate)
    VALUES (in_orderType, in_userNo, in_goodsNo, in_goodsPrice, in_amount, in_discount, in_orderPrice, in_cashType, in_userIP, 0, NOW());

    SET out_ORDER_NO = LAST_INSERT_ID();

    IF out_ORDER_NO > 0 THEN
        COMMIT;
    ELSE
        SET out_ORDER_NO = 0;
        ROLLBACK;
    END IF;

    SELECT out_ORDER_NO AS ORDER_NO;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 32. usp_Item_setSerialKey
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Item_setSerialKey`;
DELIMITER $$
CREATE PROCEDURE `usp_Item_setSerialKey`(
    IN in_userNo INT,
    IN in_serialKey BIGINT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRet = -1;

    IF in_userNo = 0 OR in_serialKey = 0 THEN
        SET out_nRet = -1;
    ELSEIF NOT EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        SET out_nRet = -1;
    ELSEIF EXISTS(SELECT 1 FROM tSerialKey WHERE serialKey = in_serialKey) THEN
        -- Ya existe: rechazar cupón duplicado sin sobrescribir
        SET out_nRet = -1;
    ELSE
        INSERT INTO tSerialKey (userNo, serialKey, registerDate)
        VALUES (in_userNo, in_serialKey, NOW());
        SET out_nRet = 0;
    END IF;

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 33. usp_Pay_chkLimit
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_chkLimit`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_chkLimit`(
    IN in_userNo INT,
    IN in_chkType INT,
    IN in_amount INT
)
BEGIN
    DECLARE v_totalAmount INT DEFAULT 0;
    DECLARE v_ret INT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM VT_PAY WHERE userNo = in_userNo) THEN
        IF in_chkType = 1 THEN
            SELECT IFNULL(SUM(cash), 0) INTO v_totalAmount
            FROM VT_PAY
            WHERE userNo = in_userNo
              AND resultDate > DATE_FORMAT(NOW(), '%Y-%m-%d 00:00:00')
              AND resultDate <= NOW();

            IF v_totalAmount >= in_amount THEN
                SET v_ret = 1;
            ELSE
                SET v_ret = 0;
            END IF;
        ELSE
            SET v_ret = 0;
        END IF;
    ELSE
        SET v_ret = 0;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 34. usp_Pay_getAgreement
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_getAgreement`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_getAgreement`(
    IN in_userNo INT
)
BEGIN
    SELECT status
    FROM tPayAgreement
    WHERE userNo = in_userNo AND vaildDate >= NOW()
    ORDER BY rowNo DESC
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 35. usp_Pay_getOrder
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_getOrder`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_getOrder`(
    IN in_orderNo INT
)
BEGIN
    SELECT p.orderNo, p.userNo, p.payNo, p.cash, p.userIP, p.registerDate,
           pr.result, pr.resultDate, pr.pgTID
    FROM tPay p
    LEFT JOIN tPayResult pr ON p.orderNo = pr.orderNo
    WHERE p.orderNo = in_orderNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 36. usp_Pay_setAgreement
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_setAgreement`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_setAgreement`(
    IN in_userNo INT,
    IN in_name VARCHAR(10),
    IN in_juminNo VARCHAR(20),
    IN in_phoneNo VARCHAR(20),
    IN in_vaildDate INT
)
BEGIN
    INSERT INTO tPayAgreement (userNo, name, juminNo, phoneNo, vaildDate)
    VALUES (in_userNo, in_name, in_juminNo, in_phoneNo, DATE_ADD(NOW(), INTERVAL in_vaildDate MONTH));

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 37. usp_Pay_setOrder
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_setOrder`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_setOrder`(
    IN in_userNo INT,
    IN in_payNo INT,
    IN in_cash INT,
    IN in_userIP VARCHAR(30),
    OUT out_ORDER_NO INT
)
BEGIN
    SET out_ORDER_NO = 0;

    INSERT INTO tPay (userNo, payNo, cash, userIP, registerDate)
    VALUES (in_userNo, in_payNo, in_cash, in_userIP, NOW());

    SET out_ORDER_NO = LAST_INSERT_ID();
    SELECT out_ORDER_NO AS ORDER_NO;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 38. usp_Pay_setResult
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pay_setResult`;
DELIMITER $$
CREATE PROCEDURE `usp_Pay_setResult`(
    IN in_orderNo INT,
    IN in_result INT,
    IN in_pgTID VARCHAR(50)
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_cash INT DEFAULT 0;
    DECLARE v_tid VARCHAR(50);
    DECLARE v_ret INT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM tPayResult WHERE orderNo = in_orderNo) THEN
        SET v_ret = -1;
    ELSE
        SET v_tid = in_pgTID;
        IF v_tid = '' THEN
            SET v_tid = NULL;
        END IF;

        START TRANSACTION;

        INSERT INTO tPayResult (orderNo, result, pgTID, resultDate)
        VALUES (in_orderNo, in_result, v_tid, NOW());

        IF in_result = 0 THEN
            SELECT userNo, cash INTO v_userNo, v_cash FROM tPay WHERE orderNo = in_orderNo LIMIT 1;

            INSERT INTO tCash (userNo, cashType, cash, status)
            VALUES (v_userNo, 1, v_cash, 0)
            ON DUPLICATE KEY UPDATE cash = cash + VALUES(cash), status = 0;

            INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
            VALUES (v_userNo, 1, 1, v_cash, CAST(in_orderNo AS CHAR));
        END IF;

        COMMIT;
        SET v_ret = 0;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 39. usp_Pcbang_chkIP
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Pcbang_chkIP`;
DELIMITER $$
CREATE PROCEDURE `usp_Pcbang_chkIP`(
    IN in_sIP VARCHAR(30),
    OUT out_nPcbangNo INT
)
BEGIN
    DECLARE v_ipNum BIGINT UNSIGNED;
    SET out_nPcbangNo = 0;
    SET v_ipNum = INET_ATON(in_sIP);

    SELECT pcbangNo INTO out_nPcbangNo
    FROM tPcbang
    WHERE bIsActive = 1
      AND v_ipNum BETWEEN INET_ATON(sIPStart) AND INET_ATON(sIPEnd)
    LIMIT 1;

    IF out_nPcbangNo IS NULL THEN
        SET out_nPcbangNo = 0;
    END IF;

    SELECT out_nPcbangNo AS nPcbangNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 40. usp_ReturnUser_GetDatediffLastlog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ReturnUser_GetDatediffLastlog`;
DELIMITER $$
CREATE PROCEDURE `usp_ReturnUser_GetDatediffLastlog`(
    IN in_nUserNo INT,
    OUT out_nDiffSec INT
)
BEGIN
    SET out_nDiffSec = -1;

    SELECT TIMESTAMPDIFF(SECOND, dDate, NOW()) INTO out_nDiffSec
    FROM tReturnUser
    WHERE nUserNo = in_nUserNo
    ORDER BY dDate DESC
    LIMIT 1;

    IF out_nDiffSec IS NULL THEN
        SET out_nDiffSec = -1;
    END IF;

    SELECT out_nDiffSec AS nDiffSec;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 41. usp_ReturnUser_SetUserItem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_ReturnUser_SetUserItem`;
DELIMITER $$
CREATE PROCEDURE `usp_ReturnUser_SetUserItem`(
    IN in_nUserNo INT,
    IN in_nItemNo INT,
    OUT out_nRet INT
)
BEGIN
    DECLARE v_item INT;
    SET v_item = IFNULL(in_nItemNo, 30001);
    IF v_item <= 0 THEN
        SET v_item = 30001;
    END IF;

    START TRANSACTION;

    INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount, isDraw, registerDate)
    VALUES (in_nUserNo, 0, v_item, 1, 0, NOW());

    INSERT INTO tReturnUser (dDate, nUserNo, nItemNo)
    VALUES (NOW(), in_nUserNo, v_item);

    COMMIT;
    SET out_nRet = 0;
    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 42. usp_Stat_sellitem
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Stat_sellitem`;
DELIMITER $$
CREATE PROCEDURE `usp_Stat_sellitem`()
BEGIN
    DECLARE v_start VARCHAR(10);
    DECLARE v_end VARCHAR(10);

    SET v_start = DATE_FORMAT(DATE_SUB(NOW(), INTERVAL 1 DAY), '%Y-%m-%d');
    SET v_end = DATE_FORMAT(NOW(), '%Y-%m-%d');

    INSERT INTO tbSellitem (baseDate, goodsNo, amount)
    SELECT DATE_FORMAT(registerDate, '%Y-%m-%d') AS baseDate,
           goodsNo,
           SUM(amount) AS amount
    FROM tChargeItem
    WHERE registerDate >= v_start AND registerDate < v_end
    GROUP BY DATE_FORMAT(registerDate, '%Y-%m-%d'), goodsNo
    ORDER BY DATE_FORMAT(registerDate, '%Y-%m-%d'), goodsNo;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 43. usp_Tutorial_getCharInfo
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Tutorial_getCharInfo`;
DELIMITER $$
CREATE PROCEDURE `usp_Tutorial_getCharInfo`(
    IN in_nCharNo INT,
    OUT out_nRet INT
)
BEGIN
    SET out_nRet = 0;
    SELECT nState, nStep
    FROM tTutorial
    WHERE nCharNo = in_nCharNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 44. usp_Tutorial_getIsSkip
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_Tutorial_getIsSkip`;
DELIMITER $$
CREATE PROCEDURE `usp_Tutorial_getIsSkip`(
    IN in_nUserNo INT,
    OUT out_nRet INT
)
BEGIN
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
-- 45. usp_Tutorial_setCharInfo
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
BEGIN
    SET out_nRet = 0;

    INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep)
    VALUES (in_nUserNo, in_nCharNo, in_nState, in_nStep)
    ON DUPLICATE KEY UPDATE nState = VALUES(nState), nStep = VALUES(nStep);

    SELECT out_nRet AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 46. usp_User_GetDatediffCreated
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_GetDatediffCreated`;
DELIMITER $$
CREATE PROCEDURE `usp_User_GetDatediffCreated`(
    IN in_nUserNo INT,
    OUT out_nDiffSec INT
)
BEGIN
    SET out_nDiffSec = 0;

    SELECT TIMESTAMPDIFF(SECOND, dDate, NOW()) INTO out_nDiffSec
    FROM tUser
    WHERE nUserNo = in_nUserNo
    LIMIT 1;

    IF out_nDiffSec IS NULL THEN
        SET out_nDiffSec = 0;
    END IF;

    SELECT out_nDiffSec AS diffDate;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 47. usp_User_LoginGame_5
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_LoginGame_5`;
DELIMITER $$
CREATE PROCEDURE `usp_User_LoginGame_5`(
    IN in_sUserID VARCHAR(60),
    IN in_sUserPW VARCHAR(32),
    OUT out_nUserNo INT,
    OUT out_nAuthID TINYINT,
    OUT out_nResult INT
)
BEGIN
    CALL usp_User_loginGame(in_sUserID, in_sUserPW, '127.0.0.1', out_nUserNo, out_nAuthID, out_nResult);
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 48. usp_User_LoginNetmarble
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_LoginNetmarble`;
DELIMITER $$
CREATE PROCEDURE `usp_User_LoginNetmarble`(
    IN in_sUserID VARCHAR(60),
    IN in_sUserPW VARCHAR(32),
    OUT out_nUserNo INT,
    OUT out_nAuthID TINYINT,
    OUT out_nResult INT
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_dbPW VARCHAR(32);
    DECLARE v_isBlock TINYINT DEFAULT 0;
    DECLARE v_isDelete TINYINT DEFAULT 0;
    DECLARE v_authID TINYINT DEFAULT 0;
    DECLARE v_isLoginAble TINYINT DEFAULT 0;
    DECLARE v_blockRet INT DEFAULT 0;

    SET out_nUserNo = 0;
    SET out_nAuthID = 0;
    SET out_nResult = -1;

    SELECT nUserNo, sUserPW, bIsBlock, bIsDelete, nAuthID
    INTO v_userNo, v_dbPW, v_isBlock, v_isDelete, v_authID
    FROM tUser
    WHERE LOWER(sUserID) = LOWER(in_sUserID)
    LIMIT 1;

    IF v_userNo IS NULL OR v_userNo = 0 THEN
        SET out_nUserNo = 0;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF v_isDelete = 1 THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF LOWER(v_dbPW) <> LOWER(in_sUserPW) AND v_dbPW <> in_sUserPW THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -2;
    ELSE
        IF v_isBlock = 1 THEN
            UPDATE tUserBlock
            SET bStatus = 0
            WHERE nUserNo = v_userNo
              AND bStatus = 1
              AND dEndDate IS NOT NULL
              AND dEndDate <= NOW()
              AND nType <> 99;

            SELECT COUNT(*) INTO v_blockRet
            FROM tUserBlock
            WHERE nUserNo = v_userNo AND bStatus = 1;

            IF v_blockRet = 0 THEN
                UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = v_userNo;
                SET v_isBlock = 0;
            ELSE
                UPDATE tUser SET bIsBlock = 1 WHERE nUserNo = v_userNo;
                SET v_isBlock = 1;
            END IF;
        END IF;

        IF v_isBlock = 1 THEN
            SET out_nUserNo = v_userNo;
            SET out_nAuthID = 0;
            SET out_nResult = -3;
        ELSE
            SELECT bIsLoginAble INTO v_isLoginAble FROM tUserAuth WHERE nAuthID = v_authID LIMIT 1;
            IF v_isLoginAble IS NULL OR v_isLoginAble = 0 THEN
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = -4;
            ELSE
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = 0;
            END IF;
        END IF;
    END IF;

    SELECT out_nResult AS nResult, out_nUserNo AS nUserNo, out_nAuthID AS nAuthID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 49. usp_User_USloginGame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_USloginGame`;
DELIMITER $$
CREATE PROCEDURE `usp_User_USloginGame`(
    IN in_sUserID VARCHAR(60),
    IN in_sUserPW VARCHAR(32),
    IN in_sUserIP VARCHAR(30),
    OUT out_nUserNo INT,
    OUT out_nAuthID TINYINT,
    OUT out_nResult INT
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_dbPW VARCHAR(32);
    DECLARE v_isBlock TINYINT DEFAULT 0;
    DECLARE v_isDelete TINYINT DEFAULT 0;
    DECLARE v_authID TINYINT DEFAULT 0;
    DECLARE v_isLoginAble TINYINT DEFAULT 0;
    DECLARE v_blockRet INT DEFAULT 0;
    DECLARE v_trimmedID VARCHAR(60);
    DECLARE v_allowCreation INT DEFAULT 0;

    SET out_nUserNo = 0;
    SET out_nAuthID = 0;
    SET out_nResult = -1;

    -- 1. Auto Account Creation ('r_' prefix)
    IF in_sUserID LIKE 'r[_]%' OR in_sUserID LIKE 'r\_%' THEN
        SELECT Allow_Creation INTO v_allowCreation FROM Auto_Account_Creation LIMIT 1;
        IF v_allowCreation = 1 THEN
            SET v_trimmedID = SUBSTRING(in_sUserID, 3);
            IF NOT EXISTS(SELECT 1 FROM tUser WHERE LOWER(sUserID) = LOWER(v_trimmedID)) THEN
                INSERT INTO tUser (sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate)
                VALUES (v_trimmedID, in_sUserPW, v_trimmedID, 0, 0, 1, IFNULL(in_sUserIP, '127.0.0.1'), NOW());

                SET v_userNo = LAST_INSERT_ID();

                INSERT INTO tUserOption (nUserNo, sEmail)
                VALUES (v_userNo, CONCAT(v_trimmedID, '@xelion.local'));
            END IF;
            SET in_sUserID = v_trimmedID;
        END IF;
    END IF;

    -- 2. Lookup in tUser
    SELECT nUserNo, sUserPW, bIsBlock, bIsDelete, nAuthID
    INTO v_userNo, v_dbPW, v_isBlock, v_isDelete, v_authID
    FROM tUser
    WHERE LOWER(sUserID) = LOWER(in_sUserID)
    LIMIT 1;

    -- 3. Fallback: Lookup in tUserOption by sEmail (with bIsDelete = 0)
    IF v_userNo IS NULL OR v_userNo = 0 THEN
        SELECT u.nUserNo, u.sUserPW, u.bIsBlock, u.bIsDelete, u.nAuthID
        INTO v_userNo, v_dbPW, v_isBlock, v_isDelete, v_authID
        FROM tUser u
        JOIN tUserOption o ON u.nUserNo = o.nUserNo
        WHERE LOWER(o.sEmail) = LOWER(in_sUserID) AND u.bIsDelete = 0
        LIMIT 1;
    END IF;

    -- 4. Result evaluation
    IF v_userNo IS NULL OR v_userNo = 0 THEN
        SET out_nUserNo = 0;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF v_isDelete = 1 THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF LOWER(v_dbPW) <> LOWER(in_sUserPW) AND v_dbPW <> in_sUserPW THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -2;
    ELSE
        IF v_isBlock = 1 THEN
            UPDATE tUserBlock
            SET bStatus = 0
            WHERE nUserNo = v_userNo
              AND bStatus = 1
              AND dEndDate IS NOT NULL
              AND dEndDate <= NOW()
              AND nType <> 99;

            SELECT COUNT(*) INTO v_blockRet
            FROM tUserBlock
            WHERE nUserNo = v_userNo AND bStatus = 1;

            IF v_blockRet = 0 THEN
                UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = v_userNo;
                SET v_isBlock = 0;
            ELSE
                UPDATE tUser SET bIsBlock = 1 WHERE nUserNo = v_userNo;
                SET v_isBlock = 1;
            END IF;
        END IF;

        IF v_isBlock = 1 THEN
            SET out_nUserNo = v_userNo;
            SET out_nAuthID = 0;
            SET out_nResult = -3;
        ELSE
            SELECT bIsLoginAble INTO v_isLoginAble FROM tUserAuth WHERE nAuthID = v_authID LIMIT 1;
            IF v_isLoginAble IS NULL OR v_isLoginAble = 0 THEN
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = -4;
            ELSE
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = 0;
            END IF;
        END IF;
    END IF;

    SELECT out_nResult AS nResult, out_nUserNo AS nUserNo, out_nAuthID AS nAuthID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 50. usp_User_USloginGame_O
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_USloginGame_O`;
DELIMITER $$
CREATE PROCEDURE `usp_User_USloginGame_O`(
    IN in_sUserID VARCHAR(60),
    IN in_sUserPW VARCHAR(32),
    OUT out_nUserNo INT,
    OUT out_nAuthID TINYINT,
    OUT out_nResult INT
)
BEGIN
    CALL usp_User_LoginGame_5(in_sUserID, in_sUserPW, out_nUserNo, out_nAuthID, out_nResult);
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 51. usp_User_chkBlock
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_chkBlock`;
DELIMITER $$
CREATE PROCEDURE `usp_User_chkBlock`(
    IN in_userNo INT,
    OUT out_RET INT
)
BEGIN
    DECLARE v_activeBlocks INT DEFAULT 0;

    -- Expirar bloqueos temporales cuyo dEndDate <= NOW() y nType != 99
    UPDATE tUserBlock
    SET bStatus = 0
    WHERE nUserNo = in_userNo
      AND bStatus = 1
      AND dEndDate IS NOT NULL
      AND dEndDate <= NOW()
      AND nType <> 99;

    -- Verificar si quedan otros bloqueos activos
    SELECT COUNT(*) INTO v_activeBlocks
    FROM tUserBlock
    WHERE nUserNo = in_userNo AND bStatus = 1;

    IF v_activeBlocks = 0 THEN
        UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = in_userNo;
        SET out_RET = 1;
    ELSE
        UPDATE tUser SET bIsBlock = 1 WHERE nUserNo = in_userNo;
        SET out_RET = 0;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 52. usp_User_chkRegister
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_chkRegister`;
DELIMITER $$
CREATE PROCEDURE `usp_User_chkRegister`(
    IN in_juminNo VARCHAR(20),
    OUT out_age INT,
    OUT out_registDate VARCHAR(10)
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_century VARCHAR(2);
    DECLARE v_birthday VARCHAR(10);
    DECLARE v_genderDigit CHAR(1);
    DECLARE v_ret INT DEFAULT 0;

    SET out_age = 0;
    SET out_registDate = '';

    IF in_juminNo IS NULL OR in_juminNo = '' THEN
        SET v_ret = -1;
    ELSE
        SET v_genderDigit = SUBSTRING(in_juminNo, 8, 1);
        IF v_genderDigit IN ('1', '2') THEN
            SET v_century = '19';
        ELSEIF v_genderDigit IN ('3', '4') THEN
            SET v_century = '20';
        ELSE
            SET v_century = '';
        END IF;

        IF v_century = '' THEN
            SET v_ret = -1;
        ELSE
            SET v_birthday = CONCAT(v_century, SUBSTRING(in_juminNo, 1, 2), '-', SUBSTRING(in_juminNo, 3, 2), '-', SUBSTRING(in_juminNo, 5, 2));
            SET out_age = TIMESTAMPDIFF(YEAR, STR_TO_DATE(v_birthday, '%Y-%m-%d'), NOW());

            IF out_age < 15 THEN
                SET v_ret = -2;
            ELSE
                SELECT u.nUserNo, DATE_FORMAT(u.dDate, '%Y-%m-%d')
                INTO v_userNo, out_registDate
                FROM tUser u
                JOIN tUserOption o ON u.nUserNo = o.nUserNo
                WHERE o.sJuminNo = in_juminNo AND u.bIsDelete = 0
                LIMIT 1;

                IF v_userNo > 100 THEN
                    SET v_ret = -3;
                ELSE
                    SET v_ret = 0;
                END IF;
            END IF;
        END IF;
    END IF;

    SELECT v_ret AS nRet, out_age AS age, out_registDate AS registDate;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 53. usp_User_chkUserID
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_chkUserID`;
DELIMITER $$
CREATE PROCEDURE `usp_User_chkUserID`(
    IN in_userID VARCHAR(60)
)
BEGIN
    DECLARE v_ret INT DEFAULT 0;

    IF in_userID IS NULL OR CHAR_LENGTH(in_userID) < 5 OR CHAR_LENGTH(in_userID) > 16 THEN
        SET v_ret = -1;
    ELSEIF EXISTS(SELECT 1 FROM tUser WHERE sUserID = in_userID AND bIsDelete = 0) THEN
        SET v_ret = -2;
    ELSE
        SET v_ret = 0;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 54. usp_User_delete
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_delete`;
DELIMITER $$
CREATE PROCEDURE `usp_User_delete`(
    IN in_userNo INT
)
BEGIN
    DECLARE v_dateStr VARCHAR(8);
    DECLARE v_ret INT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        SET v_dateStr = DATE_FORMAT(NOW(), '%Y%m%d');

        START TRANSACTION;

        UPDATE tUser
        SET sUserID = CONCAT('~', sUserID, '_', v_dateStr),
            bIsDelete = 1
        WHERE nUserNo = in_userNo;

        UPDATE tUserOption
        SET sJuminNo = CONCAT('~', sJuminNo, '_', v_dateStr)
        WHERE nUserNo = in_userNo;

        COMMIT;
        SET v_ret = 0;
    ELSE
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 55. usp_User_getAge
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getAge`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getAge`(
    IN in_userNo INT,
    OUT out_AGE INT
)
BEGIN
    DECLARE v_juminNo VARCHAR(20);
    DECLARE v_genderDigit CHAR(1);
    DECLARE v_century VARCHAR(2);
    DECLARE v_birthday VARCHAR(10);

    SET out_AGE = 0;

    SELECT sJuminNo INTO v_juminNo
    FROM tUserOption
    WHERE nUserNo = in_userNo
    LIMIT 1;

    IF v_juminNo IS NOT NULL AND CHAR_LENGTH(v_juminNo) >= 8 THEN
        SET v_genderDigit = SUBSTRING(v_juminNo, 8, 1);
        IF v_genderDigit IN ('1', '2') THEN
            SET v_century = '19';
        ELSEIF v_genderDigit IN ('3', '4') THEN
            SET v_century = '20';
        ELSE
            SET v_century = '';
        END IF;

        IF v_century <> '' THEN
            SET v_birthday = CONCAT(v_century, SUBSTRING(v_juminNo, 1, 2), '-', SUBSTRING(v_juminNo, 3, 2), '-', SUBSTRING(v_juminNo, 5, 2));
            SET out_AGE = TIMESTAMPDIFF(YEAR, STR_TO_DATE(v_birthday, '%Y-%m-%d'), NOW());
        END IF;
    END IF;

    SELECT out_AGE AS AGE;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 56. usp_User_getAuthList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getAuthList`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getAuthList`()
BEGIN
    SELECT nAuthID, sAuthName, bIsLoginAble FROM tUserAuth;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 57. usp_User_getBuyCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getBuyCount`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getBuyCount`(
    IN in_userNo INT,
    OUT out_CNT INT
)
BEGIN
    SET out_CNT = 0;
    SELECT IFNULL(COUNT(orderNo), 0) INTO out_CNT
    FROM tOrder
    WHERE userNo = in_userNo AND isPay = 1;

    SELECT out_CNT AS CNT;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 58. usp_User_getBuyList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getBuyList`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getBuyList`(
    IN in_userNo INT,
    IN in_pageNo INT,
    IN in_pageLimit INT
)
BEGIN
    DECLARE v_offset INT DEFAULT 0;
    SET v_offset = (GREATEST(1, in_pageNo) - 1) * in_pageLimit;

    SELECT DATE_FORMAT(o.orderDate, '%Y-%m-%d') AS orderDate,
           o.orderType,
           o.goodsNo,
           i.name,
           i.imageName,
           i.unit,
           i.unitDescription,
           o.orderPrice
    FROM tOrder o
    JOIN tAccountItem i ON o.goodsNo = i.goodsNo
    WHERE o.userNo = in_userNo AND o.isPay = 1
    ORDER BY o.orderNo DESC
    LIMIT in_pageLimit OFFSET v_offset;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 59. usp_User_getLogin
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getLogin`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getLogin`(
    IN in_userID VARCHAR(60),
    IN in_userPW VARCHAR(32)
)
BEGIN
    SELECT u.nUserNo, u.sUserID, u.sUserPW, u.sUserName, u.bIsBlock, u.nAuthID,
           o.sEmail, o.sJuminNo, o.sPhoneNo
    FROM tUser u
    LEFT JOIN tUserOption o ON u.nUserNo = o.nUserNo
    WHERE u.sUserID = in_userID AND u.sUserPW = in_userPW AND u.bIsDelete = 0
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 60. usp_User_getOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getOption`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getOption`(
    IN in_userNo INT
)
BEGIN
    SELECT o.sEmail, o.sJuminNo, o.sPhoneNo, o.nZipNo,
           z.sZIPCODE, z.sSIDO, z.sGUGUN, z.sDONG,
           o.sAddress, o.bIsMail
    FROM tUserOption o
    LEFT JOIN tZipCode z ON o.nZipNo = z.nZipNo
    WHERE o.nUserNo = in_userNo
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 61. usp_User_getPayCount
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getPayCount`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getPayCount`(
    IN in_userNo INT,
    OUT out_CNT INT
)
BEGIN
    SET out_CNT = 0;
    SELECT IFNULL(COUNT(orderNo), 0) INTO out_CNT
    FROM VT_PAY
    WHERE userNo = in_userNo;

    SELECT out_CNT AS CNT;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 62. usp_User_getPayList
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getPayList`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getPayList`(
    IN in_userNo INT,
    IN in_pageNo INT,
    IN in_pageLimit INT
)
BEGIN
    DECLARE v_offset INT DEFAULT 0;
    SET v_offset = (GREATEST(1, in_pageNo) - 1) * in_pageLimit;

    SELECT DATE_FORMAT(p.resultDate, '%Y-%m-%d') AS resultDate,
           pt.name,
           p.cash
    FROM VT_PAY p
    JOIN tPayType pt ON p.payNo = pt.payNo
    WHERE p.userNo = in_userNo
    ORDER BY p.orderNo DESC
    LIMIT in_pageLimit OFFSET v_offset;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 63. usp_User_getUser
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_getUser`;
DELIMITER $$
CREATE PROCEDURE `usp_User_getUser`(
    IN in_user VARCHAR(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci
)
BEGIN
    SELECT nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate
    FROM tUser
    WHERE sUserID = in_user OR CAST(nUserNo AS CHAR) = in_user
    LIMIT 1;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 64. usp_User_insert
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_insert`;
DELIMITER $$
CREATE PROCEDURE `usp_User_insert`(
    IN in_userID VARCHAR(20),
    IN in_userPW VARCHAR(32),
    IN in_nAffiliateID VARCHAR(50),
    IN in_szAffiliateSubID VARCHAR(50),
    IN in_szCountryCode VARCHAR(50),
    IN in_userName VARCHAR(60),
    IN in_userIP VARCHAR(30),
    IN in_eMail VARCHAR(70),
    IN in_juminNo VARCHAR(20),
    IN in_phoneNo VARCHAR(20),
    IN in_zipNo INT,
    IN in_address VARCHAR(90),
    IN in_isMail TINYINT(1),
    IN in_channel VARCHAR(20),
    IN in_isOption TINYINT(1),
    OUT out_userNo INT
)
BEGIN
    DECLARE v_itemNo INT DEFAULT 30000;

    SET out_userNo = 0;

    IF CHAR_LENGTH(in_userID) < 5 OR CHAR_LENGTH(in_userID) > 20 THEN
        SELECT -2 AS nRet, 0 AS userNo;
    ELSEIF EXISTS(SELECT 1 FROM tUser WHERE sUserID = in_userID AND bIsDelete = 0) THEN
        SELECT -2 AS nRet, 0 AS userNo;
    ELSE
        START TRANSACTION;

        INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP, bIsBlock, bIsDelete, nAuthID, dDate)
        VALUES (in_userID, in_userPW, in_userName, in_userIP, 0, 0, 1, NOW());

        SET out_userNo = LAST_INSERT_ID();

        IF in_isOption = 1 THEN
            INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel, nAffiliateID, szAffiliateSubID, szCountryCode)
            VALUES (out_userNo, in_eMail, in_juminNo, in_phoneNo, in_zipNo, in_address, in_isMail, IFNULL(in_channel, 'ONSON'), in_nAffiliateID, in_szAffiliateSubID, in_szCountryCode);
        END IF;

        -- Reward starter item
        INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount, isDraw, registerDate)
        VALUES (out_userNo, 0, v_itemNo, 1, 0, NOW());

        COMMIT;
        SELECT 0 AS nRet, out_userNo AS userNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 65. usp_User_insertOption
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_insertOption`;
DELIMITER $$
CREATE PROCEDURE `usp_User_insertOption`(
    IN in_userNo INT,
    IN in_eMail VARCHAR(256),
    IN in_juminNo VARCHAR(20),
    IN in_phoneNo VARCHAR(20),
    IN in_zipNo INT,
    IN in_address VARCHAR(90),
    IN in_isMail TINYINT(1),
    IN in_channel VARCHAR(20),
    IN in_nAffiliateID VARCHAR(50)
)
BEGIN
    INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel, nAffiliateID)
    VALUES (in_userNo, in_eMail, in_juminNo, in_phoneNo, in_zipNo, in_address, in_isMail, in_channel, in_nAffiliateID);

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 66. usp_User_insertOption_Old
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_insertOption_Old`;
DELIMITER $$
CREATE PROCEDURE `usp_User_insertOption_Old`(
    IN in_userNo INT,
    IN in_eMail VARCHAR(70),
    IN in_juminNo VARCHAR(20),
    IN in_phoneNo VARCHAR(20),
    IN in_zipNo INT,
    IN in_address VARCHAR(90),
    IN in_isMail TINYINT(1),
    IN in_channel VARCHAR(20)
)
BEGIN
    INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel)
    VALUES (in_userNo, in_eMail, in_juminNo, in_phoneNo, in_zipNo, in_address, in_isMail, in_channel);

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 67. usp_User_insert_Old
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_insert_Old`;
DELIMITER $$
CREATE PROCEDURE `usp_User_insert_Old`(
    IN in_userID VARCHAR(60),
    IN in_userPW VARCHAR(32),
    IN in_userName VARCHAR(60),
    IN in_userIP VARCHAR(30),
    IN in_eMail VARCHAR(70),
    IN in_juminNo VARCHAR(20),
    IN in_phoneNo VARCHAR(20),
    IN in_zipNo INT,
    IN in_address VARCHAR(90),
    IN in_isMail TINYINT(1),
    IN in_channel VARCHAR(20),
    IN in_isOption TINYINT(1),
    OUT out_userNo INT
)
BEGIN
    SET out_userNo = 0;

    IF CHAR_LENGTH(in_userID) < 5 OR CHAR_LENGTH(in_userID) > 60 THEN
        SELECT -2 AS nRet;
    ELSEIF EXISTS(SELECT 1 FROM tUser WHERE sUserID = in_userID AND bIsDelete = 0) THEN
        SELECT -2 AS nRet;
    ELSE
        START TRANSACTION;

        INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP, bIsBlock, bIsDelete, nAuthID, dDate)
        VALUES (in_userID, in_userPW, in_userName, in_userIP, 0, 0, 1, NOW());

        SET out_userNo = LAST_INSERT_ID();

        IF in_isOption = 1 THEN
            INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel)
            VALUES (out_userNo, in_eMail, in_juminNo, in_phoneNo, in_zipNo, in_address, in_isMail, IFNULL(in_channel, 'ONSON'));
        END IF;

        COMMIT;
        SELECT 0 AS nRet, out_userNo AS userNo;
    END IF;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 68. usp_User_loginGame
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_loginGame`;
DELIMITER $$
CREATE PROCEDURE `usp_User_loginGame`(
    IN in_sUserID VARCHAR(60),
    IN in_sUserPW VARCHAR(32),
    IN in_sUserIP VARCHAR(30),
    OUT out_nUserNo INT,
    OUT out_nAuthID TINYINT,
    OUT out_nResult INT
)
BEGIN
    DECLARE v_userNo INT DEFAULT 0;
    DECLARE v_dbPW VARCHAR(32);
    DECLARE v_isBlock TINYINT DEFAULT 0;
    DECLARE v_isDelete TINYINT DEFAULT 0;
    DECLARE v_authID TINYINT DEFAULT 0;
    DECLARE v_isLoginAble TINYINT DEFAULT 0;
    DECLARE v_blockRet INT DEFAULT 0;

    SET out_nUserNo = 0;
    SET out_nAuthID = 0;
    SET out_nResult = -1;

    SELECT nUserNo, sUserPW, bIsBlock, bIsDelete, nAuthID
    INTO v_userNo, v_dbPW, v_isBlock, v_isDelete, v_authID
    FROM tUser
    WHERE LOWER(sUserID) = LOWER(in_sUserID)
    LIMIT 1;

    IF v_userNo IS NULL OR v_userNo = 0 THEN
        SET out_nUserNo = 0;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF v_isDelete = 1 THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -1;
    ELSEIF LOWER(v_dbPW) <> LOWER(in_sUserPW) AND v_dbPW <> in_sUserPW THEN
        SET out_nUserNo = v_userNo;
        SET out_nAuthID = 0;
        SET out_nResult = -2;
    ELSE
        IF v_isBlock = 1 THEN
            UPDATE tUserBlock
            SET bStatus = 0
            WHERE nUserNo = v_userNo
              AND bStatus = 1
              AND dEndDate IS NOT NULL
              AND dEndDate <= NOW()
              AND nType <> 99;

            SELECT COUNT(*) INTO v_blockRet
            FROM tUserBlock
            WHERE nUserNo = v_userNo AND bStatus = 1;

            IF v_blockRet = 0 THEN
                UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = v_userNo;
                SET v_isBlock = 0;
            ELSE
                UPDATE tUser SET bIsBlock = 1 WHERE nUserNo = v_userNo;
                SET v_isBlock = 1;
            END IF;
        END IF;

        IF v_isBlock = 1 THEN
            SET out_nUserNo = v_userNo;
            SET out_nAuthID = 0;
            SET out_nResult = -3;
        ELSE
            SELECT bIsLoginAble INTO v_isLoginAble FROM tUserAuth WHERE nAuthID = v_authID LIMIT 1;
            IF v_isLoginAble IS NULL OR v_isLoginAble = 0 THEN
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = -4;
            ELSE
                SET out_nUserNo = v_userNo;
                SET out_nAuthID = v_authID;
                SET out_nResult = 0;
            END IF;
        END IF;
    END IF;

    SELECT out_nResult AS nResult, out_nUserNo AS nUserNo, out_nAuthID AS nAuthID;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 69. usp_User_search
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_search`;
DELIMITER $$
CREATE PROCEDURE `usp_User_search`(
    IN in_userName VARCHAR(60),
    IN in_juminNo VARCHAR(20)
)
BEGIN
    SELECT u.nUserNo, u.sUserID, u.sUserPW, o.sEmail
    FROM tUser u
    JOIN tUserOption o ON u.nUserNo = o.nUserNo
    WHERE u.sUserName = in_userName
      AND o.sJuminNo = in_juminNo
      AND u.bIsDelete = 0;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 70. usp_User_searchZipcode
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_searchZipcode`;
DELIMITER $$
CREATE PROCEDURE `usp_User_searchZipcode`(
    IN in_dong VARCHAR(50)
)
BEGIN
    SELECT *
    FROM tZipCode
    WHERE sDONG LIKE CONCAT(in_dong, '%');
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 71. usp_User_setChannelLog
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_setChannelLog`;
DELIMITER $$
CREATE PROCEDURE `usp_User_setChannelLog`(
    IN in_userNo INT,
    IN in_channel VARCHAR(20),
    IN in_issueDate VARCHAR(20)
)
BEGIN
    INSERT INTO tChannelLog (userNo, channel, issueDate)
    VALUES (in_userNo, in_channel, in_issueDate);

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 72. usp_User_setDrawComment
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_setDrawComment`;
DELIMITER $$
CREATE PROCEDURE `usp_User_setDrawComment`(
    IN in_userNo INT,
    IN in_causeNo TINYINT,
    IN in_comment VARCHAR(1000),
    IN in_userIP VARCHAR(30)
)
BEGIN
    INSERT INTO tDrawComment (nUserNo, nCauseNo, sComment, sUserIP, dDate)
    VALUES (in_userNo, in_causeNo, in_comment, in_userIP, NOW());

    SELECT 0 AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 73. usp_User_setLoginable
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_setLoginable`;
DELIMITER $$
CREATE PROCEDURE `usp_User_setLoginable`(
    IN in_authID INT,
    IN in_isLoginable INT,
    OUT out_RET INT
)
BEGIN
    SET out_RET = 0;

    IF EXISTS(SELECT 1 FROM tUserAuth WHERE nAuthID = in_authID) THEN
        UPDATE tUserAuth SET bIsLoginAble = in_isLoginable WHERE nAuthID = in_authID;
        IF ROW_COUNT() = 1 THEN
            SET out_RET = 1;
        ELSE
            SET out_RET = 0;
        END IF;
    ELSE
        SET out_RET = 0;
    END IF;

    SELECT out_RET AS RET;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 74. usp_User_update
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_update`;
DELIMITER $$
CREATE PROCEDURE `usp_User_update`(
    IN in_userNo INT,
    IN in_eMail VARCHAR(200),
    IN in_phoneNo VARCHAR(20),
    IN in_zipNo INT,
    IN in_address VARCHAR(90),
    IN in_isMail TINYINT(1)
)
BEGIN
    DECLARE v_ret INT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        UPDATE tUserOption
        SET sEmail = in_eMail,
            sPhoneNo = in_phoneNo,
            nZipNo = in_zipNo,
            sAddress = in_address,
            bIsMail = in_isMail
        WHERE nUserNo = in_userNo;

        SET v_ret = 0;
    ELSE
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;

-- ----------------------------------------------------------------------------
-- 75. usp_User_updatePassword
-- ----------------------------------------------------------------------------
DROP PROCEDURE IF EXISTS `usp_User_updatePassword`;
DELIMITER $$
CREATE PROCEDURE `usp_User_updatePassword`(
    IN in_userNo INT,
    IN in_prevPW VARCHAR(32),
    IN in_userPW VARCHAR(32),
    IN in_userIP VARCHAR(30)
)
BEGIN
    DECLARE v_currPW VARCHAR(32);
    DECLARE v_ret INT DEFAULT 0;

    IF EXISTS(SELECT 1 FROM tUser WHERE nUserNo = in_userNo AND bIsDelete = 0) THEN
        SELECT sUserPW INTO v_currPW FROM tUser WHERE nUserNo = in_userNo LIMIT 1;

        IF LOWER(v_currPW) <> LOWER(in_prevPW) AND v_currPW <> in_prevPW THEN
            SET v_ret = -2;
        ELSE
            START TRANSACTION;

            UPDATE tUser SET sUserPW = in_userPW WHERE nUserNo = in_userNo;

            INSERT INTO tPasswordLog (nUserNo, sPrevPW, sCurrPW, sUserIP, dDate)
            VALUES (in_userNo, v_currPW, in_userPW, in_userIP, NOW());

            COMMIT;
            SET v_ret = 0;
        END IF;
    ELSE
        SET v_ret = -1;
    END IF;

    SELECT v_ret AS nRet;
END$$
DELIMITER ;
