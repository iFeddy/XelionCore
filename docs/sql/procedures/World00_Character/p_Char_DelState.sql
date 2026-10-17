CREATE PROCEDURE [dbo].[p_Char_DelState]
/*
Char Delete state
2004.10 By CJC
 input:
	nCharNo, nStateNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nStateNo smallint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 해당 상태 삭제
DELETE FROM tCharacterState WHERE nCharNo = @nCharNo AND nStateNo = @nStateNo
SET @nRet = @@ERROR
-- end

0씁Ł
GO
