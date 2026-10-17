CREATE PROCEDURE [dbo].[p_KingdomQuest_GetAll]
/*
Kingdom Quest Get All
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

AS

SET NOCOUNT ON

SELECT nQuestNo, sData FROM tKingdomQuest

-- end

0㈣ź
GO
