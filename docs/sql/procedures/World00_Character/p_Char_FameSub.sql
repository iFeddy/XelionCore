CREATE PROCEDURE [dbo].[p_Char_FameSub]
/*
Sub fame
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nCharNo int,
@nSubFame int,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharFame int
EXEC p_Char_FameGet @nCharNo, @nCharFame OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub fame
SET @nCharFame = @nCharFame - @nSubFame
-- Check minus
IF @nCharFame < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_Char_FameSet @nCharNo, @nCharFame, @nRet OUTPUT
-- End

0뉿ņ
GO
