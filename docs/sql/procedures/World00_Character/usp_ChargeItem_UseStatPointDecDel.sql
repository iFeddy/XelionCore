CREATE PROCEDURE [dbo].[usp_ChargeItem_UseStatPointDecDel]
/*
Charge item use - Stat point decrement and use item delete
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Error decrement character stat point.
	        2 = Error delete used charge item
*/

@nCharNo int,
@nStatPointType tinyint,
@nDeleteItemKey bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

BEGIN TRAN

-- CHECK Type and Decrement stat-point
DECLARE @nStatPoint tinyint

IF @nStatPointType = 0
BEGIN
	SELECT @nStatPoint = nStrength FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nStrength	= nStrength - 1					WHERE nCharNo = @nCharNo AND nStrength > 0
END

IF @nStatPointType = 1
BEGIN
	SELECT @nStatPoint = nConstitute FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nConstitute = nConstitute - 1				WHERE nCharNo = @nCharNo AND nConstitute > 0
END
	
IF @nStatPointType = 2
BEGIN
	SELECT @nStatPoint = nDexterity FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nDexterity = nDexterity - 1					WHERE nCharNo = @nCharNo AND nDexterity > 0
END

IF @nStatPointType = 3
BEGIN
	SELECT @nStatPoint = nIntelligence FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nIntelligence = nIntelligence - 1			WHERE nCharNo = @nCharNo AND nIntelligence > 0
END

IF @nStatPointType = 4 
BEGIN
	SELECT @nStatPoint = nMentalPower FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
	
	UPDATE tCharacter	SET nMentalPower = nMentalPower - 1				WHERE nCharNo = @nCharNo AND nMentalPower > 0
END

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

0舁醕G
GO
