CREATE PROCEDURE [dbo].[usp_Item_getEventList]
	@eventType TINYINT = 1,
	@pageNo    INT = 1,
	@pageLimit INT = 3
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
	IF @eventType = 1
		SET @WHERE = 'WHERE isNew = 1 '
	ELSE IF @eventType = 2
		SET @WHERE = 'WHERE isPopular = 1 '
	ELSE IF @eventType = 3
		SET @WHERE = 'WHERE isRecommand = 1 '

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT @MIN = ISNULL(MIN(goodsNo), 0) '
	SET @STMT = @STMT + '  FROM (	SELECT TOP ' + CAST(@TOP AS NVARCHAR) + ' goodsNo '
	SET @STMT = @STMT + '			  FROM VT_ITEM_SELL '
	SET @STMT = @STMT + @WHERE
	SET @STMT = @STMT + ') AS T1'
	SET @PARAMS = N'@MIN AS INT OUTPUT'
	EXEC sp_executesql @STMT, @PARAMS, @MIN = @START OUTPUT

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR) + ' goodsNo, groupNo, name, imageName, unit, unitDescription, price, discount, sellPrice = CASE '
	SET @STMT = @STMT + '		WHEN discountStartDate <= ''' + @TODAY + ''' AND discountEndDate >= ''' + @TODAY + ''' THEN price - ((price * discount) / 100) ELSE price END, description '
	SET @STMT = @STMT + '  FROM VT_ITEM_SELL '
	SET @STMT = @STMT + @WHERE + ' AND goodsNo <= ' + CAST(@START AS NVARCHAR(10))
	EXEC sp_executesql @STMT

	SET NOCOUNT OFF
END

0
GO
