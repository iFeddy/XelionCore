CREATE PROCEDURE [dbo].[p_Quest_Get]
/*
Quest Get
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/
@nCharNo int,
@nQuestNo int,
@nStatus tinyint OUTPUT,
@sData varbinary(100) OUTPUT,
@nRet int OUTPUT

AS

SET NOCOUNT ON

SELECT 
	@nStatus = nStatus, 
	@sData = sData
FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo

SET @nRet = @@ROWCOUNT

-- end

0鶕Ċ
GO
