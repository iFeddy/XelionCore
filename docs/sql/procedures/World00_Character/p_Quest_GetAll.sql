CREATE PROCEDURE [dbo].[p_Quest_GetAll]
/*
Quest Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nQuestNo, nStatus, sData FROM tQuest WHERE nCharNo = @nCharNo

-- end

0蜁㎁Ċ
GO
