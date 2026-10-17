CREATE PROCEDURE [dbo].[p_Char_FameSet]
/*
Set fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nSetFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tCharacter SET nFame = @nSetFame WHERE nCharNo = @nCharNo
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end
RETURN

0ఁņ
GO
