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
