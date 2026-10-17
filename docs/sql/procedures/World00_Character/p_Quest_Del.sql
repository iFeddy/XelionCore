CREATE PROCEDURE [dbo].[p_Quest_Del]
/*
Quest Delete
2004.10 By CJC
 input:
	nCharNo, nQuestNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nQuestNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo
IF (@@ERROR=0) AND (@@ROWCOUNT=1)
	SET @nRet = 0
ELSE
	SET @nRet = 1
-- end

0ᔁ䬹Ĉ
GO
