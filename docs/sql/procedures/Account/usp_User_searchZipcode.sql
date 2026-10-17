CREATE PROCEDURE [dbo].[usp_User_searchZipcode]
	@dong VARCHAR(50)
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @STMT AS NVARCHAR(1000)
	
	SET @STMT = 'SELECT * FROM tZipCode WITH(NOLOCK) WHERE sDONG LIKE ''' + @dong + '%'''
	EXEC sp_executesql @STMT

	SET NOCOUNT OFF
END

0
GO
