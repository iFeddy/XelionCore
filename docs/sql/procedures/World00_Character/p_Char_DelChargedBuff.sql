CREATE PROCEDURE [dbo].[p_Char_DelChargedBuff]
/*
Char Delete charged buff
2006.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nKey int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo AND nKey = @nKey
SET @nRet = @@ERROR

RETURN

-- end

0ᨁᅞĿ
GO
