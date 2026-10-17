CREATE PROCEDURE [dbo].[usp_ChargeItem_UseSkillPowerInitLotChange]
/*
Charge item use - clear character all skill power and use item lot change
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character all skill power
	        2 = Error used charge item change lot
*/

@nCharNo int,
@nChangeItemKey bigint,
@nChangeItemOptionType smallint,
@nChangeItemOptionData bigint,		
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- Clear character all skill power
UPDATE tSkill SET 
	nSkillPowerDemage	= 0,
	nSkillPowerSP		= 0,
	nSkillPowerKeepTime = 0,
	nSkillPowerCoolTime = 0
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 1
	RETURN
END

-- Used charge item change lot
DECLARE @TempRet int
EXEC p_Item_SetOption @nChangeItemKey, @nChangeItemOptionType, @nChangeItemOptionData, @TempRet OUTPUT
IF @TempRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 2
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0猁퉻F
GO
