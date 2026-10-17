CREATE PROCEDURE [dbo].[p_Char_FameAdd]
/*
Add fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nAddFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharFame int
EXEC p_Char_FameGet @nCharNo, @nCharFame OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add fame
SET @nCharFame = @nCharFame + @nAddFame
EXEC p_Char_FameSet @nCharNo, @nCharFame, @nRet OUTPUT
-- End

0瀁쨷ń
GO
