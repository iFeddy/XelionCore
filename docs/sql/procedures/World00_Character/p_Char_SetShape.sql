CREATE PROCEDURE [dbo].[p_Char_SetShape]
/*
Character Set shape
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/
@nCharNo int,
@nRace tinyint, 
@nClass tinyint, 
@nGender tinyint, 
@nHairType tinyint, 
@nHairColor tinyint, 
@nFaceShape tinyint, 
-- Output var
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tCharacterShape SET
	nRace		= @nRace,
	nClass		= @nClass,
	nGender		= @nGender,
	nHairType	= @nHairType,
	nHairColor	= @nHairColor,
	nFaceShape	= @nFaceShape
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1
-- end

0ԁž
GO
