CREATE PROCEDURE [dbo].[usp_Item_getBuyList]
	@userNo    INT = 0,
	@pageNo    INT = 1,
	@pageLimit INT = 8
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @STMT AS NVARCHAR(1000)
	DECLARE @PARAMS AS NVARCHAR(100)
	DECLARE @TOP AS INT
	DECLARE @START AS INT

	SET @TOP = ((@pageNo - 1) * @pageLimit) + 1

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT @MIN = ISNULL(MIN(orderNo), 0) '
	SET @STMT = @STMT + '  FROM (	SELECT TOP ' + CAST(@TOP AS NVARCHAR) + ' orderNo '
	SET @STMT = @STMT + '			  FROM tOrder WITH (NOLOCK) '
	SET @STMT = @STMT + '			WHERE userNo = ' + CAST(@userNo AS NVARCHAR(10))
	SET @STMT = @STMT + '			ORDER BY orderNo DESC '
	SET @STMT = @STMT + ') AS T1'
	SET @PARAMS = N'@MIN AS INT OUTPUT'
	EXEC sp_executesql @STMT, @PARAMS, @MIN = @START OUTPUT

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR(2)) + ' CONVERT(NVARCHAR(10), orderDate, 120) AS orderDate, orderType, T1.goodsNo, name, imageName, unit, unitDescription, orderPrice '
	SET @STMT = @STMT + '  FROM tOrder AS T1, tItem AS T2 '
	SET @STMT = @STMT + 'WHERE userNo = ' + CAST(@userNo AS NVARCHAR(10)) +  ' AND isPay = 1 AND T1.goodsNo = T2.goodsNo AND orderNo <= ' + CAST(@START AS NVARCHAR(10))
	SET @STMT = @STMT + ' ORDER BY orderNo DESC'
	EXEC sp_executesql @STMT

	SET NOCOUNT OFF
END

0
GO
