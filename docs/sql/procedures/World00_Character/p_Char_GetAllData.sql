CREATE PROCEDURE [dbo].[p_Char_GetAllData]
/*
Character Get all data
2004.6 By CJC
 input:
   nCharNo = To be get login data
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON

/*-------------------------------------------------------------------------------------
Project: FiestaHeroes - NA2016
Author: xkl

Rebirth of characters that can be customized through the settings in:
Options.dbo.Character_Rebirth
------------------------------------------------------------------------------------*/

-- Retrieve values from Character_Rebirth.
DECLARE @Trigger_ItemID INT, @Trigger_Level INT, @Reward_Abstate INT, @Unequip_Character_Items INT, @CharacterLevel INT;
SELECT @Trigger_ItemID = Trigger_ItemID, @Trigger_Level = Trigger_Level, @Reward_Abstate = Reward_Abstate, @Unequip_Character_Items = Unequip_Character_Items FROM Options..Character_Rebirth WHERE Enabled = 1;

-- Get character level.
SELECT @CharacterLevel = nLevel FROM tCharacter WHERE nCharNo = @nCharNo;

-- Update the rebirth Abstate time for characters that already have it.
UPDATE tCharacterState SET nStateTime = 100000000, nRestTime = 100000000 WHERE nCharNo = @nCharNo AND nStateNo = @Reward_Abstate;

-- Begin the character rebirth process.
IF EXISTS (SELECT 1 FROM Options..Character_Rebirth WHERE Enabled = 1)
BEGIN
	-- Check if the character has the required item.
	IF EXISTS (SELECT 1 FROM tItem WHERE nOwner = @nCharNo AND nItemID = @Trigger_ItemID)
	BEGIN
	    -- Check if the character meets the required level.
	    IF @CharacterLevel = @Trigger_Level
	    BEGIN
	        -- Apply rebirth.
	        DELETE FROM tItem WHERE nOwner = @nCharNo AND nItemID = @Trigger_ItemID;
	        UPDATE tCharacter SET nLevel = 1, nExp = 0 WHERE nCharNo = @nCharNo;
	
	        -- Check for unequip setting.
	        IF @Unequip_Character_Items = 1
	            UPDATE tItem SET nStorageType = 2 WHERE nOwner = @nCharNo AND nStorageType = 8;
	
	        -- Reward the Abstate.
	        IF NOT EXISTS (SELECT 1 FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @Reward_Abstate)
	            INSERT INTO tCharacterState (nCharNo, nStateNo, nStateLevel, nStateTime, nRestTime) 
	            VALUES (@nCharNo, @Reward_Abstate, 1, 100000000, 100000000);
	    END
	END
END
/*-------------------------------------------------------------------------------------
 Continue as normal ~ 
------------------------------------------------------------------------------------*/

SELECT
nUserNo,
nSlotNo,
sID,
nLevel,
nExp,
nHP,
nSP,
nAP,
nHPS,
nSPS,
nPS,
nGS,
nMoney,
nFame,
nFlags,
nAdminLevel,
sLoginZone,
nLoginZoneX,
nLoginZoneY,
nLoginZoneD,
nPrisonMin,
nPKCount,
nPKYellowTime,
nStrength,
nConstitute,
nDexterity,
nIntelligence,
nMentalPower,
nRedistributePoint,
nPlayMin,
nLoginCount,
dLastLoginDate,
dCreateDate,
nCreateWorld,
bDeleted,
dDeletedDate
FROM tCharacter
WHERE nCharNo = @nCharNo
-- end
GO
