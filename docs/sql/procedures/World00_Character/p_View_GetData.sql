CREATE PROCEDURE [dbo].[p_View_GetData]
	@QueryExpr nvarchar(512)

AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END

0轿
GO
