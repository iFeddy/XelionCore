CREATE PROCEDURE [dbo].[p_KingdomQuest_Get]
/*
Kingdom Quest Get
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

@nQuestNo int,
@sData varchar(256) OUTPUT,
@nRet int OUTPUT

AS

SET NOCOUNT ON

SELECT @sData = sData FROM tKingdomQuest WHERE nQuestNo = @nQuestNo

SET @nRet = @@ROWCOUNT

-- end

0ጁ₭ŷ
GO
