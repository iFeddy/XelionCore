CREATE PROCEDURE [dbo].[p_KingdomQuest_GetNumOfKQ]
/*
Kingdom Quest Get number of KQ
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/

@nNumOfKQ int OUTPUT

AS

SET NOCOUNT ON

SELECT @nNumOfKQ = COUNT(*) FROM tKingdomQuest

-- end

0蔁ࣵŹ
GO
