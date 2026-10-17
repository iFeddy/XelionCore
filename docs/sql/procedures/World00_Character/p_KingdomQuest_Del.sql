CREATE PROCEDURE [dbo].[p_KingdomQuest_Del]
/*
Kingdom Quest Delete
2005.8 By CJC
 input:
	nQuestNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nQuestNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tKingdomQuest WHERE nQuestNo = @nQuestNo

IF (@@ERROR=0) AND (@@ROWCOUNT=1)
	SET @nRet = 0
ELSE
	SET @nRet = 1
	
-- end

0瘁䧛Ÿ
GO
