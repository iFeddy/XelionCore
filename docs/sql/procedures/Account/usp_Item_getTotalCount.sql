CREATE PROCEDURE [dbo].[usp_Item_getTotalCount]
	@eventType INT = 0,
	@category  NVARCHAR(10) = '',
	@TOTAL     INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @STMT AS NVARCHAR(1000)
	DECLARE @PARAMS AS NVARCHAR(100)

	IF @eventType > 0
		BEGIN
			SET @STMT = 'SELECT @CNT = COUNT(goodsNo) FROM VT_ITEM_SELL '
			IF @eventType = 1
				SET @STMT = @STMT +  'WHERE isNew = 1 '
			ELSE IF @eventType = 2
				SET @STMT = @STMT +  'WHERE isPopular = 1 '
			ELSE IF @eventType = 3
				SET @STMT = @STMT +  'WHERE isRecommand = 1 '
		END
	ELSE
		BEGIN 
			SET @STMT = 'SELECT @CNT = COUNT(goodsNo) FROM VT_ITEM_GROUP '
			IF @category <> ''
				SET @STMT = @STMT + ' WHERE category LIKE ''' + @category + '%'''			
		END
	SET @PARAMS = N'@CNT AS INT OUTPUT'
	EXEC sp_executesql @STMT, @PARAMS, @CNT = @TOTAL OUTPUT

	SET @TOTAL = ISNULL(@TOTAL, 0)

	SET NOCOUNT OFF
END

0
GO
