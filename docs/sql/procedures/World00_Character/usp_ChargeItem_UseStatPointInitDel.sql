CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointInitDel]
/*
Charge item use - Stat point init and use item delete
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error clear character stat point.
	        2 = Error delete used charge item
*/

@nCharNo int,
@nDeleteItemKey bigint,
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

-- Delete used charge item
DECLARE @TempRet int
EXEC p_Item_Delete @nDeleteItemKey, @TempRet OUTPUT
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

0뫃H
GO
