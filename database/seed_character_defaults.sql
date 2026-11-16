-- ============================================================================
-- XELION ONLINE - SEMILLAS OFICIALES DE WORLD00_CHARACTER (MYSQL 8.4)
-- Inicialización de Parámetros Oficiales de Personajes (DefaultCharacterData.txt / TutorialCharacterData.txt)
-- Idempotente: seguro para re-ejecución sin pérdida de datos
-- ============================================================================

USE `xelion_online`;

-- 1. tGameOptionInit: Opciones oficiales de cliente y configuración de juego (extraídas de World00_Character.bak)
TRUNCATE TABLE `tGameOptionInit`;
INSERT INTO `tGameOptionInit` (`nOptionNo`, `sys_create_date`, `sys_update_date`, `nDefaultValue`, `bIsActive`, `sDescription`, `sDescriptionLocal`) VALUES
(0, NOW(), NOW(), 1, 1, "Whisper Allow", "Whisper"),
(1, NOW(), NOW(), 1, 1, "Friend Invite Allow", "FriendInvite"),
(2, NOW(), NOW(), 1, 1, "NPC Name Display", "NPCName"),
(3, NOW(), NOW(), 1, 1, "Monster Name Display", "MonsterName"),
(4, NOW(), NOW(), 0, 1, "Player Name Display", "PlayerName"),
(5, NOW(), NOW(), 1, 1, "HP Bar Display", "HPBar"),
(6, NOW(), NOW(), 0, 1, "SP Bar Display", "SPBar"),
(7, NOW(), NOW(), 1, 1, "Chat Balloon Display", "ChatBalloon"),
(8, NOW(), NOW(), 1, 1, "Booth Sign Display", "BoothSign"),
(9, NOW(), NOW(), 1, 1, "3D Sound", "3DSound"),
(10, NOW(), NOW(), 1, 1, "UI Tip Display", "UITip"),
(11, NOW(), NOW(), 1, 1, "Auto Target", "AutoTarget"),
(12, NOW(), NOW(), 0, 1, "Trade Reject", "TradeReject"),
(13, NOW(), NOW(), 0, 1, "Party Reject", "PartyReject"),
(14, NOW(), NOW(), 0, 1, "Guild Invite Reject", "GuildReject"),
(15, NOW(), NOW(), 0, 1, "Master Pupil Reject", "MasterReject"),
(16, NOW(), NOW(), 1, 1, "My Guild Mark Display", "MyGuildMark"),
(17, NOW(), NOW(), 0, 1, "Other Guild Mark Hide", "OtherGuildMark"),
(18, NOW(), NOW(), 0, 1, "Camera Shake Off", "CameraShake"),
(19, NOW(), NOW(), 0, 1, "Weather Effect Off", "WeatherEffect"),
(20, NOW(), NOW(), 1, 1, "In-Game Help Display", "InGameHelp"),
(21, NOW(), NOW(), 0, 1, "Other Effect Off", "OtherEffect"),
(22, NOW(), NOW(), 0, 1, "MiniHouse Invite Reject", "MiniHouseReject"),
(23, NOW(), NOW(), 0, 1, "Pet Hide", "PetHide"),
(24, NOW(), NOW(), 1, 1, "mover Hide", "MoverHide"),
(25, NOW(), NOW(), 1, 1, "Costume Display", "CostumeDisplay"),
(26, NOW(), NOW(), 1, 1, "Title Display", "TitleDisplay"),
(27, NOW(), NOW(), 1, 1, "Damage Number Display", "DamageDisplay")
ON DUPLICATE KEY UPDATE
    `nDefaultValue` = VALUES(`nDefaultValue`),
    `bIsActive` = VALUES(`bIsActive`);

-- 2. tKeyMapInit: Mapeo oficial de funciones de teclado 0..94 (extraído byte a byte de World00_Character.bak)
DELETE FROM `tKeyMapInit`;
INSERT INTO `tKeyMapInit` (`nFunctionNo`, `sys_create_date`, `sys_update_date`, `nExtendedKey`, `nASCIICode`, `sDescription`) VALUES
(0, NOW(), NOW(), 0, 121, "System Menu (F10)"),
(1, NOW(), NOW(), 0, 27, "Cancel / Escape (ESC)"),
(2, NOW(), NOW(), 0, 67, "Character Info (C)"),
(3, NOW(), NOW(), 0, 73, "Inventory (I)"),
(4, NOW(), NOW(), 0, 75, "Skill Window (K)"),
(5, NOW(), NOW(), 0, 76, "Quest Window (L)"),
(6, NOW(), NOW(), 0, 70, "Community / Friends (F)"),
(7, NOW(), NOW(), 0, 72, "MiniHouse (H)"),
(8, NOW(), NOW(), 0, 86, "Action Window (V)"),
(9, NOW(), NOW(), 0, 13, "Chat Input (Enter)"),
(10, NOW(), NOW(), 17, 78, "Chat Mode Normal (Ctrl+N)"),
(11, NOW(), NOW(), 17, 71, "Chat Mode Guild (Ctrl+G)"),
(12, NOW(), NOW(), 17, 80, "Chat Mode Party (Ctrl+P)"),
(13, NOW(), NOW(), 17, 87, "Chat Mode Whisper (Ctrl+W)"),
(14, NOW(), NOW(), 0, 222, "Reply Whisper (')"),
(15, NOW(), NOW(), 0, 88, "Sit / Rest (X)"),
(16, NOW(), NOW(), 0, 71, "Guild Window (G)"),
(17, NOW(), NOW(), 0, 0, "Reserved 17"),
(18, NOW(), NOW(), 0, 0, "Reserved 18"),
(19, NOW(), NOW(), 0, 0, "Reserved 19"),
(20, NOW(), NOW(), 0, 0, "Reserved 20"),
(21, NOW(), NOW(), 0, 82, "Auto Run (R)"),
(22, NOW(), NOW(), 17, 65, "Select Self (Ctrl+A)"),
(23, NOW(), NOW(), 0, 87, "Move Forward (W)"),
(24, NOW(), NOW(), 0, 83, "Move Backward (S)"),
(25, NOW(), NOW(), 0, 0, "Reserved 25"),
(26, NOW(), NOW(), 0, 65, "Turn Left (A)"),
(27, NOW(), NOW(), 0, 68, "Turn Right (D)"),
(28, NOW(), NOW(), 0, 90, "Walk / Run Toggle (Z)"),
(29, NOW(), NOW(), 0, 32, "Jump (Space)"),
(30, NOW(), NOW(), 0, 38, "Camera Up (Up)"),
(31, NOW(), NOW(), 0, 40, "Camera Down (Down)"),
(32, NOW(), NOW(), 0, 37, "Camera Left (Left)"),
(33, NOW(), NOW(), 0, 39, "Camera Right (Right)"),
(34, NOW(), NOW(), 0, 36, "Camera Reset (Home)"),
(35, NOW(), NOW(), 0, 84, "Target Nearest (T)"),
(36, NOW(), NOW(), 0, 81, "Use HP Stone (Q)"),
(37, NOW(), NOW(), 0, 69, "Use SP Stone (E)"),
(38, NOW(), NOW(), 0, 245, "Auto Attack / Pick"),
(39, NOW(), NOW(), 0, 66, "Booth / Store (B)"),
(40, NOW(), NOW(), 0, 80, "Party Window (P)"),
(41, NOW(), NOW(), 0, 77, "World Map (M)"),
(42, NOW(), NOW(), 0, 85, "Upgrade / Title (U)"),
(43, NOW(), NOW(), 16, 90, "Hide UI (Shift+Z)"),
(44, NOW(), NOW(), 0, 0, "PK Mode"),
(45, NOW(), NOW(), 0, 0, "Reserved 45"),
(46, NOW(), NOW(), 0, 35, "Camera View (End)"),
(47, NOW(), NOW(), 0, 49, "Slot 1-1"),
(48, NOW(), NOW(), 0, 50, "Slot 1-2"),
(49, NOW(), NOW(), 0, 51, "Slot 1-3"),
(50, NOW(), NOW(), 0, 52, "Slot 1-4"),
(51, NOW(), NOW(), 0, 53, "Slot 1-5"),
(52, NOW(), NOW(), 0, 54, "Slot 1-6"),
(53, NOW(), NOW(), 0, 55, "Slot 1-7"),
(54, NOW(), NOW(), 0, 56, "Slot 1-8"),
(55, NOW(), NOW(), 0, 57, "Slot 1-9"),
(56, NOW(), NOW(), 0, 48, "Slot 1-10"),
(57, NOW(), NOW(), 0, 189, "Slot 1-11"),
(58, NOW(), NOW(), 0, 187, "Slot 1-12"),
(59, NOW(), NOW(), 16, 49, "Slot 2-1"),
(60, NOW(), NOW(), 16, 50, "Slot 2-2"),
(61, NOW(), NOW(), 16, 51, "Slot 2-3"),
(62, NOW(), NOW(), 16, 52, "Slot 2-4"),
(63, NOW(), NOW(), 16, 53, "Slot 2-5"),
(64, NOW(), NOW(), 16, 54, "Slot 2-6"),
(65, NOW(), NOW(), 16, 55, "Slot 2-7"),
(66, NOW(), NOW(), 16, 56, "Slot 2-8"),
(67, NOW(), NOW(), 16, 57, "Slot 2-9"),
(68, NOW(), NOW(), 16, 48, "Slot 2-10"),
(69, NOW(), NOW(), 16, 189, "Slot 2-11"),
(70, NOW(), NOW(), 16, 187, "Slot 2-12"),
(71, NOW(), NOW(), 18, 49, "Slot 3-1"),
(72, NOW(), NOW(), 18, 50, "Slot 3-2"),
(73, NOW(), NOW(), 18, 51, "Slot 3-3"),
(74, NOW(), NOW(), 18, 52, "Slot 3-4"),
(75, NOW(), NOW(), 18, 53, "Slot 3-5"),
(76, NOW(), NOW(), 18, 54, "Slot 3-6"),
(77, NOW(), NOW(), 18, 55, "Slot 3-7"),
(78, NOW(), NOW(), 18, 56, "Slot 3-8"),
(79, NOW(), NOW(), 18, 57, "Slot 3-9"),
(80, NOW(), NOW(), 18, 48, "Slot 3-10"),
(81, NOW(), NOW(), 18, 189, "Slot 3-11"),
(82, NOW(), NOW(), 18, 187, "Slot 3-12"),
(83, NOW(), NOW(), 0, 0, "Slot 4-1"),
(84, NOW(), NOW(), 0, 0, "Slot 4-2"),
(85, NOW(), NOW(), 0, 0, "Slot 4-3"),
(86, NOW(), NOW(), 0, 0, "Slot 4-4"),
(87, NOW(), NOW(), 0, 0, "Slot 4-5"),
(88, NOW(), NOW(), 0, 0, "Slot 4-6"),
(89, NOW(), NOW(), 0, 0, "Slot 4-7"),
(90, NOW(), NOW(), 0, 0, "Slot 4-8"),
(91, NOW(), NOW(), 0, 0, "Slot 4-9"),
(92, NOW(), NOW(), 0, 0, "Slot 4-10"),
(93, NOW(), NOW(), 0, 0, "Slot 4-11"),
(94, NOW(), NOW(), 0, 0, "Slot 4-12")
ON DUPLICATE KEY UPDATE
    `nExtendedKey` = VALUES(`nExtendedKey`),
    `nASCIICode` = VALUES(`nASCIICode`);

-- 3. tExp: Tabla oficial de niveles y progresión de experiencia (Niveles 1 a 150)
INSERT INTO `tExp` (`level`, `NextExp`, `sExp`, `eExp`) VALUES
(1, 10, 0, 10),
(2, 40, 10, 50),
(3, 100, 50, 150),
(4, 200, 150, 350),
(5, 350, 350, 700),
(6, 560, 700, 1260),
(7, 840, 1260, 2100),
(8, 1200, 2100, 3300),
(9, 1650, 3300, 4950),
(10, 2350, 4950, 7300),
(11, 3340, 7300, 10640),
(12, 4671, 10640, 15311),
(13, 6401, 15311, 21712),
(14, 8597, 21712, 30309),
(15, 11416, 30309, 41725),
(16, 14959, 41725, 56684),
(17, 19987, 56684, 76671),
(18, 26203, 76671, 102874),
(19, 33843, 102874, 136717),
(20, 43539, 136717, 180256),
(21, 54952, 180256, 235208),
(22, 68451, 235208, 303659),
(23, 84419, 303659, 388078),
(24, 103307, 388078, 491385),
(25, 125647, 491385, 617032),
(26, 152070, 617032, 769102),
(27, 183322, 769102, 952424),
(28, 220284, 952424, 1172708),
(29, 263998, 1172708, 1436706),
(30, 314987, 1436706, 1751693),
(31, 374745, 1751693, 2126438),
(32, 444441, 2126438, 2570879),
(33, 525720, 2570879, 3096599),
(34, 620499, 3096599, 3717098),
(35, 731012, 3717098, 4448110),
(36, 859859, 4448110, 5307969),
(37, 1010070, 5307969, 6318039),
(38, 1185174, 6318039, 7503213),
(39, 1389279, 7503213, 8892492),
(40, 1628577, 8892492, 10521069),
(41, 1903735, 10521069, 12424804),
(42, 2220074, 12424804, 14644878),
(43, 2583699, 14644878, 17228577),
(44, 3001612, 17228577, 20230189),
(45, 3481845, 20230189, 23712034),
(46, 4033608, 23712034, 27745642),
(47, 4667463, 27745642, 32413105),
(48, 5395518, 32413105, 37808623),
(49, 6231658, 37808623, 44040281),
(50, 7186574, 44040281, 51226855),
(51, 8280054, 51226855, 59506909),
(52, 9519838, 59506909, 69026747),
(53, 10925018, 69026747, 79951765),
(54, 12517133, 79951765, 92468898),
(55, 14320476, 92468898, 106789374),
(56, 16362442, 106789374, 123151816),
(57, 18673920, 123151816, 141825736),
(58, 21289729, 141825736, 163115465),
(59, 24346145, 163115465, 187461610),
(60, 27892956, 187461610, 215354566),
(61, 31941383, 215354566, 247295949),
(62, 36553632, 247295949, 283849581),
(63, 41799105, 283849581, 325648686),
(64, 47755192, 325648686, 373403878),
(65, 54508159, 373403878, 427912037),
(66, 62154117, 427912037, 490066154),
(67, 70800112, 490066154, 560866266),
(68, 80565316, 560866266, 641431582),
(69, 91582354, 641431582, 733013936),
(70, 103941051, 733013936, 836954987),
(71, 117311787, 836954987, 954266774),
(72, 131750092, 954266774, 1086016866),
(73, 147308993, 1086016866, 1233325859),
(74, 164037442, 1233325859, 1397363301),
(75, 181978388, 1397363301, 1579341689),
(76, 201166415, 1579341689, 1780508104),
(77, 221624868, 1780508104, 2002132972),
(78, 243362385, 2002132972, 2245495357),
(79, 266368711, 2245495357, 2511864068),
(80, 290846185, 2511864068, 2802710253),
(81, 319376085, 2802710253, 3122086338),
(82, 354507454, 3122086338, 3476593792),
(83, 397048348, 3476593792, 3873642140),
(84, 448664634, 3873642140, 4322306774),
(85, 511477682, 4322306774, 4833784456),
(86, 588199335, 4833784456, 5421983791),
(87, 682311228, 5421983791, 6104295019),
(88, 798304137, 6104295019, 6902599156),
(89, 941998882, 6902599156, 7844598038),
(90, 1109768883, 7844598038, 8954366921),
(91, 1318270915, 8954366921, 10272637836),
(92, 1578831197, 10272637836, 11851469033),
(93, 1906316596, 11851469033, 13757785629),
(94, 2222646006, 13757785629, 15980431635),
(95, 2574361418, 15980431635, 18554793053),
(96, 2952943980, 18554793053, 21507737033),
(97, 3376799595, 21507737033, 24884536628),
(98, 3816589234, 24884536628, 28701125862),
(99, 4256887017, 28701125862, 32958012879),
(100, 4701456361, 32958012879, 37659469240),
(101, 5230623378, 37659469240, 42890092618),
(102, 5791994175, 42890092618, 48682086793),
(103, 6521104030, 48682086793, 55203190823),
(104, 7476783069, 55203190823, 62679973892),
(105, 8616244808, 62679973892, 71296218700),
(106, 10008870423, 71296218700, 81305089123),
(107, 11551859528, 81305089123, 92856948651),
(108, 13305976287, 92856948651, 106162924938),
(109, 15301872730, 106162924938, 121464797668),
(110, 17595762560, 121464797668, 139060560228),
(111, 20302802954, 139060560228, 159363363182),
(112, 23664947123, 159363363182, 183028310305),
(113, 27642185010, 183028310305, 210670495315),
(114, 32552737877, 210670495315, 243223233192),
(115, 39063285452, 243223233192, 282286518644),
(116, 46626986736, 282286518644, 328913505380),
(117, 54726830429, 328913505380, 383640335809),
(118, 63550633157, 383640335809, 447190968966),
(119, 72987672292, 447190968966, 520178641258),
(120, 83037947834, 520178641258, 603216589092),
(121, 93617185246, 603216589092, 696833774338),
(122, 104753224627, 696833774338, 801586998965),
(123, 116475371344, 801586998965, 918062370309),
(124, 128814473151, 918062370309, 1046876843460),
(125, 141803001369, 1046876843460, 1188679844829),
(126, 155816939710, 1188679844829, 1344496784539),
(127, 170983622642, 1344496784539, 1515480407181),
(128, 187303050164, 1515480407181, 1702783457345),
(129, 205254420439, 1702783457345, 1908037877784),
(130, 225000927742, 1908037877784, 2133038805526),
(131, 246722085775, 2133038805526, 2379760891301),
(132, 270615359611, 2379760891301, 2650376250912),
(133, 296897960830, 2650376250912, 2947274211742),
(134, 325808822172, 2947274211742, 3273083033914),
(135, 357610769648, 3273083033914, 3630693803562),
(136, 392592911871, 3630693803562, 4023286715433),
(137, 431073268317, 4023286715433, 4454359983750),
(138, 473401660407, 4454359983750, 4927761644157),
(139, 519962891706, 4927761644157, 5447724535863),
(140, 571180246136, 5447724535863, 6018904781999),
(141, 627519336008, 6018904781999, 6646424118007),
(142, 689492334867, 6646424118007, 7335916452874),
(143, 757662633612, 7335916452874, 8093579086486),
(144, 832649962232, 8093579086486, 8926229048718),
(145, 915136023714, 8926229048718, 9841365072432),
(146, 1005870691344, 9841365072432, 10847235763776),
(147, 1105678825737, 10847235763776, 11952914589513),
(148, 1215467773569, 11952914589513, 13168382363082),
(149, 1336235616185, 13168382363082, 14504617979267),
(150, 1469080243062, 14504617979267, 15973698222329)
ON DUPLICATE KEY UPDATE
    `NextExp` = VALUES(`NextExp`),
    `sExp` = VALUES(`sExp`),
    `eExp` = VALUES(`eExp`);

-- 4. Saneamiento idempotente: eliminar registros falsos de plantillas bajo nUserNo = 1
-- (En Character.exe oficial, DefaultCharacterData.txt y TutorialCharacterData.txt se cargan
-- desde data/ en memoria y se aplican al personaje real vía CALL p_Char_CreateSetDefaultData)
DELETE FROM `tItem`
WHERE `nOwner` IN (
    SELECT `nCharNo` FROM `tCharacter`
    WHERE `nUserNo` = 1 AND `sID` IN (
        'Fighter_Template', 'Cleric_Template', 'Archer_Template',
        'Mage_Template', 'Trickster_Template', 'Sentinel_Template',
        'Tutorial_Fighter', 'Tutorial_Cleric', 'Tutorial_Archer',
        'Tutorial_Mage', 'Tutorial_Trickster'
    )
);

DELETE FROM `tCharacterShape`
WHERE `nCharNo` IN (
    SELECT `nCharNo` FROM `tCharacter`
    WHERE `nUserNo` = 1 AND `sID` IN (
        'Fighter_Template', 'Cleric_Template', 'Archer_Template',
        'Mage_Template', 'Trickster_Template', 'Sentinel_Template',
        'Tutorial_Fighter', 'Tutorial_Cleric', 'Tutorial_Archer',
        'Tutorial_Mage', 'Tutorial_Trickster'
    )
);

DELETE FROM `tCharacter`
WHERE `nUserNo` = 1 AND `sID` IN (
    'Fighter_Template', 'Cleric_Template', 'Archer_Template',
    'Mage_Template', 'Trickster_Template', 'Sentinel_Template',
    'Tutorial_Fighter', 'Tutorial_Cleric', 'Tutorial_Archer',
    'Tutorial_Mage', 'Tutorial_Trickster'
);

-- 5. num: Tabla auxiliar oficial de secuencia entera (0..1000) requerida por SPs de World00_Character.bak
SET @@cte_max_recursion_depth = 2000;

INSERT IGNORE INTO `num` (`num`)
WITH RECURSIVE `seq` (`n`) AS (
    SELECT 0
    UNION ALL
    SELECT `n` + 1 FROM `seq` WHERE `n` < 1000
)
SELECT `n` FROM `seq`;

-- 6. tItemStorageFindEmpty: Índices de slots (0..255) requeridos por p_Item_GetStorageEmptyList / p_Item_Create_UseEmptyStorage
INSERT IGNORE INTO `tItemStorageFindEmpty` (`nStorage`)
WITH RECURSIVE `slot_seq` (`n`) AS (
    SELECT 0
    UNION ALL
    SELECT `n` + 1 FROM `slot_seq` WHERE `n` < 255
)
SELECT `n` FROM `slot_seq`;
