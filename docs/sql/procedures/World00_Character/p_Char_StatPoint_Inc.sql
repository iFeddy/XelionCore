CREATE PROCEDURE [dbo].[p_Char_StatPoint_Inc]
/*
Character increment statpoint
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

IF @nStatPointType = 0 UPDATE tCharacter SET nStrength		= nStrength + 1		WHERE nCharNo = @nCharNo
IF @nStatPointType = 1 UPDATE tCharacter SET nConstitute	= nConstitute + 1	WHERE nCharNo = @nCharNo
IF @nStatPointType = 2 UPDATE tCharacter SET nDexterity		= nDexterity + 1	WHERE nCharNo = @nCharNo
IF @nStatPointType = 3 UPDATE tCharacter SET nIntelligence	= nIntelligence + 1 WHERE nCharNo = @nCharNo
IF @nStatPointType = 4 UPDATE tCharacter SET nMentalPower	= nMentalPower + 1	WHERE nCharNo = @nCharNo

-- check error
IF @@ERROR <> 0 OR @nStatPointType > 4
BEGIN
	SET @nRet = 0
	RETURN
END

SET @nRet = 1

-- end of file

0鐁ꏖą
GO
