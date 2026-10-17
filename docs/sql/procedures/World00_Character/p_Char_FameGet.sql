CREATE PROCEDURE [dbo].[p_Char_FameGet]
/*
Get fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Fame return
	        0 = Error
*/
@nCharNo int,
@nCharFame int OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
-- Get fame
SELECT @nCharFame = nFame FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 1
	SET @nRet = 1
	
-- End
RETURN

0꤁빛Ņ
GO
