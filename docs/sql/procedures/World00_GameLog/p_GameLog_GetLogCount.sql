CREATE PROCEDURE [dbo].[p_GameLog_GetLogCount]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
	
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*)
			  FROM tGameLog a INNER JOIN [W00CHAR].World00_Character.dbo.tCharacter b ON a.nCharNo = b.nCharNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END

0缁攘Q
GO
