CREATE PROCEDURE [p_Report_Get_Count_OPT]
	@QueryExpr nvarchar(256),
	@Count int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT @RowCount = COUNT(*) 
		FROM tReport ' + @QueryExpr
	EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
	SET NOCOUNT OFF
END
GO
