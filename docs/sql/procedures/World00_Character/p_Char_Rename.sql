CREATE PROCEDURE [dbo].[p_Char_Rename]
/*
Char Rename
2006.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@sID nvarchar(40),
-- Output var
@nRet int OUTPUT

AS
BEGIN

SET NOCOUNT ON

IF NOT EXISTS(SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = @nCharNo) BEGIN
	SET @nRet = @@ERROR
	RETURN
END

UPDATE tCharacter SET sID = @sID WHERE nCharNo = @nCharNo
SET @nRet = @@ERROR

END
-- end
GO
