CREATE PROCEDURE [dbo].[p_Item_SetMobList]
/*
Item Set Mob List
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nMobID smallint,
@nGrade tinyint,
@nKillCount int,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemMobList SET	nGrade = @nGrade, nKillCount = @nKillCount
WHERE nItemKey = @nItemKey AND nMobID = @nMobID
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemMobList (  nItemKey,  nMobID,  nGrade,  nKillCount )
	VALUES				( @nItemKey, @nMobID, @nGrade, @nKillCount )
END
SET @nRet = @@ERROR
-- end

0舁⓪ū
GO
