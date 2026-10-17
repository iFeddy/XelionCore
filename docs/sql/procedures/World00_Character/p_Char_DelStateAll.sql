CREATE PROCEDURE [dbo].[p_Char_DelStateAll]
/*
Char Delete state all
2004.10 By CJC
 input:
	nCharNo, nStateNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterState WHERE nCharNo = @nCharNo
SET @nRet = @@ERROR
-- end

0㜁혓Ń
GO
