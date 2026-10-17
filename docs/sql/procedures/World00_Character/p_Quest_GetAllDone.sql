CREATE PROCEDURE [dbo].[p_Quest_GetAllDone]
/*
Quest Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nStatusDone tinyint

AS

SET NOCOUNT ON

SELECT nQuestNo FROM tQuest WHERE nCharNo = @nCharNo AND nStatus = @nStatusDone

-- end

0锁㤜Ď
GO
