CREATE PROCEDURE [dbo].[p_Char_SetOptVideo]
/*
Character Set options
2004.6 By CJC
 input:
   All args
 output: 
   0 = error
   1 = update ok
*/

@nCharNo int,
@sData varbinary(60),
@nRet int OUTPUT
AS

SET NOCOUNT ON

UPDATE tCharacterOptions SET sVideo = @sData WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tCharacterOptions (nCharNo, sVideo) VALUES (@nCharNo, @sData)
	SET @nRet = @@ROWCOUNT
	RETURN
END

SET @nRet = 1

-- end

0봁䎚Ž
GO
