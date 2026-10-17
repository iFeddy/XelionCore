CREATE PROCEDURE [p_AccountLog_GetLogCount]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*)
			  FROM tAccountLog  a INNER JOIN AccountDB.Account.dbo.tUser b ON a.nUserNo = b.nUserNo ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END

0෪Ć
GO
