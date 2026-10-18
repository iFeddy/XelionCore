-- ============================================================================
-- XELION ONLINE - ESQUEMA DE BASE DE DATOS OFICIAL (MYSQL)
-- Core Tables para LoginService (Account & AccountLog)
-- Fidelidad binaria y estructural con AstraServer / Fiesta Online Legacy
-- ============================================================================

CREATE DATABASE IF NOT EXISTS `xelion_online`
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE `xelion_online`;

-- ----------------------------------------------------------------------------
-- 1. tUserAuth: Roles y permisos de acceso
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tUserAuth` (
    `nAuthID` TINYINT UNSIGNED NOT NULL,
    `sAuthName` VARCHAR(20) NOT NULL,
    `bIsLoginAble` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    PRIMARY KEY (`nAuthID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 2. tUser: Cuentas de usuario y credenciales de autenticación
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tUser` (
    `nUserNo` INT NOT NULL AUTO_INCREMENT,
    `sUserID` VARCHAR(60) NOT NULL,
    `sUserPW` VARCHAR(32) NOT NULL,
    `sUserName` VARCHAR(60) NOT NULL DEFAULT '',
    `bIsBlock` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `bIsLoginAble` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    `bIsDelete` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nAuthID` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    `sUserIP` VARCHAR(30) NOT NULL DEFAULT '',
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`),
    UNIQUE KEY `uk_tuser_userid` (`sUserID`),
    KEY `idx_tuser_auth` (`nAuthID`),
    CONSTRAINT `fk_tuser_auth` FOREIGN KEY (`nAuthID`) REFERENCES `tUserAuth` (`nAuthID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 3. tUserBlock: Sanciones y bloqueos temporales/permanentes de cuentas
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tUserBlock` (
    `nBlockNo` INT NOT NULL AUTO_INCREMENT,
    `nUserNo` INT NOT NULL,
    `nType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dBlockDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `nTimeValue` INT NOT NULL DEFAULT 0,
    `dEndDate` DATETIME DEFAULT NULL,
    `bStatus` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    `sComment` VARCHAR(100) NOT NULL DEFAULT '',
    PRIMARY KEY (`nBlockNo`),
    KEY `idx_tuserblock_userno_status` (`nUserNo`, `bStatus`),
    CONSTRAINT `fk_tuserblock_user` FOREIGN KEY (`nUserNo`) REFERENCES `tUser` (`nUserNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 4. tAccountLog: Auditoría de intentos de inicio de sesión y errores
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tAccountLog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `sUserID` VARCHAR(60) NOT NULL,
    `nType` TINYINT UNSIGNED NOT NULL, -- 1: Login Success, 2: Login Fail
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nPlayMin` SMALLINT NOT NULL DEFAULT 0,
    `sUserPassword` VARCHAR(32) NOT NULL DEFAULT '',
    `nWorldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nErrorCode` SMALLINT NOT NULL DEFAULT 0,
    `sIP` VARCHAR(30) NOT NULL DEFAULT '',
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_taccountlog_userno` (`nUserNo`),
    KEY `idx_taccountlog_date` (`dDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 5. tIPBlock: Firewall de rangos IP bloqueados
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tIPBlock` (
    `nBlockNo` INT NOT NULL AUTO_INCREMENT,
    `sIPStart` VARCHAR(30) NOT NULL,
    `sIPEnd` VARCHAR(30) NOT NULL,
    `sDesc` VARCHAR(100) NOT NULL DEFAULT '',
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nBlockNo`),
    KEY `idx_tipblock_range` (`sIPStart`, `sIPEnd`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 6. tLogout: Registro de desconexión y tiempo de juego
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tLogout` (
    `rowNo` BIGINT NOT NULL AUTO_INCREMENT,
    `logDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `worldNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `userNo` INT NOT NULL,
    `charNo` INT NOT NULL DEFAULT 0,
    `class` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `level` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    `experience` BIGINT NOT NULL DEFAULT 0,
    `playTime` SMALLINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`rowNo`),
    KEY `idx_tlogout_userno` (`userNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 7. tTransferToken: Sesiones volátiles de transferencia (100% en RAM, 0 disco)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tTransferToken` (
    `sToken` VARCHAR(64) NOT NULL,
    `nUserNo` INT NOT NULL,
    `sUserID` VARCHAR(60) NOT NULL,
    `nAuthID` TINYINT UNSIGNED NOT NULL DEFAULT 1,
    `sIP` VARCHAR(30) NOT NULL DEFAULT '',
    `dCreatedAt` DATETIME NOT NULL,
    PRIMARY KEY (`sToken`),
    KEY `idx_transfer_created` (`dCreatedAt`)
) ENGINE=MEMORY DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- WORLD & CHARACTER TABLES (World00_Character)
-- Fidelidad binaria con AstraServer / WorldManager.exe & Character.exe
-- ============================================================================

-- ----------------------------------------------------------------------------
-- 8. tCharacter: Maestro de personajes del mundo
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacter` (
    `nCharNo` INT NOT NULL AUTO_INCREMENT,
    `sID` VARCHAR(40) NOT NULL,
    `nUserNo` INT NOT NULL,
    `nAdminLevel` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nPS` INT NOT NULL DEFAULT 0,
    `nGS` INT NOT NULL DEFAULT 0,
    `dCreateDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `bDeleted` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nLoginZoneX` INT NOT NULL DEFAULT 0,
    `nLoginZoneY` INT NOT NULL DEFAULT 0,
    `nLoginZoneD` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nPKYellowTime` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDeletedDate` DATETIME DEFAULT NULL,
    `nAP` SMALLINT NOT NULL DEFAULT 0,
    `sLoginZone` VARCHAR(16) NOT NULL DEFAULT 'RouN',
    `nPKCount` INT NOT NULL DEFAULT 0,
    `nLevel` SMALLINT NOT NULL DEFAULT 1,
    `nExp` BIGINT NOT NULL DEFAULT 0,
    `nSlotNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nHP` SMALLINT NOT NULL DEFAULT 100,
    `nSP` SMALLINT NOT NULL DEFAULT 100,
    `nHPS` SMALLINT NOT NULL DEFAULT 100,
    `nSPS` SMALLINT NOT NULL DEFAULT 100,
    `dLastLoginDate` DATETIME DEFAULT NULL,
    `nLoginCount` INT NOT NULL DEFAULT 0,
    `nCreateWorld` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nFame` INT NOT NULL DEFAULT 0,
    `nPrisonMin` SMALLINT NOT NULL DEFAULT 0,
    `nPlayMin` INT NOT NULL DEFAULT 0,
    `nFlags` INT NOT NULL DEFAULT 0,
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    `nStrength` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nConstitute` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nDexterity` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nIntelligence` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nMentalPower` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nRedistributePoint` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nKQHandle` INT NOT NULL DEFAULT 0,
    `sKQMap` VARCHAR(16) NOT NULL DEFAULT '',
    `nKQX` INT NOT NULL DEFAULT 0,
    `nKQY` INT NOT NULL DEFAULT 0,
    `nKQD` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dKQDate` DATETIME DEFAULT NULL,
    `nFriendPoint` INT NOT NULL DEFAULT 0,
    `nLastFriendPoint` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`),
    KEY `idx_tcharacter_userno_slot` (`nUserNo`, `nSlotNo`, `bDeleted`),
    KEY `idx_tcharacter_id` (`sID`),
    CONSTRAINT `fk_tcharacter_user` FOREIGN KEY (`nUserNo`) REFERENCES `tUser` (`nUserNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 9. tCharacterShape: Rasgos visuales 3D y apariencia del personaje
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterShape` (
    `nCharNo` INT NOT NULL,
    `nClass` TINYINT UNSIGNED NOT NULL,
    `nRace` TINYINT UNSIGNED NOT NULL,
    `nGender` TINYINT UNSIGNED NOT NULL,
    `nHairType` TINYINT UNSIGNED NOT NULL,
    `nHairColor` TINYINT UNSIGNED NOT NULL,
    `nFaceShape` TINYINT UNSIGNED NOT NULL,
    PRIMARY KEY (`nCharNo`),
    CONSTRAINT `fk_tcharactershape_char` FOREIGN KEY (`nCharNo`) REFERENCES `tCharacter` (`nCharNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 10. tCharacterOptions: Preferencias de cliente (Video, Audio, UI, Teclas)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterOptions` (
    `nCharNo` INT NOT NULL,
    `sShortCutData` VARBINARY(2048) DEFAULT NULL,
    `sShortCutSize` VARBINARY(24) DEFAULT NULL,
    `sVideo` VARBINARY(60) DEFAULT NULL,
    `sSound` VARBINARY(1) DEFAULT NULL,
    `sGame` VARBINARY(64) DEFAULT NULL,
    `sWindowsPos` VARBINARY(136) DEFAULT NULL,
    `sKeyMapping` VARBINARY(52) DEFAULT NULL,
    PRIMARY KEY (`nCharNo`),
    CONSTRAINT `fk_tcharacteroptions_char` FOREIGN KEY (`nCharNo`) REFERENCES `tCharacter` (`nCharNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 11. tCharacterShortCut: Barra de accesos directos de habilidades y acciones
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterShortCut` (
    `nCharNo` INT NOT NULL,
    `nSlotNo` TINYINT UNSIGNED NOT NULL,
    `sys_create_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `nCodeType` SMALLINT NOT NULL,
    `nCodeNo` SMALLINT NOT NULL,
    `nValue` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nSlotNo`),
    CONSTRAINT `fk_tcharactershortcut_char` FOREIGN KEY (`nCharNo`) REFERENCES `tCharacter` (`nCharNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- 12. tItem: Inventario y equipamiento del personaje
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItem` (
    `nItemKey` BIGINT NOT NULL AUTO_INCREMENT,
    `nStorageType` TINYINT UNSIGNED NOT NULL DEFAULT 0, -- 0: Equip, 1: Inv, etc.
    `nOwner` INT NOT NULL,
    `nStorage` SMALLINT NOT NULL,
    `nItemID` INT NOT NULL,
    `nFlags` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nItemKey`),
    KEY `idx_titem_owner_storage` (`nOwner`, `nStorageType`, `nStorage`),
    CONSTRAINT `fk_titem_char` FOREIGN KEY (`nOwner`) REFERENCES `tCharacter` (`nCharNo`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


-- ============================================================================
-- TABLAS OFICIALES DE GAMEPLAY (WORLD00_CHARACTER / ZONE)
-- ============================================================================

-- ----------------------------------------------------------------------------
-- Table: tArenaChallenge
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tArenaChallenge` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nMonsterLevel` SMALLINT NOT NULL DEFAULT 0,
    `nMonsterMKC` SMALLINT NOT NULL DEFAULT 0,
    `nMonsterKPS` SMALLINT NOT NULL DEFAULT 0,
    `nMonsterKills` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nMonsterLevel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tArenaFBZ
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tArenaFBZ` (
    `nFBZWin` INT NOT NULL DEFAULT 0,
    `nFBZLose` INT NOT NULL DEFAULT 0,
    `nFBZLadder` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tArenaGuild
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tArenaGuild` (
    `nGuildLose` INT NOT NULL DEFAULT 0,
    `nGuildDraw` INT NOT NULL DEFAULT 0,
    `nGuildNo` INT NOT NULL DEFAULT 0,
    `nGuildWin` INT NOT NULL DEFAULT 0,
    `nGuildLadder` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nGuildNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tArenaPvP
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tArenaPvP` (
    `nPvPLadder` INT NOT NULL DEFAULT 0,
    `nPvPWin` INT NOT NULL DEFAULT 0,
    `nPvPLose` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nPvPDraw` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tAttendance
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tAttendance` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dAttendanceDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `dAttendanceDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tAttendanceReward
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tAttendanceReward` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dEventDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dRewardDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nRewardID` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `dEventDate`, `nRewardID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tAutoMon
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tAutoMon` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSlot` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nItemID` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nType`, `nSlot`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tAutoPick
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tAutoPick` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `IsAutoPick` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterChangeID
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterChangeID` (
    `nRowNo` INT NOT NULL AUTO_INCREMENT,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `bFlag` TINYINT(1) NOT NULL DEFAULT 0,
    `bInit` TINYINT(1) NOT NULL DEFAULT 0,
    `sOldID` VARCHAR(255) NOT NULL DEFAULT '',
    `sNewID` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nRowNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterChargedBuffer
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterChargedBuffer` (
    `nNo` INT NOT NULL AUTO_INCREMENT,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nKey` INT NOT NULL DEFAULT 0,
    `nID` SMALLINT NOT NULL DEFAULT 0,
    `dUseDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dEndDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterGameOption
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterGameOption` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nOptionNo` SMALLINT NOT NULL DEFAULT 0,
    `nValue` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nOptionNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterKeyMap
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterKeyMap` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nFunctionNo` SMALLINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nExtendedKey` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nASCIICode` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nFunctionNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterLoading
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterLoading` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nSeqNo` INT NOT NULL DEFAULT 0,
    `dWorldTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dDBTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sLoadingMapName` VARCHAR(255) NOT NULL DEFAULT '',
    `nFailType` SMALLINT NOT NULL DEFAULT 0,
    `nLoadingFailCnt` SMALLINT NOT NULL DEFAULT 0,
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterSingleOption
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterSingleOption` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `bMysteryBolt` TINYINT(1) NOT NULL DEFAULT 0,
    `bAutoPick` TINYINT(1) NOT NULL DEFAULT 0,
    `dRestExpLastExecuteTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nChatBan` INT NOT NULL DEFAULT 0,
    `nAutoBuff` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nAutoChargedBuff` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterState
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterState` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nStateNo` SMALLINT NOT NULL DEFAULT 0,
    `nStateLevel` SMALLINT NOT NULL DEFAULT 0,
    `nStateTime` INT NOT NULL DEFAULT 0,
    `nStateCount` SMALLINT NOT NULL DEFAULT 0,
    `nStatePowerDemage` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nStatePowerSP` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nStatePowerKeepTime` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nStatePowerCoolTime` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSetItemIndexA` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemValueA` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemIndexB` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemValueB` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemIndexC` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemValueC` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemIndexD` SMALLINT NOT NULL DEFAULT 0,
    `nSetItemValueD` SMALLINT NOT NULL DEFAULT 0,
    `nRestTime` INT NOT NULL DEFAULT 0,
    `nCasterLevel` SMALLINT NOT NULL DEFAULT 0,
    `nCasterAttackPower` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nStateNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCharacterTitle
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCharacterTitle` (
    `nNo` BIGINT NOT NULL AUTO_INCREMENT,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nElementNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nElementValue` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tChatBlock
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tChatBlock` (
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`, `nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tChatColor
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tChatColor` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nFontColorID` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nBalloonColorID` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCollectionCard
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCollectionCard` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nCardID` INT NOT NULL DEFAULT 0,
    `nSerialNo` BIGINT NOT NULL DEFAULT 0,
    `nStarNumber` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nCardID`, `nSerialNo`),
    KEY `idx_card` (`nCharNo`, `nCardID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCollectionFavorites
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCollectionFavorites` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nSlotNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nCardID` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nSlotNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCollectionReward
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCollectionReward` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nReward` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nReward`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tCollectionSerialCount
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tCollectionSerialCount` (
    `nCardID` INT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nCount` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCardID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tDice
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tDice` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `sDice` VARCHAR(255) NOT NULL DEFAULT '',
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nGameNo`, `nCharNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tDiceLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tDiceLog` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `nCustomCharNo` INT NOT NULL DEFAULT 0,
    `nSlotNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nBattingMoney` BIGINT NOT NULL DEFAULT 0,
    `nBattingRate` INT NOT NULL DEFAULT 0,
    `bStatus` TINYINT(1) NOT NULL DEFAULT 0,
    `nTaxRate` INT NOT NULL DEFAULT 0,
    `nDividend` BIGINT NOT NULL DEFAULT 0,
    KEY `idx_game` (`nGameNo`),
    KEY `idx_char` (`nCustomCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tDicePayLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tDicePayLog` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nTotalDividend` BIGINT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_char` (`nCharNo`),
    KEY `idx_item` (`nItemKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tDiceResult
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tDiceResult` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `nCustomCharNo` INT NOT NULL DEFAULT 0,
    `nBeforeMoney` BIGINT NOT NULL DEFAULT 0,
    `nTotalDividend` BIGINT NOT NULL DEFAULT 0,
    `nAfterMoney` BIGINT NOT NULL DEFAULT 0,
    KEY `idx_game` (`nGameNo`),
    KEY `idx_char` (`nCustomCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tExp
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tExp` (
    `level` INT NOT NULL DEFAULT 0,
    `NextExp` BIGINT NOT NULL DEFAULT 0,
    `sExp` BIGINT NOT NULL DEFAULT 0,
    `eExp` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`level`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tFriend
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tFriend` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nFriendCharNo` INT NOT NULL DEFAULT 0,
    `nFlags` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nFriendCharNo`),
    KEY `idx_friend` (`nFriendCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGMItemReward
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGMItemReward` (
    `nGroupNo` BIGINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sOperID` VARCHAR(255) NOT NULL DEFAULT '',
    `sIP` VARCHAR(255) NOT NULL DEFAULT '',
    `sCause` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nGroupNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGMItemRewardLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGMItemRewardLog` (
    `nGroupNo` BIGINT NOT NULL DEFAULT 0,
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nStorageType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nItemID` INT NOT NULL DEFAULT 0,
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_group` (`nGroupNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleCoin
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleCoin` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nCoin` BIGINT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleCoinExchangeLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleCoinExchangeLog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nExchangedCoin` BIGINT NOT NULL DEFAULT 0,
    `nStatus` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nBeforeCoin` BIGINT NOT NULL DEFAULT 0,
    `nAfterCoin` BIGINT NOT NULL DEFAULT 0,
    `nBeforeMoney` BIGINT NOT NULL DEFAULT 0,
    `nAfterMoney` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleDice
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleDice` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `sDice` VARCHAR(255) NOT NULL DEFAULT '',
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nGameNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleDiceLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleDiceLog` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nSlotNo` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nBattingMoney` BIGINT NOT NULL DEFAULT 0,
    `nBattingRate` INT NOT NULL DEFAULT 0,
    `bStatus` TINYINT(1) NOT NULL DEFAULT 0,
    `nTaxRate` INT NOT NULL DEFAULT 0,
    `nDividend` BIGINT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_game` (`nGameNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleDiceResult
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleDiceResult` (
    `nGameNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nBeforeMoney` BIGINT NOT NULL DEFAULT 0,
    `nTotalDividend` BIGINT NOT NULL DEFAULT 0,
    `nAfterMoney` BIGINT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    KEY `idx_game` (`nGameNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleJackpot
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleJackpot` (
    `sGameName` VARCHAR(255) NOT NULL DEFAULT '',
    `nBattingUnit` BIGINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nJackpotCoin` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`sGameName`, `nBattingUnit`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleJackpotLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleJackpotLog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `sGameName` VARCHAR(255) NOT NULL DEFAULT '',
    `nBattingUnit` BIGINT NOT NULL DEFAULT 0,
    `sID` VARCHAR(255) NOT NULL DEFAULT '',
    `nJackpotCoin` BIGINT NOT NULL DEFAULT 0,
    `nGameNo` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nLogNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGambleSlotMachineLog
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGambleSlotMachineLog` (
    `nGameNo` BIGINT NOT NULL AUTO_INCREMENT,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nBattingUnit` BIGINT NOT NULL DEFAULT 0,
    `nBatting` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sGameResult` VARCHAR(255) NOT NULL DEFAULT '',
    `bIsWin` TINYINT(1) NOT NULL DEFAULT 0,
    `nPrizeCoin` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nGameNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGameOptionInit
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGameOptionInit` (
    `nOptionNo` SMALLINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nDefaultValue` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `bIsActive` TINYINT(1) NOT NULL DEFAULT 0,
    `sDescription` VARCHAR(255) NOT NULL DEFAULT '',
    `sDescriptionLocal` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nOptionNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuild
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuild` (
    `nNo` INT NOT NULL AUTO_INCREMENT,
    `sName` VARCHAR(255) NOT NULL DEFAULT '',
    `sPassword` VARCHAR(255) NOT NULL DEFAULT '',
    `nMaxMembers` SMALLINT NOT NULL DEFAULT 0,
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    `nType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nGrade` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nFame` INT NOT NULL DEFAULT 0,
    `nStoneLevel` SMALLINT NOT NULL DEFAULT 0,
    `nExp` BIGINT NOT NULL DEFAULT 0,
    `dCreateDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dWarRequestDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dWarStartDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dWarEndDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nWarEnemy` INT NOT NULL DEFAULT 0,
    `nWarWinCount` INT NOT NULL DEFAULT 0,
    `nWarLoseCount` INT NOT NULL DEFAULT 0,
    `nWarDrawCount` INT NOT NULL DEFAULT 0,
    `nDismissStatus` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDismissDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sIntro` VARCHAR(255) NOT NULL DEFAULT '',
    `dNotifyDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sNotifyCharID` VARCHAR(255) NOT NULL DEFAULT '',
    `sNotify` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nNo`),
    UNIQUE KEY `idx_name` (`sName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademy
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademy` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nMasterCharNo` INT NOT NULL DEFAULT 0,
    `nMaxMember` SMALLINT NOT NULL DEFAULT 0,
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    `nPoint` INT NOT NULL DEFAULT 0,
    `dBuffEndTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dNotifyDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sNotify` VARCHAR(255) NOT NULL DEFAULT '',
    `sIntro` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nNo`),
    KEY `idx_master` (`nMasterCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyGraduate
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyGraduate` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sID` VARCHAR(255) NOT NULL DEFAULT '',
    `nClass` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nJoinLevel` SMALLINT NOT NULL DEFAULT 0,
    `dJoinDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nGraduateLevel` SMALLINT NOT NULL DEFAULT 0,
    `dGraduateDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nNo`, `nCharNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyHistory
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyHistory` (
    `nAcademyID` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sID` VARCHAR(255) NOT NULL DEFAULT '',
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nAction` INT NOT NULL DEFAULT 0,
    KEY `idx_academy` (`nAcademyID`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyMember
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyMember` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nJoinLevel` SMALLINT NOT NULL DEFAULT 0,
    `dJoinDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nGrade` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sIntro` VARCHAR(255) NOT NULL DEFAULT '',
    `nChatBanTime` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nNo`, `nCharNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyRanking
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyRanking` (
    `nRowNo` INT NOT NULL AUTO_INCREMENT,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nRanking` INT NOT NULL DEFAULT 0,
    `nGuildNo` INT NOT NULL DEFAULT 0,
    `sGuildName` VARCHAR(255) NOT NULL DEFAULT '',
    `nMasterNo` INT NOT NULL DEFAULT 0,
    `sMasterID` VARCHAR(255) NOT NULL DEFAULT '',
    `nPoint` INT NOT NULL DEFAULT 0,
    `nCount` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nRowNo`),
    KEY `idx_ranking` (`dDate`, `nRanking`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyReward
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyReward` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nLevel` SMALLINT NOT NULL DEFAULT 0,
    `nTypeID` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nValue` BIGINT NOT NULL DEFAULT 0,
    `nItemID` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nNo`, `nLevel`, `nTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildAcademyRewardType
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildAcademyRewardType` (
    `nTypeID` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sTypeName` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nTypeID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildEmblem
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildEmblem` (
    `nGuildNo` INT NOT NULL DEFAULT 0,
    `nColor` INT NOT NULL DEFAULT 0,
    `nIcon` INT NOT NULL DEFAULT 0,
    `dCreated` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dUpdated` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nUnk` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nGuildNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildMember
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildMember` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dJoinDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nGrade` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nFlags` INT NOT NULL DEFAULT 0,
    `nExpRatio` SMALLINT NOT NULL DEFAULT 0,
    `sIntro` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nNo`, `nCharNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildToken
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildToken` (
    `nNo` INT NOT NULL DEFAULT 0,
    `nToken` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildTournament
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildTournament` (
    `nGTNo` INT NOT NULL DEFAULT 0,
    `nType` SMALLINT NOT NULL DEFAULT 0,
    `nData` INT NOT NULL DEFAULT 0,
    `nGuildNo` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nGTNo`, `nGuildNo`),
    KEY `idx_guild` (`nGuildNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGuildTournamentRegist
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGuildTournamentRegist` (
    `nGTNo` INT NOT NULL DEFAULT 0,
    `nNo` INT NOT NULL DEFAULT 0,
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nSeed` INT NOT NULL DEFAULT 0,
    `bChoice` TINYINT(1) NOT NULL DEFAULT 0,
    PRIMARY KEY (`nGTNo`, `nNo`),
    KEY `idx_guild` (`nNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tHolyPromise
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tHolyPromise` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nUpCharNo` INT NOT NULL DEFAULT 0,
    `nFlags` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dRegDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dRejoinableDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dMemberAcceptableDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nUpCharNo`),
    KEY `idx_upchar` (`nUpCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tHolyPromiseReward
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tHolyPromiseReward` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nSeqNo` INT NOT NULL DEFAULT 0,
    `nType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    `nRemainMoney` BIGINT NOT NULL DEFAULT 0,
    `nReceiptCharNo` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nSeqNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tInGameGuide
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tInGameGuide` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nGuideNo` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nGuideNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemActionCoolTime
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemActionCoolTime` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nActionGroup` SMALLINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dEndTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nActionGroup`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemEnchant
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemEnchant` (
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nSlot` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nGemItemID` INT NOT NULL DEFAULT 0,
    `nRestCount` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nItemKey`, `nSlot`),
    KEY `idx_item` (`nItemKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemMobList
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemMobList` (
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nMobID` SMALLINT NOT NULL DEFAULT 0,
    `nGrade` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nKillCount` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nItemKey`, `nMobID`),
    KEY `idx_item` (`nItemKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemOptions
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemOptions` (
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nOptionType` SMALLINT NOT NULL DEFAULT 0,
    `nOptionData` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nItemKey`, `nOptionType`),
    KEY `idx_item` (`nItemKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemStorageFindEmpty
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemStorageFindEmpty` (
    `nStorage` INT NOT NULL DEFAULT 0,
    KEY `idx_storage` (`nStorage`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tItemTitle
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tItemTitle` (
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nTitleType` SMALLINT NOT NULL DEFAULT 0,
    `sTitle` VARCHAR(255) NOT NULL DEFAULT '',
    `nTitleCharNo` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nItemKey`, `nTitleType`),
    KEY `idx_item` (`nItemKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tKeyMapInit
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tKeyMapInit` (
    `nFunctionNo` SMALLINT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nExtendedKey` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nASCIICode` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sDescription` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nFunctionNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tKingdomQuest
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tKingdomQuest` (
    `nQuestNo` SMALLINT NOT NULL DEFAULT 0,
    `sData` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nQuestNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tLog_GuildMasterChange
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tLog_GuildMasterChange` (
    `nRowNo` BIGINT NOT NULL AUTO_INCREMENT,
    `nGuildNo` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nOldGuildMaster` INT NOT NULL DEFAULT 0,
    `nNewGuildMaster` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nRowNo`),
    KEY `idx_guild` (`nGuildNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tMarriage
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tMarriage` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nPartnerCharNo` INT NOT NULL DEFAULT 0,
    `nStatus` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSubStatus` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dWhen` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dSubWhen` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nMarriageCount` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nDivorceCancel` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`),
    KEY `idx_partner` (`nPartnerCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tMasterPupils
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tMasterPupils` (
    `nPupils` INT NOT NULL DEFAULT 0,
    `nMaster` INT NOT NULL DEFAULT 0,
    `nFlags` INT NOT NULL DEFAULT 0,
    `dRegDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nPupils`, `nMaster`),
    KEY `idx_master` (`nMaster`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tMinihouse
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tMinihouse` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nVisitor` INT NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tMinihouseGate
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tMinihouseGate` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nItemID` INT NOT NULL DEFAULT 0,
    `bType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `dEndDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nItemID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tMinihouseGuest
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tMinihouseGuest` (
    `nArticleNo` INT NOT NULL AUTO_INCREMENT,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nWriterNo` INT NOT NULL DEFAULT 0,
    `sWriterID` VARCHAR(255) NOT NULL DEFAULT '',
    `sContent` VARCHAR(255) NOT NULL DEFAULT '',
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nArticleNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tParty
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tParty` (
    `nFlags` INT NOT NULL DEFAULT 0,
    `nMasterCharNo` INT NOT NULL DEFAULT 0,
    `nPartyNo` INT NOT NULL AUTO_INCREMENT,
    PRIMARY KEY (`nPartyNo`),
    KEY `idx_master` (`nMasterCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tPartyMember
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tPartyMember` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nFlag` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nPartyNo` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nPartyNo`, `nCharNo`),
    KEY `idx_char` (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tPet
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tPet` (
    `nStorage` SMALLINT NOT NULL DEFAULT 0,
    `nPetNumber` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nOwner` BIGINT NOT NULL DEFAULT 0,
    `nPetKey` BIGINT NOT NULL DEFAULT 0,
    `nPetID` INT NOT NULL DEFAULT 0,
    `sInfo` VARCHAR(255) NOT NULL DEFAULT '',
    `nFlags` INT NOT NULL DEFAULT 0,
    `nAssignName` VARCHAR(255) NOT NULL DEFAULT '',
    `nLevel` SMALLINT NOT NULL DEFAULT 0,
    `nExp` BIGINT NOT NULL DEFAULT 0,
    `nSatiety` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSpeed` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nAbilityType` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nAbilityValue` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkinColor` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nStorage`, `nPetNumber`, `nOwner`),
    KEY `idx_owner` (`nOwner`),
    KEY `idx_petkey` (`nPetKey`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tPetNew
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tPetNew` (
    `nPetNo` INT NOT NULL AUTO_INCREMENT,
    `nOwnerNo` INT NOT NULL DEFAULT 0,
    `sPetName` VARCHAR(255) NOT NULL DEFAULT '',
    `nPetMind` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nPetStress` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nPetNo`),
    KEY `idx_owner` (`nOwnerNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tPetSkill
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tPetSkill` (
    `nSkillLevel` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkillNo` INT NOT NULL DEFAULT 0,
    `nSkillExp` INT NOT NULL DEFAULT 0,
    `nPetKey` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nPetKey`, `nSkillNo`, `nSkillLevel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tPrison
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tPrison` (
    `nPrisonNum` BIGINT NOT NULL AUTO_INCREMENT,
    `nCharID` INT NOT NULL DEFAULT 0,
    `dPrisonDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nMinute` SMALLINT NOT NULL DEFAULT 0,
    `sReason` VARCHAR(255) NOT NULL DEFAULT '',
    `nRegNumber` BIGINT NOT NULL DEFAULT 0,
    `sGmID` VARCHAR(255) NOT NULL DEFAULT '',
    `sRemark` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nPrisonNum`),
    KEY `idx_char` (`nCharID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tQuest
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tQuest` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nQuestNo` INT NOT NULL DEFAULT 0,
    `nStatus` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `sData` VARCHAR(255) NOT NULL DEFAULT '',
    PRIMARY KEY (`nCharNo`, `nQuestNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tQuestTimes
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tQuestTimes` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nQuestNo` INT NOT NULL DEFAULT 0,
    `nTimes` INT NOT NULL DEFAULT 0,
    `dLastComplete` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`, `nQuestNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tRestExp
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tRestExp` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `dRestExpLastExcuteTime` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tSkill
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tSkill` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nSkillNo` SMALLINT NOT NULL DEFAULT 0,
    `nSkillLevel` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkillExp` INT NOT NULL DEFAULT 0,
    `nSkillWriteTime` BIGINT NOT NULL DEFAULT 0,
    `nSkillCoolTime` INT NOT NULL DEFAULT 0,
    `nSkillPowerDemage` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkillPowerSP` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkillPowerKeepTime` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nSkillPowerCoolTime` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`, `nSkillNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tSpamBlockTime
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tSpamBlockTime` (
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `nBlockTime` INT NOT NULL DEFAULT 0,
    `nStandbyTime` INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tSpamFilter
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tSpamFilter` (
    `nUserNo` INT NOT NULL DEFAULT 0,
    `sys_create_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sys_update_date` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `bIsSpamSpread` TINYINT(1) NOT NULL DEFAULT 0,
    `bIsSpamBlock` TINYINT(1) NOT NULL DEFAULT 0,
    `nSpamFilterRate` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tStatCharLevel
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tStatCharLevel` (
    `dDate` VARCHAR(255) NOT NULL DEFAULT '',
    `nLevel` SMALLINT NOT NULL DEFAULT 0,
    `nTotalBefore` INT NOT NULL DEFAULT 0,
    `nBefore` INT NOT NULL DEFAULT 0,
    `nTotalAfter` INT NOT NULL DEFAULT 0,
    `nAfter` INT NOT NULL DEFAULT 0,
    KEY `idx_stat` (`dDate`, `nLevel`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tStatGameMoney
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tStatGameMoney` (
    `dDate` VARCHAR(255) NOT NULL DEFAULT '',
    `nLevel030` BIGINT NOT NULL DEFAULT 0,
    `nLevel060` BIGINT NOT NULL DEFAULT 0,
    `nLevel110` BIGINT NOT NULL DEFAULT 0,
    `nLevel120` BIGINT NOT NULL DEFAULT 0,
    `nUserMoney` BIGINT NOT NULL DEFAULT 0,
    `nGuildMoney` BIGINT NOT NULL DEFAULT 0,
    KEY `idx_date` (`dDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tStatisticsGameMoney
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tStatisticsGameMoney` (
    `nRowNo` BIGINT NOT NULL AUTO_INCREMENT,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    `sCode` VARCHAR(255) NOT NULL DEFAULT '',
    `nMoney` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nRowNo`),
    KEY `idx_date` (`dDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tTutorial
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tTutorial` (
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `nState` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `nStep` TINYINT UNSIGNED NOT NULL DEFAULT 0,
    `dDate` DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`nUserNo`, `nCharNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tUserMoney
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tUserMoney` (
    `nUserNo` INT NOT NULL DEFAULT 0,
    `nUserMoney` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nUserNo`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------------------------------------------------------
-- Table: tGameLog (World00_GameLog)
-- Tabla oficial única de auditoría de gameplay (login/logout, level, comercio, crafteo, etc.)
-- Invocada originalmente por los 23 Stored Procedures de GameLog.exe (p_GameLog_Add*)
-- ----------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS `tGameLog` (
    `nLogNo` BIGINT NOT NULL AUTO_INCREMENT,
    `dDate` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `nType` INT NOT NULL DEFAULT 0,
    `nCharNo` INT NOT NULL DEFAULT 0,
    `sMap` VARCHAR(16) NOT NULL DEFAULT '',
    `nMapX` INT NOT NULL DEFAULT 0,
    `nMapY` INT NOT NULL DEFAULT 0,
    `nMapZ` INT NOT NULL DEFAULT 0,
    `nTargetCharNo` INT NOT NULL DEFAULT 0,
    `nTargetID` INT NOT NULL DEFAULT 0,
    `nItemKey` BIGINT NOT NULL DEFAULT 0,
    `nInt1` INT NOT NULL DEFAULT 0,
    `nInt2` INT NOT NULL DEFAULT 0,
    `nInt3` INT NOT NULL DEFAULT 0,
    `nBigint1` BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (`nLogNo`),
    KEY `idx_gamelog_char` (`nCharNo`, `dDate`),
    KEY `idx_gamelog_type` (`nType`, `dDate`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ============================================================================
-- SEED DATA OFICIAL Y CUENTAS DE TEST
-- ============================================================================

-- Roles oficiales
INSERT IGNORE INTO `tUserAuth` (`nAuthID`, `sAuthName`, `bIsLoginAble`) VALUES
(1, 'User', 1),
(2, 'SubGM', 1),
(3, 'GM', 1),
(4, 'Admin', 1);

-- Usuarios de prueba iniciales (passwords en hash MD5 estandar de Fiesta)
-- admin / admin -> MD5: 21232f297a57a5a743894a0e4a801fc3
-- test / test   -> MD5: 098f6bcd4621d373cade4e832627b4f6
INSERT IGNORE INTO `tUser` (`nUserNo`, `sUserID`, `sUserPW`, `sUserName`, `bIsBlock`, `bIsLoginAble`, `bIsDelete`, `nAuthID`, `sUserIP`, `dDate`) VALUES
(1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 'Administrator', 0, 1, 0, 4, '127.0.0.1', NOW()),
(2, 'test', '098f6bcd4621d373cade4e832627b4f6', 'TestUser', 0, 1, 0, 1, '127.0.0.1', NOW());
