CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointInitLotChange]
/*
Charge item use - Stat point init and use item lot change
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character stat point.
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

-- Clear character stat point.
UPDATE tCharacter SET 
	nStrength = 0, 
	nConstitute = 0, 
	nDexterity = 0, 
	nIntelligence = 0, 
	nMentalPower = 0
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

0ⴁ渂J
GO
