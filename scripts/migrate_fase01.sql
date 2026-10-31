-- ============================================================================
-- XELION ONLINE - MIGRACIÓN FASE 01 (NO DESTRUCTIVA)
-- Saneamiento 100% de tablas, columnas, índices, vistas y datos semilla
-- ============================================================================

USE `xelion_online`;

-- 1. tUser: Drop bIsLoginAble, add sUserNo, blsBlock, blsDelete
ALTER TABLE `tUser` DROP COLUMN IF EXISTS `bIsLoginAble`;
ALTER TABLE `tUser` ADD COLUMN IF NOT EXISTS `sUserNo` INT DEFAULT NULL;
ALTER TABLE `tUser` ADD COLUMN IF NOT EXISTS `blsBlock` TINYINT NOT NULL DEFAULT 0;
ALTER TABLE `tUser` ADD COLUMN IF NOT EXISTS `blsDelete` TINYINT NOT NULL DEFAULT 0;

-- 2. tCash: status TINYINT(1) NOT NULL DEFAULT 0, PK (userNo, cashType)
ALTER TABLE `tCash` MODIFY COLUMN `status` TINYINT(1) NOT NULL DEFAULT 0;
UPDATE `tCash` SET `status` = 0 WHERE `status` = 1;

-- 3. tCashLog: rename/unify validDate -> vaildDate
SET @exist := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'xelion_online' AND TABLE_NAME = 'tCashLog' AND COLUMN_NAME = 'validDate');
SET @sql := IF(@exist > 0, 'ALTER TABLE `tCashLog` CHANGE COLUMN `validDate` `vaildDate` DATETIME DEFAULT NULL', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- 4. tChargeItem: orderNo BIGINT NOT NULL DEFAULT 0
ALTER TABLE `tChargeItem` MODIFY COLUMN `orderNo` BIGINT NOT NULL DEFAULT 0;

-- 5. tChargeItemDelLog: rowNo without auto_increment, orderNo BIGINT NOT NULL DEFAULT 0
ALTER TABLE `tChargeItemDelLog` MODIFY COLUMN `rowNo` INT NOT NULL;
ALTER TABLE `tChargeItemDelLog` MODIFY COLUMN `orderNo` BIGINT NOT NULL DEFAULT 0;

-- 6. tDrawComment: sComment VARCHAR(1000)
ALTER TABLE `tDrawComment` MODIFY COLUMN `sComment` VARCHAR(1000) NOT NULL DEFAULT '';

-- 7. tOrder: pcbangNo INT DEFAULT NULL
ALTER TABLE `tOrder` MODIFY COLUMN `pcbangNo` INT DEFAULT NULL;

-- 8. tPay: orderNo with AUTO_INCREMENT and PK (orderNo)
ALTER TABLE `tPay` DROP PRIMARY KEY, ADD PRIMARY KEY (`orderNo`), MODIFY COLUMN `orderNo` INT NOT NULL AUTO_INCREMENT;

-- 9. tPayResult: pgTID VARCHAR(50) DEFAULT NULL
ALTER TABLE `tPayResult` MODIFY COLUMN `pgTID` VARCHAR(50) DEFAULT NULL;

-- 10. tSerialKey: PK (serialKey) and KEY (userNo)
ALTER TABLE `tSerialKey` DROP PRIMARY KEY, ADD PRIMARY KEY (`serialKey`), ADD KEY IF NOT EXISTS `idx_tserialkey_userno` (`userNo`);

-- 11. tUserOption: sEmail VARCHAR(256) with index, nullable fields
ALTER TABLE `tUserOption` MODIFY COLUMN `sEmail` VARCHAR(256) NOT NULL DEFAULT '';
ALTER TABLE `tUserOption` MODIFY COLUMN `sJuminNo` VARCHAR(50) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `sPhoneNo` VARCHAR(20) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `nZipNo` INT DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `sAddress` VARCHAR(90) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `sChannel` VARCHAR(20) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `nAffiliateID` VARCHAR(50) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `szAffiliateSubID` VARCHAR(50) DEFAULT NULL;
ALTER TABLE `tUserOption` MODIFY COLUMN `szCountryCode` VARCHAR(50) DEFAULT NULL;
ALTER TABLE `tUserOption` ADD KEY IF NOT EXISTS `idx_tuseroption_email` (`sEmail`);

-- 12. tbSellitem: baseDate VARCHAR(10) NOT NULL, PK (baseDate, goodsNo)
ALTER TABLE `tbSellitem` MODIFY COLUMN `baseDate` VARCHAR(10) NOT NULL;

-- 13. tAccountLog: sUserID VARCHAR(60), sUserPassword VARCHAR(32), sIP VARCHAR(30), indices
ALTER TABLE `tAccountLog` MODIFY COLUMN `sUserID` VARCHAR(60) NOT NULL;
ALTER TABLE `tAccountLog` MODIFY COLUMN `sUserPassword` VARCHAR(32) NOT NULL DEFAULT '';
ALTER TABLE `tAccountLog` MODIFY COLUMN `sIP` VARCHAR(30) NOT NULL DEFAULT '';
ALTER TABLE `tAccountLog` ADD KEY IF NOT EXISTS `idx_taccountlog_userno_type_date` (`nUserNo`, `nType`, `dDate`);
ALTER TABLE `tAccountLog` ADD KEY IF NOT EXISTS `idx_taccountlog_date_type` (`dDate`, `nType`);

-- 14. tEventWorldCup: PK (nEventCharNo, nLevel, nWorldNo)
ALTER TABLE `tEventWorldCup` DROP PRIMARY KEY, ADD PRIMARY KEY (`nEventCharNo`, `nLevel`, `nWorldNo`);

-- 15. tEvent_Friend: nIndex AUTO_INCREMENT PK, dRegDate, dCompleteDate
ALTER TABLE `tEvent_Friend` DROP PRIMARY KEY;
ALTER TABLE `tEvent_Friend` ADD COLUMN IF NOT EXISTS `nIndex` INT NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST;
ALTER TABLE `tEvent_Friend` ADD COLUMN IF NOT EXISTS `dRegDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE `tEvent_Friend` ADD COLUMN IF NOT EXISTS `dCompleteDate` DATETIME DEFAULT NULL;
ALTER TABLE `tEvent_Friend` ADD KEY IF NOT EXISTS `idx_event_friend_users` (`nEventUserNo`, `nFriendUserNo`, `nWorldNo`);

-- 16. tIPBlock: PK (sIPStart, sIPEnd), generated columns nIPStartNum, nIPEndNum
SET @has_blockno := (SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = 'xelion_online' AND TABLE_NAME = 'tIPBlock' AND COLUMN_NAME = 'nBlockNo');
SET @sql := IF(@has_blockno > 0, 'ALTER TABLE `tIPBlock` DROP PRIMARY KEY, DROP COLUMN `nBlockNo`', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

ALTER TABLE `tIPBlock` MODIFY COLUMN `sIPStart` VARCHAR(20) NOT NULL;
ALTER TABLE `tIPBlock` MODIFY COLUMN `sIPEnd` VARCHAR(20) NOT NULL;
ALTER TABLE `tIPBlock` MODIFY COLUMN `sDesc` VARCHAR(30) NOT NULL DEFAULT '';

SET @has_pk := (SELECT COUNT(*) FROM information_schema.KEY_COLUMN_USAGE WHERE TABLE_SCHEMA = 'xelion_online' AND TABLE_NAME = 'tIPBlock' AND CONSTRAINT_NAME = 'PRIMARY');
SET @sql := IF(@has_pk = 0, 'ALTER TABLE `tIPBlock` ADD PRIMARY KEY (`sIPStart`, `sIPEnd`)', 'SELECT 1');
PREPARE stmt FROM @sql; EXECUTE stmt; DEALLOCATE PREPARE stmt;

ALTER TABLE `tIPBlock` ADD COLUMN IF NOT EXISTS `nIPStartNum` INT UNSIGNED GENERATED ALWAYS AS (INET_ATON(`sIPStart`)) STORED;
ALTER TABLE `tIPBlock` ADD COLUMN IF NOT EXISTS `nIPEndNum` INT UNSIGNED GENERATED ALWAYS AS (INET_ATON(`sIPEnd`)) STORED;
ALTER TABLE `tIPBlock` ADD KEY IF NOT EXISTS `idx_tipblock_num_range` (`nIPStartNum`, `nIPEndNum`);

-- 17. tLogout: level SMALLINT UNSIGNED, index (userNo, logDate DESC)
ALTER TABLE `tLogout` MODIFY COLUMN `level` SMALLINT UNSIGNED NOT NULL DEFAULT 1;
ALTER TABLE `tLogout` ADD KEY IF NOT EXISTS `idx_tlogout_user_date` (`userNo`, `logDate` DESC);

-- 18. tReport: sReporterCharID VARCHAR(20), sHarmer VARCHAR(20), sReportType VARCHAR(20), TEXT fields, dDealDate, sOperID
ALTER TABLE `tReport` MODIFY COLUMN `sReporterCharID` VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE `tReport` MODIFY COLUMN `sHarmer` VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE `tReport` MODIFY COLUMN `sReportType` VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE `tReport` MODIFY COLUMN `sReason` TEXT;
ALTER TABLE `tReport` MODIFY COLUMN `sChatLog` MEDIUMTEXT;
ALTER TABLE `tReport` MODIFY COLUMN `sRemark` TEXT;
ALTER TABLE `tReport` ADD COLUMN IF NOT EXISTS `dDealDate` DATETIME DEFAULT NULL;
ALTER TABLE `tReport` ADD COLUMN IF NOT EXISTS `sOperID` VARCHAR(20) DEFAULT NULL;

-- 19. tTutorial: Ensure single table with PK (nUserNo, nCharNo) and index (nUserNo, nState)
ALTER TABLE `tTutorial` ADD KEY IF NOT EXISTS `idx_ttutorial_user_state` (`nUserNo`, `nState`);

-- 20. tUserBlock: index (nUserNo, bStatus, dEndDate)
ALTER TABLE `tUserBlock` ADD KEY IF NOT EXISTS `idx_tuserblock_userno_status_end` (`nUserNo`, `bStatus`, `dEndDate`);

-- 21. Create tAccountItem (22 columns)
CREATE TABLE IF NOT EXISTS `tAccountItem` (
    `goodsNo` INT NOT NULL,
    `billingNo` INT NOT NULL DEFAULT 0,
    `category` VARCHAR(10) NOT NULL,
    `groupNo` INT NOT NULL DEFAULT 0,
    `name` VARCHAR(100) NOT NULL DEFAULT '',
    `imageName` VARCHAR(100) NOT NULL DEFAULT '',
    `itemType` VARCHAR(20) NOT NULL DEFAULT '',
    `unit` INT NOT NULL DEFAULT 1,
    `unitDescription` VARCHAR(20) NOT NULL DEFAULT '',
    `price` INT NOT NULL DEFAULT 0,
    `discount` INT NOT NULL DEFAULT 0,
    `discountStartDate` DATETIME DEFAULT NULL,
    `discountEndDate` DATETIME DEFAULT NULL,
    `isNew` TINYINT(1) NOT NULL DEFAULT 0,
    `isPopular` TINYINT(1) NOT NULL DEFAULT 0,
    `isRecommand` TINYINT(1) NOT NULL DEFAULT 0,
    `isSell` TINYINT(1) NOT NULL DEFAULT 1,
    `sellStartDate` DATETIME DEFAULT NULL,
    `sellEndDate` DATETIME DEFAULT NULL,
    `sellCount` INT NOT NULL DEFAULT 0,
    `registerDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `description` VARCHAR(104) NOT NULL DEFAULT '',
    PRIMARY KEY (`goodsNo`),
    KEY `idx_taccountitem_category` (`category`),
    KEY `idx_taccountitem_group` (`groupNo`),
    KEY `idx_taccountitem_sell` (`isSell`),
    CONSTRAINT `fk_taccountitem_category` FOREIGN KEY (`category`) REFERENCES `tItemCategory` (`category`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 22. Create tAccountDemographics
CREATE TABLE IF NOT EXISTS `tAccountDemographics` (
    `userNo` INT NOT NULL,
    `gender` INT NOT NULL DEFAULT 0,
    `age` INT NOT NULL DEFAULT 0,
    `registerDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`userNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 23. Missing auxiliary / legacy tables
CREATE TABLE IF NOT EXISTS `Auto_Account_Creation` (
    `Allow_Creation` TINYINT(1) NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tPcbang` (
    `pcbangNo` INT NOT NULL AUTO_INCREMENT,
    `sIPStart` VARCHAR(20) NOT NULL,
    `sIPEnd` VARCHAR(20) NOT NULL,
    `sName` VARCHAR(50) NOT NULL DEFAULT '',
    `bIsActive` TINYINT(1) NOT NULL DEFAULT 1,
    PRIMARY KEY (`pcbangNo`),
    KEY `idx_tpcbang_range` (`sIPStart`, `sIPEnd`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tUserKey` (
    `nUserNo` INT NOT NULL,
    `sWebKey` VARCHAR(64) NOT NULL DEFAULT '',
    `dIssueDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dExpireDate` DATETIME DEFAULT NULL,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tUserPlayRemainTime` (
    `nUserNo` INT NOT NULL,
    `nRemainSec` INT NOT NULL DEFAULT 0,
    `bIsTeenager` TINYINT(1) NOT NULL DEFAULT 0,
    `dLastUpdate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tZipcode` (
    `nZipNo` INT NOT NULL AUTO_INCREMENT,
    `sZipCode` VARCHAR(10) NOT NULL DEFAULT '',
    `sSido` VARCHAR(30) NOT NULL DEFAULT '',
    `sGugun` VARCHAR(30) NOT NULL DEFAULT '',
    `sDong` VARCHAR(60) NOT NULL DEFAULT '',
    `sBunji` VARCHAR(40) NOT NULL DEFAULT '',
    PRIMARY KEY (`nZipNo`),
    KEY `idx_tzipcode_zip` (`sZipCode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tbl_user` (
    `user_id` VARCHAR(50) NOT NULL,
    `user_pwd` VARCHAR(50) DEFAULT NULL,
    `user_name` VARCHAR(50) DEFAULT NULL,
    `user_email` VARCHAR(100) DEFAULT NULL,
    `reg_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tbl_user_test` (
    `user_id` VARCHAR(50) NOT NULL,
    `user_pwd` VARCHAR(50) DEFAULT NULL,
    `user_name` VARCHAR(50) DEFAULT NULL,
    `user_email` VARCHAR(100) DEFAULT NULL,
    `reg_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tUserTemp` (
    `nUserNo` INT NOT NULL AUTO_INCREMENT,
    `sUserID` VARCHAR(60) NOT NULL,
    `sUserPW` VARCHAR(32) NOT NULL,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tUserTemp2` (
    `nUserNo` INT NOT NULL AUTO_INCREMENT,
    `sUserID` VARCHAR(60) NOT NULL,
    `sUserPW` VARCHAR(32) NOT NULL,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tUser_ori` (
    `nUserNo` INT NOT NULL,
    `sUserID` VARCHAR(60) NOT NULL,
    `sUserPW` VARCHAR(32) NOT NULL,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `temp_newUser` (
    `sUserID` VARCHAR(60) NOT NULL,
    `sUserPW` VARCHAR(32) NOT NULL,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`sUserID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `last2_onson_users_1` (
    `Column 0` VARCHAR(50) DEFAULT NULL,
    `Column 1` VARCHAR(50) DEFAULT NULL,
    `Column 2` VARCHAR(200) DEFAULT NULL,
    `Column 3` VARCHAR(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `last2_onson_users_2` (
    `Column 0` VARCHAR(50) DEFAULT NULL,
    `Column 1` VARCHAR(50) DEFAULT NULL,
    `Column 2` VARCHAR(200) DEFAULT NULL,
    `Column 3` VARCHAR(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `login_issue_check` (
    `userid` INT DEFAULT NULL,
    `username` VARCHAR(255) DEFAULT NULL,
    `usedpwd` VARCHAR(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `onson_0419` (
    `userid` VARCHAR(50) DEFAULT NULL,
    `sparkid` VARCHAR(50) DEFAULT NULL,
    `email` VARCHAR(200) DEFAULT NULL,
    `password` VARCHAR(200) DEFAULT NULL,
    `bisBlock` VARCHAR(50) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `onson_users_0314` (
    `Column 0` VARCHAR(50) DEFAULT NULL,
    `Column 1` VARCHAR(200) DEFAULT NULL,
    `Column 2` VARCHAR(200) DEFAULT NULL,
    `Column 3` VARCHAR(200) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tuseroption_2` (
    `nUserNo` INT NOT NULL,
    `sEmail` VARCHAR(256) NOT NULL DEFAULT '',
    `sJuminNo` VARCHAR(50) DEFAULT NULL,
    `sPhoneNo` VARCHAR(20) DEFAULT NULL,
    `nZipNo` INT DEFAULT NULL,
    `sAddress` VARCHAR(90) DEFAULT NULL,
    `bIsMail` TINYINT(1) NOT NULL DEFAULT 0,
    `sChannel` VARCHAR(20) DEFAULT NULL,
    `nAffiliateID` VARCHAR(50) DEFAULT NULL,
    `szAffiliateSubID` VARCHAR(50) DEFAULT NULL,
    `szCountryCode` VARCHAR(50) DEFAULT NULL,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tEvent_L20` (
    `nRowNo` INT NOT NULL AUTO_INCREMENT,
    `nUserNo` INT NOT NULL,
    `sUserID` VARCHAR(60) NOT NULL DEFAULT '',
    `nWorldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `bRewarded` TINYINT(1) NOT NULL DEFAULT 0,
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nRowNo`),
    KEY `idx_event_l20_user` (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tHackScanLog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nWorldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nHackType` SMALLINT NOT NULL DEFAULT 0,
    `sDetail` VARCHAR(64) NOT NULL DEFAULT '',
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_thackscanlog_user` (`nUserNo`, `dDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tSyslog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `nOpCode` SMALLINT UNSIGNED NOT NULL DEFAULT 0,
    `nUserNo` INT NOT NULL DEFAULT 0,
    `sUserID` VARCHAR(60) NOT NULL DEFAULT '',
    `sIP` VARCHAR(30) NOT NULL DEFAULT '',
    `nWorldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nCCU` INT NOT NULL DEFAULT 0,
    `nErrorCode` SMALLINT NOT NULL DEFAULT 0,
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_tsyslog_opcode_date` (`nOpCode`, `dDate`),
    KEY `idx_tsyslog_user` (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tWorldServer` (
    `nWorldNo` TINYINT UNSIGNED NOT NULL,
    `sWorldName` VARCHAR(16) NOT NULL DEFAULT '',
    `sWorldIP` VARCHAR(16) NOT NULL DEFAULT '',
    `nWorldPort` SMALLINT UNSIGNED NOT NULL DEFAULT 9014,
    `nMaxUser` SMALLINT UNSIGNED NOT NULL DEFAULT 1500,
    `nUserRatable` TINYINT UNSIGNED NOT NULL DEFAULT 100,
    `bIsMaintenance` TINYINT(1) NOT NULL DEFAULT 0,
    `nCurrentCCU` INT NOT NULL DEFAULT 0,
    `dLastHeartbeat` DATETIME DEFAULT NULL,
    PRIMARY KEY (`nWorldNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `tActiveAccountSession` (
    `nUserNo` INT NOT NULL,
    `sUserID` VARCHAR(60) NOT NULL DEFAULT '',
    `nWorldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sClientIP` VARCHAR(30) NOT NULL DEFAULT '',
    `abValidateToken` BINARY(64) DEFAULT NULL,
    `sOTP` CHAR(32) DEFAULT NULL,
    `dTokenExpire` DATETIME DEFAULT NULL,
    `dOTPExpire` DATETIME DEFAULT NULL,
    `nState` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dLoginTime` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `dUpdatedAt` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`),
    KEY `idx_session_otp` (`sOTP`),
    KEY `idx_session_world` (`nWorldNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 24. Create or replace Views
CREATE OR REPLACE VIEW `VT_ITEM_SELL` AS
SELECT *
FROM `tAccountItem`
WHERE `isSell` = 1
  AND (`sellStartDate` IS NULL OR `sellStartDate` <= NOW())
  AND (`sellEndDate` IS NULL OR `sellEndDate` >= NOW());

CREATE OR REPLACE VIEW `VT_ITEM_GROUP` AS
SELECT *
FROM `tAccountItem`;

CREATE OR REPLACE VIEW `VT_PAY` AS
SELECT 
    p.orderNo,
    p.userNo,
    p.payNo,
    p.cash,
    p.userIP,
    p.registerDate,
    pr.result,
    pr.resultDate,
    pr.pgTID
FROM `tPay` p
JOIN `tPayResult` pr ON p.orderNo = pr.orderNo;
