CREATE PROCEDURE [dbo].[usp_User_getPayList]
	@userNo    INT = 0,
	@pageNo    INT = 1,
	@pageLimit INT = 10
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
	SET @STMT = @STMT + '			  FROM VT_PAY '
	SET @STMT = @STMT + '			 WHERE userNo = @userNo '
	SET @STMT = @STMT + '			 ORDER BY orderNo DESC '
	SET @STMT = @STMT + ') AS T1'
	SET @PARAMS = N'@userNo AS INT, @MIN AS INT OUTPUT'
	EXEC sp_executesql @STMT, @PARAMS, @userNo, @MIN = @START OUTPUT

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR(2)) + ' CONVERT(NVARCHAR(10), resultDate, 120) AS resultDate, T2.name, T1.cash '
	SET @STMT = @STMT + '  FROM VT_PAY AS T1, tPayType AS T2 '
	SET @STMT = @STMT + ' WHERE T1.orderNo <= @orderNo AND userNo = @userNo AND T1.payNo = T2.payNo '
	SET @STMT = @STMT + ' ORDER BY orderNo DESC'
	SET @PARAMS = N'@userNo AS INT, @orderNo AS INT'
	EXEC sp_executesql @STMT, @PARAMS, @userNo, @START

	SET NOCOUNT OFF
END

0
GO
