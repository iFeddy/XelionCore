CREATE PROCEDURE [dbo].[p_Item_GetMobListAll]
/*
Item Get Mob List All
2004.8 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nItemKey bigint
AS
SET NOCOUNT ON
SELECT nMobID, nGrade, nKillCount FROM tItemMobList WHERE nItemKey = @nItemKey
-- end

0먁ꠁŗ
GO
