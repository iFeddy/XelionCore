CREATE PROCEDURE [dbo].[usp_Item_getGroupList]
	@category NVARCHAR(10) = '',
	@pageNo INT = 1,
	@pageLimit INT = 6
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @STMT AS NVARCHAR(1000)
	DECLARE @PARAMS AS NVARCHAR(100)
	DECLARE @WHERE AS NVARCHAR(500)
	DECLARE @TODAY AS NVARCHAR(10)
	DECLARE @TOP AS INT
	DECLARE @START AS INT

	SET @TODAY = CONVERT(NVARCHAR(10), GETDATE(), 120)
	SET @TOP = ((@pageNo - 1) * @pageLimit) + 1
	SET @WHERE = ''

	IF @pageNo = 1
		BEGIN
			IF @category <> ''
				SET @WHERE = ' WHERE category LIKE ''' + @category + '%'''

			SET @STMT = ''
			SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR) + ' goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount, sellPrice = CASE '
			SET @STMT = @STMT + '		WHEN discountStartDate <= ''' + @TODAY + ''' AND discountEndDate >= ''' + @TODAY + ''' THEN price - ((price * discount) / 100) ELSE price END '
			SET @STMT = @STMT + '  FROM VT_ITEM_GROUP '
			SET @STMT = @STMT + @WHERE
			EXEC sp_executesql @STMT
		END
	ELSE
		BEGIN
			IF @category <> ''
				SET @WHERE = ' AND category LIKE ''' + @category + '%'''

			SET @STMT = ''
			SET @STMT = @STMT + 'SELECT @MIN = ISNULL(MIN(goodsNo), 0) '
			SET @STMT = @STMT + '  FROM (	SELECT TOP ' + CAST(@TOP AS NVARCHAR) + ' goodsNo '
			SET @STMT = @STMT + '			  FROM VT_ITEM_GROUP '
			SET @STMT = @STMT + '			WHERE isSell = 1 '
			SET @STMT = @STMT + @WHERE
			SET @STMT = @STMT + ') AS T1'
			SET @PARAMS = N'@MIN AS INT OUTPUT'
			EXEC sp_executesql @STMT, @PARAMS, @MIN = @START OUTPUT
		
			SET @STMT = ''
			SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR) + ' goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount, sellPrice = CASE '
			SET @STMT = @STMT + '		WHEN discountStartDate <= ''' + @TODAY + ''' AND discountEndDate >= ''' + @TODAY + ''' THEN price - ((price * discount) / 100) ELSE price END '
			SET @STMT = @STMT + '  FROM VT_ITEM_GROUP '
			SET @STMT = @STMT + 'WHERE goodsNo <= ' + CAST(@START AS NVARCHAR(10)) + @WHERE
			EXEC sp_executesql @STMT
		END

	SET NOCOUNT OFF
END
GO
