CREATE PROCEDURE [dbo].[p_Item_GetMobList]
/*
Item Get Mob List
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data fount
	        ? = Number data here(must to be 1)
	        
	@nGrade = grade
	@nKillCount = kill count(progress)
		        
*/
@nItemKey bigint,
@nMobID smallint,
-- Output var
@nGrade tinyint OUTPUT,
@nKillCount int OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT @nGrade = nGrade, @nKillCount = nKillCount 
FROM tItemMobList
WHERE nItemKey = @nItemKey AND nMobID = @nMobID
SET @nRet = @@ROWCOUNT
-- end

0픁ᷱŗ
GO
