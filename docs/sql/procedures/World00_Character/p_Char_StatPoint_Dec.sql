CREATE PROCEDURE [dbo].[p_Char_StatPoint_Dec]
/*
Character decrement statpoint
2005.11 By CJC

 input:
   See SQL
   
 output: 
   0 = error
   1 = update ok
   
*/

@nCharNo int,
@nStatPointType tinyint,

-- Output var
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- CHECK nRedistributePoint

/* DO NOT USE nRedistributePoint

DECLARE @nRedistributePoint tinyint
SET @nRedistributePoint = 0
SELECT @nRedistributePoint = nRedistributePoint FROM tCharacter WHERE nCharNo = @nCharNo
IF @nRedistributePoint < 1
BEGIN
	SET @nRet = 0
	RETURN
END
*/
-- CHECK Type and Decrement stat-point
DECLARE @nStatPoint tinyint

IF @nStatPointType = 0
BEGIN
	SELECT @nStatPoint = nStrength FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nStrength	= nStrength - 1					WHERE nCharNo = @nCharNo AND nStrength > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 1
BEGIN
	SELECT @nStatPoint = nConstitute FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nConstitute = nConstitute - 1				WHERE nCharNo = @nCharNo AND nConstitute > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1	WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END
	
IF @nStatPointType = 2
BEGIN
	SELECT @nStatPoint = nDexterity FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nDexterity = nDexterity - 1					WHERE nCharNo = @nCharNo AND nDexterity > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1	WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 3
BEGIN
	SELECT @nStatPoint = nIntelligence FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nIntelligence = nIntelligence - 1			WHERE nCharNo = @nCharNo AND nIntelligence > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

IF @nStatPointType = 4 
BEGIN
	SELECT @nStatPoint = nMentalPower FROM tCharacter WHERE nCharNo = @nCharNo
	IF @nStatPoint < 1
	BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tCharacter	SET nMentalPower = nMentalPower - 1				WHERE nCharNo = @nCharNo AND nMentalPower > 0
--	UPDATE tCharacter	SET nRedistributePoint = nRedistributePoint - 1 WHERE nCharNo = @nCharNo AND nRedistributePoint > 0
END

-- Check Error
IF @@ERROR <> 0 OR @nStatPointType > 4
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

-- end of file

0蔁Ą
GO
