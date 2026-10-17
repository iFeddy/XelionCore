CREATE PROCEDURE [dbo].[p_Char_SetTitle]
/*
Character Set shape
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update/insert ok
*/
@nCharNo		int,
@nType			tinyint, 
@nElementNo		tinyint, 
@nElementValue	bigint,
@nRet			tinyint OUTPUT

AS

SET NOCOUNT ON

-- SET ERROR FIRST
SET @nRet = 0

IF EXISTS( SELECT nCharNo FROM tCharacterTitle WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo )
BEGIN
	UPDATE tCharacterTitle SET nElementValue = @nElementValue
	WHERE nCharNo = @nCharNo AND nType = @nType AND nElementNo = @nElementNo

	IF @@ERROR = 0 AND @@ROWCOUNT = 1
	BEGIN
		SET @nRet = 1
		RETURN
	END
END
ELSE
BEGIN
	INSERT tCharacterTitle (  nCharNo,  nType,  nElementNo,  nElementValue )
	VALUES                 ( @nCharNo, @nType, @nElementNo, @nElementValue )

	IF @@ERROR = 0
	BEGIN
		SET @nRet = 1
		RETURN
	END
END

-- end

0欁⡋
GO
