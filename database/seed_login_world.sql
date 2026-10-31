-- ============================================================================
-- XELION ONLINE - SEMILLAS OFICIALES (MYSQL)
-- LoginService, WorldServer, Autenticación y Versión de Cliente
-- Idempotente: seguro para re-ejecución sin pérdida de datos
-- ============================================================================

USE `xelion_online`;

-- 1. FOUS: Clave de versión del cliente (sin pisar clave personalizada existente)
INSERT INTO `FOUS` (`GameAccountTitle`, `GameAccountGameServerReference`, `GameAccountKey`, `GameVersionKey`)
VALUES ('FiestaOnline', 0, 0, '3dca91289ab564b3ce1f5c8ca1b362a6')
ON DUPLICATE KEY UPDATE
    `GameAccountTitle` = `GameAccountTitle`;

-- 2. tUserAuth: Niveles oficiales de privilegios y autorización de login
INSERT INTO `tUserAuth` (`nAuthID`, `sAuthName`, `bIsLoginAble`) VALUES
(0, 'Blocked', 0),
(1, 'Player', 1),
(2, 'Tester', 1),
(3, 'GM', 1),
(4, 'Admin', 1),
(9, 'Admin', 1),
(100, 'Admin', 1)
ON DUPLICATE KEY UPDATE
    `sAuthName` = VALUES(`sAuthName`),
    `bIsLoginAble` = VALUES(`bIsLoginAble`);

-- 3. Auto_Account_Creation: Creación automática de cuentas r_*
INSERT INTO `Auto_Account_Creation` (`Allow_Creation`)
SELECT 1
WHERE NOT EXISTS (SELECT 1 FROM `Auto_Account_Creation`);

-- 4. tWorldServer: Configuración de mundos gobernados 100% por MySQL (sin pisar ediciones manuales)
INSERT INTO `tWorldServer` (
    `nWorldNo`, `sWorldName`, `sWorldIP`, `nWorldPort`, 
    `nMaxUser`, `nUserRatable`, `bIsMaintenance`, `nCurrentCCU`, `dLastHeartbeat`
) VALUES
(0, 'Xelion', '192.168.0.4', 9014, 1500, 100, 0, 0, NOW()),
(1, 'Bijou', '127.0.0.1', 9015, 1500, 100, 0, 0, NOW())
ON DUPLICATE KEY UPDATE
    `nWorldNo` = `nWorldNo`;

-- 5. Cuentas base iniciales (admin y test) sin pisar cuentas existentes
-- admin / admin -> MD5: 21232f297a57a5a743894a0e4a801fc3
-- test / test   -> MD5: 098f6bcd4621d373cade4e832627b4f6
INSERT INTO `tUser` (
    `nUserNo`, `sUserID`, `sUserPW`, `sUserName`, 
    `bIsBlock`, `bIsDelete`, `nAuthID`, `sUserIP`, `dDate`, `sUserNo`, `blsBlock`, `blsDelete`
) VALUES
(1, 'admin', '21232f297a57a5a743894a0e4a801fc3', 'Administrator', 0, 0, 4, '127.0.0.1', NOW(), NULL, 0, 0)
ON DUPLICATE KEY UPDATE `sUserID` = `sUserID`;

INSERT INTO `tUserOption` (`nUserNo`, `sEmail`)
VALUES (1, 'admin@xelion.online')
ON DUPLICATE KEY UPDATE `sEmail` = VALUES(`sEmail`);

INSERT INTO `tUser` (
    `nUserNo`, `sUserID`, `sUserPW`, `sUserName`, 
    `bIsBlock`, `bIsDelete`, `nAuthID`, `sUserIP`, `dDate`, `sUserNo`, `blsBlock`, `blsDelete`
) VALUES
(2, 'test', '098f6bcd4621d373cade4e832627b4f6', 'TestUser', 0, 0, 1, '127.0.0.1', NOW(), NULL, 0, 0)
ON DUPLICATE KEY UPDATE `sUserID` = `sUserID`;

INSERT INTO `tUserOption` (`nUserNo`, `sEmail`)
VALUES (2, 'test@xelion.online')
ON DUPLICATE KEY UPDATE `sEmail` = VALUES(`sEmail`);
