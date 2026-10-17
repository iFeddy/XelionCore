CREATE PROCEDURE [dbo].[usp_OPT_searchUser]
	@key   int,
	@value nvarchar(100)
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @STMT       AS NVARCHAR(1000)
	DECLARE @SQL_WHERE	AS NVARCHAR(200)

	SET @SQL_WHERE = ' WHERE ' +
	CASE @key
		WHEN 1 THEN ' sUserName = ''' + @value + ''''
		WHEN 3 THEN ' nAuthID = ' + CAST(@value AS NVARCHAR(100))
		WHEN 5 THEN ' bIsBlock = ' + CAST(@value AS NVARCHAR(100))
	END
	SET @STMT = 'SELECT * FROM tUser(NOLOCK)' + @SQL_WHERE
	EXEC sp_executesql @STMT
	
	SET NOCOUNT OFF
END
GO
