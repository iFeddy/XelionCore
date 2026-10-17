CREATE PROCEDURE [dbo].[p_GameLog_GetLog]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT a.nLogNo, a.dDate, a.nCharNo, a.nType, a.sMap, a.nMapX, a.nMapY, a.nTargetCharNo, a.nTargetID, a.nItemKey, a.nInt1, a.nInt2, a.nInt3, a.nBigint1, b.sID
			  FROM tGameLog a INNER JOIN [W00CHAR].World00_Character.dbo.tCharacter b ON a.nCharNo = b.nCharNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END
GO
