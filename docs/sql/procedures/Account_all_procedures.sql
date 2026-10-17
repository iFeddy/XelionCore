-- ========================================================================
-- Stored Procedures Oficiales de Fiesta Online - Base de Datos: Account
-- Extraídos directamente de Account.bak
-- Total procedimientos recuperados: 77 de 77
-- ========================================================================

-- ========================================================================
-- Procedure: p_UserGetAuthList
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserGetAuthList]
AS
BEGIN
	SET NOCOUNT ON

	SELECT nAuthID, sAuthName, bIsLoginAble  FROM tUserAuth

	SET NOCOUNT OFF
END


0
GO

-- ========================================================================
-- Procedure: p_UserSetAuthLoginable
-- ========================================================================
CREATE PROCEDURE [dbo].[p_UserSetAuthLoginable]
	@userNo   int
,	@AuthID   tinyint
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	
	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			UPDATE tUser 
			   SET nAuthID = @AuthID
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				RETURN (-9)
			ELSE
				RETURN (0)
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END


0
GO

-- ========================================================================
-- Procedure: usp_Cash_getCash
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Cash_getCash]
	@userNo   int
,	@cashType tinyint = 1
,	@CASH	  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @CASH = 0
	SELECT @CASH = ISNULL(cash, 0) FROM tCash 
	 WHERE userNo = @userNo AND cashType = @cashType AND status = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Cash_setCash
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Cash_setCash]
	@userNo    int
,	@cashType  tinyint
,	@cash      int
,	@reference nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	BEGIN TRAN

	IF EXISTS(SELECT userNo FROM tCash WHERE userNo = @userNo AND cashType = @cashType)
		BEGIN
			UPDATE tCash SET cash = cash + @cash WHERE userNo = @userNo AND cashType = @cashType
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-1)
				END
		END
	ELSE
		BEGIN
			INSERT INTO tCash (userNo, cashType, cash)
			VALUES (@userNo, @cashType, @cash)
			IF @@ERROR <> 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-1)
				END
		END

	INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
	VALUES (@userNo, 1, @cashType, @cash, @reference)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		BEGIN
			ROLLBACK TRAN
			RETURN (-1)
		END
	ELSE
		BEGIN
			COMMIT TRAN
			RETURN (0)
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Cash_setCash_EVENT_L20
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Cash_setCash_EVENT_L20]
	@userNo    int
,	@cashType  tinyint
,	@cash      int
,	@reference nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT rowNo FROM tCashLog WHERE userNo = @userNo AND cashType = @cashType AND reference = '200708_EVENT_L20')
		RETURN(-1)
	ELSE
		BEGIN
			BEGIN TRAN

			IF EXISTS(SELECT userNo FROM tCash WHERE userNo = @userNo AND cashType = @cashType)
				BEGIN
					UPDATE tCash SET cash = cash + @cash WHERE userNo = @userNo AND cashType = @cashType
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						BEGIN
							ROLLBACK TRAN
							RETURN (-1)
						END
				END
			ELSE
				BEGIN
					INSERT INTO tCash (userNo, cashType, cash)
					VALUES (@userNo, @cashType, @cash)
					IF @@ERROR <> 0
						BEGIN
							ROLLBACK TRAN
							RETURN (-1)
						END
				END

			INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
			VALUES (@userNo, 1, @cashType, @cash, @reference)
			IF @@ERROR <> 0 OR @@IDENTITY = 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-1)
				END
			ELSE
				BEGIN
					COMMIT TRAN
					RETURN (0)
				END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Cash_useCash
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Cash_useCash]
	@userNo    int
,	@cashType  tinyint
,	@cash      int
,	@reference nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @REMAIN_CASH INT
	SET @REMAIN_CASH = 0

	SELECT @REMAIN_CASH = ISNULL(cash, 0) FROM tCash WHERE userNo = @userNo AND cashType = @cashType AND status = 0
	IF @REMAIN_CASH >= @cash
		BEGIN
			BEGIN TRAN

			UPDATE tCash SET cash = @REMAIN_CASH - @cash WHERE userNo = @userNo AND cashType = @cashType
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			INSERT INTO tCashLog (userNo, logType, cashType, cash, reference)
			VALUES (@userNo, 2, @cashType, @cash, @reference)
			IF @@ERROR <> 0 OR @@IDENTITY = 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END
			ELSE
				BEGIN
					COMMIT TRAN
					RETURN (0)
				END
		END	
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Charge_ItemDel
-- ========================================================================
CREATE  PROCEDURE [dbo].[usp_Charge_ItemDel]
	@rowNo INT,
	@userNo INT,
	@RET	INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	DECLARE @IS_DRAW AS BIT
	--Check User and rowNo
	IF EXISTS(SELECT rowNo FROM tChargeItem WHERE rowNo = @rowNo AND userNo = @userNo)
		BEGIN
			-- check used
			SELECT @IS_DRAW = ISNULL(isDraw,1) FROM tChargeItem WHERE rowNo = @rowNo
			IF @IS_DRAW = 0
				BEGIN
					BEGIN TRAN
					-- insert Log
					INSERT INTO tChargeItemDelLog (rowNo, userNo, orderNo, goodsNo, amount, registerDate, DeleteDate)
					SELECT rowNo, userNo, orderNo, goodsNo, amount, registerDate, getdate()
						FROM tChargeItem
						WHERE @rowNo = rowNo AND @userNo = userNo
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						BEGIN
							-- db Error
							SET @RET = 0
							ROLLBACK TRAN
						END
					ELSE
						BEGIN
							-- item Delete
							DELETE tChargeItem WHERE rowNo = @rowNo AND userNo = @userNo
							IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
								BEGIN
									-- db Error
									SET @RET = 0
									ROLLBACK TRAN
								END
							ELSE
								BEGIN
								-- success
								SET @RET = 1
								COMMIT TRAN
								END
						END
				END
			ELSE
				-- used item
				SET @RET = -1
		END		
	ELSE
		-- mismatch rowNo / userNo
		SET @RET = -2

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
0
GO

-- ========================================================================
-- Procedure: usp_Charge_ItemDraw
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Charge_ItemDraw]
	@rowNo INT,
	@userNo INT,
	@RET	INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @IS_DRAW AS BIT

	-- check user
	IF EXISTS(SELECT rowNo FROM tChargeItem WHERE rowNo = @rowNo AND userNo = @userNo)
		BEGIN
			-- check used
			SELECT @IS_DRAW = isDraw FROM tChargeItem WHERE rowNo = @rowNo
			IF @IS_DRAW = 0
				BEGIN
					BEGIN TRAN

					-- item use
					UPDATE tChargeItem SET isDraw = 1, drawDate = GETDATE() WHERE rowNo = @rowNo
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						BEGIN
							-- db Error
							SET @RET = 0
							ROLLBACK TRAN
						END
					ELSE
						BEGIN
							-- success
							SET @RET = 1
							COMMIT TRAN
						END
				END
			ELSE
				-- used item
				SET @RET = -1
		END		
	ELSE
		-- mismatch rowNo / userNo
		SET @RET = -2

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Charge_ItemInsert
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Charge_ItemInsert]
	@userNo  INT,
	@orderNo BIGINT,
	@goodsNo INT,
	@amount	 INT,
	@RET	 INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF @orderNo <> 0
		BEGIN
			IF EXISTS(SELECT orderNo FROM tChargeItem WITH(NOLOCK) WHERE orderNo = @orderNo)
				BEGIN
					SET @RET = -1
					RETURN
				END
		END
		
	INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount)
	VALUES (@userNo, @orderNo, @goodsNo, @amount)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		SET @RET = 0
	ELSE
		SET @RET = 1	

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Charge_ItemSelectList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Charge_ItemSelectList]
	@userNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT rowNo, goodsNo, amount, registerDate
	FROM tChargeItem WITH (NOLOCK)
	WHERE userNo = @userNo AND isDraw = 0
    ORDER BY rowNo ASC

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getBuyCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getBuyCount]
	@userNo INT = 0,
	@CNT INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CNT = ISNULL(COUNT(orderNo), 0)
	  FROM tOrder WITH (NOLOCK)
	WHERE userNo = @userNo AND isPay = 1

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_Item_getBuyList
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_Item_getCategoryName
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getCategoryName]
	@category NVARCHAR(10) = '',
	@CATEGORY_NAME NVARCHAR(10) OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CATEGORY_NAME = name FROM tItemCategory WITH(NOLOCK) WHERE category = @category

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getEventList
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_Item_getGroup
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getGroup]
	@groupNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT category, name, imageName, itemType, description 
	  FROM VT_ITEM_GROUP WITH (NOLOCK) 
	 WHERE groupNo = @groupNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getGroupList
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_Item_getGroupPrice
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getGroupPrice]
	@groupNo INT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TODAY AS NVARCHAR(10)
	SET @TODAY = CONVERT(NVARCHAR(10), GETDATE(), 120)

	SELECT goodsNo, unit, unitDescription, price, discount, sellPrice = 
		CASE
			WHEN discountStartDate <= @TODAY AND discountEndDate >= @TODAY THEN price - ((price * discount) / 100) 
			ELSE price 
		END
	  FROM VT_ITEM_SELL
	WHERE groupNo = @groupNo
	ORDER BY goodsNo

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_Item_getItem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getItem]
	@goodsNo INT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TODAY AS NVARCHAR(10)
	SET @TODAY = CONVERT(NVARCHAR(10), GETDATE(), 120)

	SELECT goodsNo, billingNo, category, groupNo, name, imageName, itemType, unit, unitDescription, price, discount, sellPrice = 
		CASE
			WHEN discountStartDate <= @TODAY AND discountEndDate >= @TODAY THEN price - ((price * discount) / 100) 
			ELSE price 
		END, isSell, description 
	  FROM tItem WITH (NOLOCK)
	WHERE goodsNo = @goodsNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getOrder
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getOrder]
	@orderNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderDate, orderType, userNo, goodsNo, orderPrice, cashType, isPay, userIP
	  FROM tOrder WITH (NOLOCK)
	WHERE orderNo = @orderNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getOrderList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getOrderList]
	@userNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderDate, orderType, T1.goodsNo, name, imageName, orderPrice
	  FROM tOrder AS T1, tItem AS T2
	WHERE userNo = @userNo AND isPay = 1 AND T1.goodsNo = T2.goodsNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getSubCategory
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getSubCategory]
	@parent NVARCHAR(10)
AS
BEGIN
	SET NOCOUNT ON

	SELECT * FROM tItemCategory WITH(NOLOCK) WHERE parent = @parent

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_getTotalCount
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_Item_getUser
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_getUser]
	@userNo    INT = 0,
	@serialKey INT = 0,
	@USER_ID   NVARCHAR(20) = '' OUTPUT,
	@AUTH_ID   INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @serialKey = 0
		BEGIN
			SET @USER_ID = ''
			SET @AUTH_ID = 0
		END
	ELSE
		BEGIN
			IF EXISTS(SELECT nUserNo FROM tUser WITH (NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
				BEGIN
					SELECT @USER_ID = sUserID, @AUTH_ID = nAuthID
					  FROM tUser AS T1 INNER JOIN tSerialKey AS T2 ON T1.nUserNo = T2.userNo
					 WHERE T2.userNo = @userNo AND T2.serialKey = @serialKey

					SET @USER_ID = ISNULL(@USER_ID, '')
					IF @USER_ID <> '' OR @USER_ID <> NULL
						DELETE FROM tSerialKey WHERE userNo = @userNo
				END
			ELSE
				BEGIN
					SET @USER_ID = ''
					SET @AUTH_ID = 0
				END
		END

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_setBuyOrder
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setBuyOrder]
	@orderNo    INT,
	@userNo     INT,
	@targetUser INT = NULL,
	@RET        INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	-- ÁÖ¹®°ü·Ã º¯¼ö
	DECLARE @ORDER_TYPE AS INT
	DECLARE @ORDER_USER AS NVARCHAR(20)
	DECLARE @GOODS_NO AS INT
	DECLARE @AMOUNT AS INT
	DECLARE @ORDER_PRICE AS INT
	DECLARE @CASH_TYPE AS TINYINT
	DECLARE @IS_PAY AS INT
	DECLARE @OWNER AS INT

	IF EXISTS(SELECT nUserNo FROM tUser WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN

			IF EXISTS(SELECT orderNo FROM tOrder WHERE orderNo = @orderNo)
				BEGIN
					SELECT @ORDER_TYPE = orderType, 
						   @ORDER_USER = userNo, 
						   @GOODS_NO = goodsNo, 
						   @AMOUNT = amount, 
						   @ORDER_PRICE = orderPrice, 
						   @CASH_TYPE = cashType,
						   @IS_PAY = isPay
					  FROM tOrder
					WHERE orderNo = @orderNo

					IF @ORDER_TYPE = 1 AND (@targetUser = 0 or @targetUser IS NULL)
						BEGIN
							SET @RET = -4
							GOTO GT_OUT
						END

					IF @ORDER_TYPE = 0
						SET @OWNER = @userNo
					ELSE
						SET @OWNER = @targetUser

					IF @ORDER_USER <> @userNo
						BEGIN
							SET @RET = -5
							GOTO GT_OUT
						END

					IF @IS_PAY = 1
						BEGIN
							SET @RET = -6
							GOTO GT_OUT
						END
				END
			ELSE
				BEGIN
					SET @RET = -7
					GOTO GT_OUT
				END

			UPDATE tOrder SET toUserNo = @OWNER, isPay = 1
			WHERE orderNo = @orderNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @RET = -9
					GOTO GT_OUT
				END

			EXEC @RET = usp_Cash_useCash @userNo, @CASH_TYPE, @ORDER_PRICE, @orderNo
			IF @RET <> 0
				BEGIN
					SET @RET = @RET
					GOTO GT_OUT
				END

			INSERT INTO tChargeItem (userNo, orderNo, goodsNo, amount)
			VALUES (@OWNER, @orderNo, @GOODS_NO, @AMOUNT)
			IF @@ERROR <> 0 OR @@IDENTITY = 0
				BEGIN
					SET @RET = -9
					GOTO GT_OUT
				END

			UPDATE tItem SET sellCount = sellCount + 1 WHERE goodsNo = @GOODS_NO
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @RET = -9
					GOTO GT_OUT
				END
			ELSE
				BEGIN
					SET @RET = 0
					GOTO GT_OUT
				END
		END
	ELSE
		SET @RET = -2

GT_OUT:
	IF @RET <> 0
		ROLLBACK TRAN
	ELSE
		COMMIT TRAN

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_setCategory
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setCategory]
	@category NVARCHAR(10),
	@name NVARCHAR(20),
	@parent NVARCHAR(10),
	@status BIT,
	@RET INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT category FROM tItemCategory WITH(NOLOCK) WHERE category = @parent)
		BEGIN
			IF EXISTS(SELECT category FROM tItemCategory WITH(NOLOCK) WHERE category = @category)
				BEGIN
					UPDATE tItemCategory 
						SET name = @name,
							parent = @parent,
							status = @status
					WHERE category = @category
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						SET @RET = -2
					ELSE
						SET @RET = 1
				END
			ELSE
				BEGIN
					INSERT INTO tItemCategory (category, name, parent)
					VALUES (@category, @name, @parent)
					IF @@ERROR <> 0
						SET @RET = -1
					ELSE
						SET @RET = 1
				END
		END
	ELSE
		SET @RET = -3

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_setItem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setItem]
	@goodsNo INT,
	@category NVARCHAR(10),
	@groupNo INT,
	@name NVARCHAR(100),
	@imageName NVARCHAR(100),
	@itemType NVARCHAR(20),
	@unit INT,
	@price INT,
	@discount INT,
	@isRecommand BIT,
	@isPopular BIT,
	@isNew BIT,
	@isSell BIT,
	@sellStartDate NVARCHAR(10) = '',
	@sellEndDate NVARCHAR(10) = '',
	@discountStartDate NVARCHAR(10) = '',
	@discountEndDate NVARCHAR(10) = '',
	@RET INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT goodsNo FROM tItem WHERE goodsNo = @goodsNo)
		BEGIN
			BEGIN TRAN

			UPDATE tItem SET
				category = @category,
				groupNo = @groupNo,
				name = @name,
				imageName = @imageName,
				itemType = @itemType,
				unit = @unit,
				price = @price,
				discount = @discount,
				isRecommand = @isRecommand,
				isPopular = @isPopular,
				isNew = @isNew,
				isSell = @isSell,
				sellStartDate = @sellStartDate,
				sellEndDate = @sellEndDate,
				discountStartDate = @discountStartDate,
				discountEndDate = @discountEndDate
			WHERE goodsNo = @goodsNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @RET = -1
					ROLLBACK TRAN
				END
			ELSE
				BEGIN
					SET @RET = 1
					COMMIT TRAN
				END
		END
	ELSE
		SET @RET = -2

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_Item_setOrder
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setOrder]
	@orderType TINYINT,
	@userNo INT,
	@goodsNo INT,
	@goodsPrice INT,
	@amount INT,
	@discount INT,
	@orderPrice INT,
	@cashType int,
	@userIP NVARCHAR(30),
	@ORDER_NO INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @pcbangNo AS INT

	BEGIN TRAN


	INSERT INTO tOrder (orderType, userNo, goodsNo, goodsPrice, amount, discount, orderPrice, cashType, userIP, pcbangNo)
	VALUES (@orderType, @userNo, @goodsNo, @goodsPrice, @amount, @discount, @orderPrice, @cashType, @userIP, @pcbangNo)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		BEGIN
			SET @ORDER_NO = 0
			ROLLBACK TRAN
		END
	ELSE
		BEGIN
			SET @ORDER_NO = @@IDENTITY
			COMMIT TRAN
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Item_setSerialKey
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Item_setSerialKey]
	@userNo    int = 0,
	@serialKey bigint = 0,
	@RET       int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF @userNo = 0 OR @serialKey = 0
		RETURN (-1)

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN 

			IF EXISTS(SELECT userNo FROM tSerialKey WHERE userNo = @userNo)
				DELETE FROM tSerialKey WHERE userNo = @userNo

			INSERT INTO tSerialKey (userNo, serialKey)
			VALUES (@userNo, @serialKey)
			IF @@ERROR <> 0
				BEGIN 
					ROLLBACK TRAN
					SET @RET = 0
				END
			ELSE
				BEGIN
					COMMIT TRAN
					SET @RET = 1
				END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_getBuyList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_getBuyList]
	@userNo   int
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderNo, orderDate, T1.goodsNo, name, unit, unitDescription, toUserNo, 
		   orderType = CASE orderType WHEN 0 THEN 'Buy' ELSE 'Gift' END, orderPrice,
		   cashType  = CASE cashType  WHEN 1 THEN 'Game' ELSE 'Event' END
	  FROM tOrder AS T1, tItem AS T2
	 WHERE userNo = @userNo AND isPay = 1 AND T1.goodsNo = T2.goodsNo
	 ORDER BY orderNo DESC

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_getPayAgreement
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_getPayAgreement]
	@userNo  int
,	@getType int = 0
AS
BEGIN
	SET NOCOUNT ON

	IF @getType = 0
		SELECT rowNo, userNo, name, juminNo, phoneNo, status, issueDate, vaildDate
          FROM tPayAgreement WITH(NOLOCK)
         WHERE userNo = @userNo
		 ORDER BY rowNo DESC
	ELSE
		SELECT TOP 1 rowNo, userNo, name, juminNo, phoneNo, status, issueDate, vaildDate
          FROM tPayAgreement WITH(NOLOCK)
         WHERE userNo = @userNo
		 ORDER BY rowNo DESC

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_getUserCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_getUserCount]
	@KEYFIELD AS NVARCHAR(20) = '',
	@KEYWORD  AS NVARCHAR(20) = ''
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	DECLARE @STMT AS NVARCHAR(1000)
	
	IF @KEYFIELD <> '' AND @KEYWORD <> ''
		SET @STMT = 'SELECT COUNT(nUserNo) FROM tUser(NOLOCK) WHERE ' + @KEYFIELD + ' LIKE ' + '''%' + @KEYWORD + '%'''
	ELSE
		SET @STMT = 'SELECT COUNT(nUserNo) FROM tUser(NOLOCK)'	
	
	EXEC sp_executesql @STMT
	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_getUserID
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_getUserID]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT sUserID FROM tUser WHERE nUserNo = @userNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_getUserNo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_getUserNo]
	@userID nvarchar(60)
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo FROM tUser WHERE sUserID = @userID

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_searchUser
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_OPT_setPayAgreement
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_setPayAgreement]
	@rowNo  int
,	@userNo int
,	@status tinyint
,	@RET    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @USER_NO INT

	SET @RET = 0
	IF EXISTS(SELECT rowNo FROM tPayAgreement WHERE rowNo = @rowNo)
		BEGIN
			SELECT @USER_NO = userNo FROM tPayAgreement WHERE rowNo = @rowNo
			IF @USER_NO <> @userNo
				SET @RET = -2

			IF @RET = 0
				BEGIN
					BEGIN TRAN

					UPDATE tPayAgreement SET status = @status WHERE rowNo = @rowNo
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						BEGIN
							SET @RET = -9
							ROLLBACK TRAN
						END
					ELSE
						BEGIN
							SET @RET = 0
							COMMIT TRAN
						END
				END
		END
	ELSE
		SET @RET = -1

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_setUserAuthID
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_setUserAuthID]
	@userNo   int
,	@AuthID   tinyint
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	
	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			UPDATE tUser 
			   SET nAuthID = @AuthID
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				RETURN (-9)
			ELSE
				RETURN (0)
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_OPT_setUserBlock
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_OPT_setUserBlock]
	@userNo   int
,	@block   bit
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	
	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			UPDATE tUser 
			   SET bIsBlock = @block
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				RETURN (-9)
			ELSE
				RETURN (0)
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Pay_chkLimit
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_chkLimit]
	@userNo  int
,	@chkType int
,	@amount  int
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TOTAL_AMOUNT INT

	IF EXISTS(SELECT orderNo FROM VT_PAY WHERE userNo = @userNo)
		BEGIN
			IF @chkType = 1
				BEGIN
					SELECT @TOTAL_AMOUNT = ISNULL(SUM(cash), 0) 
					  FROM VT_PAY 
					 WHERE userNo = @userNo AND resultDate > CONVERT(NVARCHAR(10), GETDATE(), 120) AND resultDate <= GETDATE()
					IF @TOTAL_AMOUNT >= @amount
						RETURN (1)
					ELSE
						RETURN (0)
				END
			ELSE
				RETURN(0)
		END
	ELSE
		RETURN (0)

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Pay_getAgreement
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_getAgreement]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT TOP 1 status
	  FROM tPayAgreement WITH(NOLOCK) 
	 WHERE userNo = @userNo AND vaildDate >= GETDATE()
	 ORDER BY rowNo DESC

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Pay_getOrder
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_getOrder]
	@orderNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.orderNo, T1.userNo, T1.payNo, T1.cash, T1.userIP, T1.registerDate,
		   T2.result, T2.resultDate, T2.pgTID
	  FROM tPay T1 LEFT JOIN tPayResult T2 ON T1.orderNo = T2.orderNo
	 WHERE T1.orderNo = @orderNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Pay_setAgreement
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_setAgreement]
	@userNo    int
,	@name      nvarchar(10)
,	@juminNo   nvarchar(20)
,	@phoneNo   nvarchar(20)
,	@vaildDate int
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tPayAgreement (userNo, name, juminNo, phoneNo, vaildDate)
	VALUES (@userNo, @name, @juminNo, @phoneNo, DATEADD(mm, @vaildDate, GETDATE()))

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_Pay_setOrder
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_setOrder]
	@userNo int
,	@payNo  int
,	@cash   int
,	@userIP nvarchar(30)
,	@ORDER_NO int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tPay (userNo, payNo, cash, userIP)
	VALUES (@userNo, @payNo, @cash, @userIP)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		SET @ORDER_NO = 0
	ELSE
		SET @ORDER_NO = @@IDENTITY

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_Pay_setResult
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Pay_setResult]
	@orderNo int
,	@result  int
,	@pgTID   nvarchar(50) = ''
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @USER_NO AS INT
	DECLARE @CASH    AS INT
	DECLARE @RET     AS INT

	IF EXISTS(SELECT orderNo FROM tPayResult WHERE orderNo = @orderNo)
		RETURN (-1)
	ELSE
		BEGIN
			BEGIN TRAN
			
			IF @pgTID = ''
				SET @pgTID = NULL

			INSERT INTO tPayResult (orderNo, result, pgTID)
			VALUES (@orderNo, @result, @pgTID)
			IF @@ERROR <> 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END
				
			IF @result = 0
				BEGIN
					SELECT @USER_NO = userNo, @CASH = cash FROM tPay WHERE orderNo = @orderNo

					EXEC @RET = usp_Cash_setCash @USER_NO, 1, @CASH, @orderNo
					IF @RET <> 0
						BEGIN
							ROLLBACK TRAN
							RETURN (@RET)
						END
				END

			COMMIT TRAN
			RETURN (0)
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_ReturnUser_GetDatediffLastlog
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ReturnUser_GetDatediffLastlog]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN 	
	
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT TOP 1 DATEDIFF(s, dDate, getdate()) AS diffDate
	  FROM tReturnUser
	 WHERE nUserNo = @nUserNo
	 ORDER BY dDate DESC
	 
	SET NOCOUNT OFF
    
END

0送亾ţ
GO

-- ========================================================================
-- Procedure: usp_ReturnUser_SetUserItem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_ReturnUser_SetUserItem]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @itemNo int
	DECLARE @RET	int

	SET @itemNo = 30001 -- 이벤트 아이템은 변경될 수 있음
	SET @nRet = 0
	
	BEGIN TRAN

	EXEC usp_Charge_ItemInsert @nUserNo, 0, @itemNo, 1, @RET OUTPUT
	IF @RET = 1 BEGIN
		INSERT INTO tReturnUser (dDate, nUserNo, nItemNo)
		VALUES (GETDATE(), @nUserNo, @itemNo)

		IF @@TRANCOUNT > 0 COMMIT TRAN
	END
	ELSE BEGIN
		IF @@TRANCOUNT > 0 ROLLBACK TRAN
		SET @nRet = -1
	END
				
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
	
END

0줁䋢Ť
GO

-- ========================================================================
-- Procedure: usp_Stat_sellitem
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Stat_sellitem]
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @start AS NVARCHAR(10), @end AS NVARCHAR(10)

	SET @start = CONVERT(NVARCHAR(10), DATEADD(dd, -1, GETDATE()), 120)
	SET @end = CONVERT(NVARCHAR(10), GETDATE(), 120)

	INSERT INTO tbSellitem (baseDate, goodsNo, amount)
	SELECT CONVERT(NVARCHAR(10), registerDate, 120) AS baseDate
	      , goodsNo AS goodsNo
	      , SUM(amount) AS amount
	  FROM tChargeItem WITH(NOLOCK)
	 WHERE registerDate >= @start AND registerDate < @end
	 GROUP BY CONVERT(NVARCHAR(10), registerDate, 120),goodsNo
	 ORDER BY CONVERT(NVARCHAR(10), registerDate, 120),goodsNo

END

0
GO

-- ========================================================================
-- Procedure: usp_Tutorial_getCharInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_getCharInfo]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState, nStep FROM tTutorial WHERE nCharNo = @nCharNo 
END
GO

-- ========================================================================
-- Procedure: usp_Tutorial_getIsSkip
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_getIsSkip]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState FROM tTutorial WHERE nUserNo = @nUserNo AND nState = 1
END

0琁὏ŧ
GO

-- ========================================================================
-- Procedure: usp_Tutorial_setCharInfo
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_Tutorial_setCharInfo]
	@nUserNo int
,	@nCharNo int
,	@nState  tinyint
,	@nStep   tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tTutorial WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tTutorial SET nState = @nState, nStep  = @nStep 
		 WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep)
		VALUES (@nUserNo, @nCharNo, @nState, @nStep)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1
			RETURN
		END		
	END
END

0괁፳Ũ
GO

-- ========================================================================
-- Procedure: usp_User_GetDatediffCreated
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_GetDatediffCreated]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN

	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT DATEDIFF(s, dDate, getdate()) AS diffDate
	  FROM tUser
	 WHERE nUserNo = @nUserNo
	 
    SET NOCOUNT OFF
    
END

0鼁o
GO

-- ========================================================================
-- Procedure: usp_User_USloginGame
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_USloginGame]
	@userID      nvarchar(60)
,	@userPW      nvarchar(32)
,	@userNo      int = 0 OUTPUT
,	@authID      tinyint = 0 OUTPUT
,	@block     	 int = 0 OUTPUT
,	@isLoginable int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET INT
	
	SET @userNo = 0

	/*-------------------------------------------------------------------------------------
		Project: FiestaHeroes - NA2016
		Author: xkl
		
		Using the "r_" parameter, automatically create the account if it doesn't exist.
	------------------------------------------------------------------------------------*/
	IF EXISTS (SELECT Enabled FROM Options..Auto_Account_Creation WHERE Enabled = 1)
	BEGIN
		IF NOT EXISTS (SELECT sUserID FROM tUser WHERE sUserID = @userID)
		BEGIN
			-- If the username starts with "r_", trim it.
			IF @userID LIKE 'r[_]%'
			BEGIN
				SET @userID = SUBSTRING(@userID, 3, LEN(@userID) - 2);
	
				-- If the trimmed username doesn't already exist, insert it into the table.
				IF NOT EXISTS (SELECT sUserID FROM tUser WHERE sUserID = @userID)
				BEGIN
					INSERT INTO tUser(sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP)
					VALUES (@userID, @userPW, @userID, 0, 0, 1, '127.0.0.1')
				END
			END
			ELSE
			BEGIN
				-- The username doesn't start with "r_", so exit.
				RETURN;
			END
		END
	END
	/*-------------------------------------------------------------------------------------
		Continue as normal ~ 
	------------------------------------------------------------------------------------*/

	
	SELECT @userNo = nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
	  FROM tUser AS U(NOLOCK) INNER JOIN tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0
	IF @userNo = 0 BEGIN
		SELECT @userNo = U.nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
		  FROM tUser AS U(NOLOCK) 
		 INNER JOIN tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
		 INNER JOIN tUserOption AS UO(NOLOCK) ON U.nUserNo = UO.nUserNo 
		 WHERE UO.sEmail = @userID AND sUserPW = @userPW AND bIsDelete = 0
		IF @userNo = 0 BEGIN
			SET @userNo = -1
			SET @authID = 0
			SET @block  = 0
			RETURN
		END
	END	
	
	IF @block = 1 BEGIN
		EXEC usp_User_chkBlock @userNo, @RET OUTPUT
		IF @RET = 1 BEGIN
			SELECT @block = bIsBlock FROM tUser(NOLOCK) WHERE nUserNo = @userNo
		END
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_User_USloginGame_O
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_USloginGame_O]
	@userID      nvarchar(60)
,	@userPW      nvarchar(32)
,	@userNo      int = 0 OUTPUT
,	@authID      tinyint = 0 OUTPUT
,	@block     	 int = 0 OUTPUT
,	@isLoginable int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET INT	
	SET @userNo = 0
	
	SELECT @userNo = nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
	  FROM dbo.tUser AS U(NOLOCK) INNER JOIN dbo.tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0
	IF @userNo = 0
		BEGIN
			SET @userNo = -1
			SET @authID = 0
			SET @block  = 0
		END
	ELSE
		BEGIN
			IF @block = 1
			BEGIN
				EXEC usp_User_chkBlock @userNo, @RET OUTPUT
				IF @RET = 1
					SELECT @block = bIsBlock FROM tUser(NOLOCK) WHERE nUserNo = @userNo
			END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_chkBlock
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_chkBlock]
	@userNo INT
,	@RET    INT OUTPUT
AS
BEGIN TRAN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @BLOCK_NO INT
	DECLARE @TYPE     BIT
	DECLARE @END_DATE DATETIME	
	
	SELECT TOP 1 @BLOCK_NO = nBlockNo, @TYPE = nType, @END_DATE = dEndDate 
	  FROM tUserBlock(NOLOCK) 
	 WHERE nUserNo = @userNo AND bStatus = 1
	 ORDER BY nBlockNo
	IF @TYPE = 1 AND @END_DATE < GETDATE()
		BEGIN
			UPDATE tUserBlock SET bStatus = 0 WHERE nBlockNo = @BLOCK_NO
			IF @@ERROR <> 0 OR @@ROWCOUNT = 0
				BEGIN
					ROLLBACK TRAN
					SET @RET = 0
					RETURN
				END
		
			-- Ãß°¡ °èÁ¤ºí·°»óÅÂ À¯¹«Ã¼Å©
			IF EXISTS(SELECT nBlockNo FROM tUserBlock(NOLOCK) WHERE nUserNo = @userNo AND bStatus = 1)
				SET @RET = 0
			ELSE
				BEGIN
					UPDATE tUser SET bIsBlock = 0 WHERE nUserNo = @userNo
					IF @@ERROR <> 0 OR @@ROWCOUNT = 0
						BEGIN
							ROLLBACK TRAN
							SET @RET = 0
							RETURN
						END
			
					SET @RET = 1
				END	
		END
	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
	
COMMIT TRAN

0
GO

-- ========================================================================
-- Procedure: usp_User_chkRegister
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_chkRegister]
	@juminNo    nvarchar(20) = ''
,	@age        int = 0 OUTPUT
,	@registDate nvarchar(10) = '' OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @USER_NO  AS INT
	DECLARE @BIRTHDAY AS NVARCHAR(10)

	IF @juminNo = ''
		RETURN -1

	SET @BIRTHDAY = CASE SUBSTRING(@juminNo, 8, 1)
						WHEN '1' THEN '19'
						WHEN '2' THEN '19'
						WHEN '3' THEN '20'
						WHEN '4' THEN '20'
						ELSE ''
					END
	IF @BIRTHDAY = ''
		RETURN -1
	ELSE
		SET @BIRTHDAY = @BIRTHDAY + SUBSTRING(@juminNo, 1, 2) + '-' + SUBSTRING(@juminNo, 3, 2) + '-' + SUBSTRING(@juminNo, 5, 2)

	SET @age = DATEDIFF(yy, @BIRTHDAY, GETDATE())
	IF @age < 15
		RETURN -2

	SELECT @USER_NO = T1.nUserNo, @registDate = T1.dDate 
	  FROM tUser T1, tUserOption T2
	 WHERE T1.nUserNo = T2.nUserNo
	   AND T2.sJuminNo = @juminNo
	   AND T1.bIsDelete = 0
	IF @USER_NO	> 100
		BEGIN
			SELECT @registDate = dDate FROM tUser WITH(NOLOCK) WHERE nUserNo = @USER_NO
			RETURN -3
		END
	ELSE
		RETURN 0

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_User_chkUserID
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_chkUserID]
	@userID nvarchar(60) = ''
AS
BEGIN
	SET NOCOUNT ON

	IF @userID = '' OR LEN(@userID) < 5 OR LEN(@userID) > 16
		RETURN -1

	IF EXISTS (SELECT sUserID FROM tUser WITH(NOLOCK) WHERE sUserID = @userID AND bIsDelete = 0)
		RETURN -2
	ELSE
		RETURN 0

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_User_delete
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_delete]
	@userNo  int
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN

			UPDATE tUser
			   SET sUserID = '~' + sUserID + '_' + CONVERT(NVARCHAR(8), GETDATE(), 112),
			       bIsDelete = 1
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			UPDATE tUserOption 
			   SET sJuminNo = '~' + sJuminNo + '_' + CONVERT(NVARCHAR(8), GETDATE(), 112)
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END
			ELSE
				BEGIN
					COMMIT TRAN
					RETURN (0)
				END	
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getAge
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getAge]
	@userNo int
,	@AGE    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @JUMIN_NO AS NVARCHAR(14)
	DECLARE @BIRTHDAY AS NVARCHAR(10)

	IF EXISTS(SELECT sJuminNo FROM tUserOption WITH(NOLOCK) WHERE nUserNo = @userNo)
		BEGIN
			SELECT @JUMIN_NO = sJuminNo FROM tUserOption WITH(NOLOCK) WHERE nUserNo = @userNo

			SET @BIRTHDAY = CASE SUBSTRING(@JUMIN_NO, 8, 1)
								WHEN '1' THEN '19'
								WHEN '2' THEN '19'
								WHEN '3' THEN '20'
								WHEN '4' THEN '20'
								ELSE ''
							END
			IF @BIRTHDAY = ''
				SET @AGE = 0
			ELSE
				BEGIN
					SET @BIRTHDAY = @BIRTHDAY + SUBSTRING(@JUMIN_NO, 1, 2) + '-' + SUBSTRING(@JUMIN_NO, 3, 2) + '-' + SUBSTRING(@JUMIN_NO, 5, 2)
					SET @AGE = DATEDIFF(yy, @BIRTHDAY, GETDATE())
				END
		END
	ELSE
		SET @AGE = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getAuthList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getAuthList]
AS
BEGIN
	SET NOCOUNT ON

	SELECT nAuthID, sAuthName, bIsLoginAble  FROM tUserAuth

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getBuyCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getBuyCount]
	@userNo INT = 0,
	@CNT    INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CNT = ISNULL(COUNT(orderNo), 0)
	  FROM tOrder WITH (NOLOCK)
	 WHERE userNo = @userNo AND isPay = 1

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getBuyList
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getBuyList]
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
	SET @STMT = @STMT + '			WHERE userNo = @userNo AND isPay = 1 '
	SET @STMT = @STMT + '			ORDER BY orderNo DESC '
	SET @STMT = @STMT + ') AS T1'
	SET @PARAMS = N'@userNo AS INT, @MIN AS INT OUTPUT'
	EXEC sp_executesql @STMT, @PARAMS, @userNo, @MIN = @START OUTPUT

	SET @STMT = ''
	SET @STMT = @STMT + 'SELECT TOP ' + CAST(@pageLimit AS NVARCHAR(2)) + ' CONVERT(NVARCHAR(10), orderDate, 120) AS orderDate, orderType, T1.goodsNo, name, imageName, unit, unitDescription, orderPrice '
	SET @STMT = @STMT + '  FROM tOrder AS T1, tItem AS T2 '
	SET @STMT = @STMT + ' WHERE T1.orderNo <= @orderNo AND T1.userNo = @userNo AND T1.isPay = 1 AND T1.goodsNo = T2.goodsNo '
	SET @STMT = @STMT + ' ORDER BY T1.orderNo DESC'
	SET @PARAMS = N'@userNo AS INT, @orderNo AS INT'
	EXEC sp_executesql @STMT, @PARAMS, @userNo, @START

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getLogin
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getLogin]
	@userID   NVARCHAR(60)
,	@userPW   NVARCHAR(32)
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.nUserNo, sUserID, sUserPW, sUserName, bIsBlock, nAuthID, sEmail, sJuminNo, sPhoneNo
	  FROM tUser AS T1 LEFT JOIN tUserOption AS T2 ON T1.nUserNo = T2.nUserNo
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getOption
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getOption]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.sEmail, T1.sJuminNo, T1.sPhoneNo, T1.nZipNo, T2.sZIPCODE, T2.sSIDO, T2.sGUGUN, T2.sDONG, T1.sAddress, T1.bIsMail
	  FROM tUserOption T1, tZipcode T2
	 WHERE T1.nZipNo = T2.nZipNo
	   AND T1.nUserNo = @userNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getPayCount
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getPayCount]
	@userNo INT = 0,
	@CNT    INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CNT = ISNULL(COUNT(orderNo), 0)
	  FROM VT_PAY WITH (NOLOCK)
	 WHERE userNo = @userNo

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_getPayList
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_User_getUser
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_getUser]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate
	  FROM tUser WITH(NOLOCK)
	 WHERE nUserNo = @userNo

	SET NOCOUNT OFF
END
GO

-- ========================================================================
-- Procedure: usp_User_insert
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_insert]
	@userID   nvarchar(20)
,	@userPW   nvarchar(32)
,	@nAffiliateID varchar(50) = NULL
,	@szAffilateSubID varchar(50) = NULL
,	@szCountryCode varchar(50) = NULL
,	@userName nvarchar(10)
,	@userIP   nvarchar(30)
,	@eMail    nvarchar(70) = NULL
,	@juminNo  nvarchar(20) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
,	@channel  nvarchar(20) = 'ONSON'
,	@isOption bit = 1
,	@userNo   int = 0 OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT
	DECLARE @itemNo AS INT

	SET @itemNo = 30000 -- 이벤트 아이템은 변경될 수 있음

	EXEC @RET = usp_User_chkUserID @userID
	IF (@RET <> 0)
		RETURN (-2)

	BEGIN TRAN

	INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP)
	VALUES (@userID, @userPW, @userName, @userIP)
	IF (@@ERROR <> 0 OR @@IDENTITY = 0)
		BEGIN
			ROLLBACK TRAN
			RETURN (-9)
		END
	ELSE
		BEGIN                        
			SET @userNo = @@IDENTITY

			IF (@isOption = 1)
				BEGIN
					EXEC @RET = usp_User_insertOption @userNo,  @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel, @nAffiliateID
					IF (@RET <> 0)
						BEGIN
							ROLLBACK TRAN
							RETURN (-3)
						END
				END
	END

	SET @RET = 0
	EXEC usp_Charge_ItemInsert @userNo, 0, @itemNo, 1, @RET OUTPUT

	COMMIT TRAN
	SET LOCK_TIMEOUT -1
	SET NOCOUNT ON
END

0Ё憓Ŷ
GO

-- ========================================================================
-- Procedure: usp_User_insertOption
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_insertOption]
	@userNo  int
,	@eMail   nvarchar(256) = NULL
,	@juminNo nvarchar(20) = NULL
,	@phoneNo nvarchar(20) = NULL
,	@zipNo   int = NULL
,	@address nvarchar(90) = NULL
,	@isMail  bit = NULL
,	@channel nvarchar(20) = NULL
,	@nAffiliateID int = NULL
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel, nAffiliateID)
	VALUES (@userNo, @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel, @nAffiliateID)
	IF @@ERROR <> 0
		RETURN (-1)
	ELSE
		RETURN (0)
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_insertOption_Old
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_insertOption_Old]
	@userNo  int
,	@eMail   nvarchar(70) = NULL
,	@juminNo nvarchar(20) = NULL
,	@phoneNo nvarchar(20) = NULL
,	@zipNo   int = NULL
,	@address nvarchar(90) = NULL
,	@isMail  bit = NULL
,	@channel nvarchar(20) = NULL
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	INSERT INTO tUserOption (nUserNo, sEmail, sJuminNo, sPhoneNo, nZipNo, sAddress, bIsMail, sChannel)
	VALUES (@userNo, @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel)
	IF @@ERROR <> 0
		RETURN (-1)
	ELSE
		RETURN (0)
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_insert_Old
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_insert_Old]
	@userID   nvarchar(60)
,	@userPW   nvarchar(32)
,	@userName nvarchar(10)
,	@userIP   nvarchar(30)
,	@eMail    nvarchar(70) = NULL
,	@juminNo  nvarchar(20) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
,	@channel  nvarchar(20) = 'ONSON'
,	@isOption bit
,	@userNo   int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT

	EXEC @RET = usp_User_chkUserID @userID
	IF (@RET <> 0)
		RETURN (-2)

	BEGIN TRAN

	INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP)
	VALUES (@userID, @userPW, @userName, @userIP)
	IF (@@ERROR <> 0 OR @@IDENTITY = 0)
		BEGIN
			ROLLBACK TRAN
			RETURN (-9)
		END
	ELSE
		BEGIN			
			SET @userNo = @@IDENTITY

			IF (@isOption = 1)
				BEGIN
					EXEC @RET = usp_User_insertOption @userNo, @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel
					IF (@RET <> 0)
						BEGIN
							ROLLBACK TRAN
							RETURN (-3)
						END
				END

			COMMIT TRAN
			RETURN (0)
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT ON
END
GO

-- ========================================================================
-- Procedure: usp_User_loginGame
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_loginGame]
	@userID      nvarchar(60)
,	@userPW      nvarchar(32)
,	@userNo      int = 0 OUTPUT
,	@authID      tinyint = 0 OUTPUT
,	@block     	 int = 0 OUTPUT
,	@isLoginable int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET INT	
	SET @userNo = 0
	
	SELECT @userNo = nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
	  FROM dbo.tUser AS U(NOLOCK) INNER JOIN dbo.tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0
	IF @userNo = 0
		BEGIN
			SET @userNo = -1
			SET @authID = 0
			SET @block  = 0
		END
	ELSE
		BEGIN
			IF @block = 1
			BEGIN
				EXEC usp_User_chkBlock @userNo, @RET OUTPUT
				IF @RET = 1
					SELECT @block = bIsBlock FROM tUser(NOLOCK) WHERE nUserNo = @userNo
			END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_search
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_search]
	@userName nvarchar(10)
,	@juminNo  nvarchar(20)
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.nUserNo, T1.sUserID, T1.sUserPW, T2.sEmail 
	  FROM tUser T1, tUserOption T2
	 WHERE T1.nUserNo = T2.nUserNo
	   AND T1.sUserName = @userName
	   AND T2.sJuminNo = @juminNo
	   AND T1.bIsDelete = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_searchZipcode
-- ========================================================================
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

-- ========================================================================
-- Procedure: usp_User_setChannelLog
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_setChannelLog]
	@userNo    int
,	@channel   nvarchar(20)
,	@issueDate nvarchar(20)
As
BEGIN
	SET NOCOUNT ON

	INSERT INTO tChannelLog (userNo, channel, issueDate)
	VALUES (@userNo, @channel, @issueDate)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		RETURN (0)
	ELSE
		RETURN (-1)

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_setDrawComment
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_setDrawComment]
	@userNo  int
,	@causeNo tinyint = 0
,	@comment nvarchar(1000) = ''
,	@userIP  nvarchar(30)
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tDrawComment (nUserNo, nCauseNo, sComment, sUserIP)
	VALUES (@userNo, @causeNo, @comment, @userIP)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		RETURN (-9)
	ELSE
		RETURN (0)

	SET NOCOUNT OFF
END

!bNôÙ¾ dF(ä¿u PF+
GO

-- ========================================================================
-- Procedure: usp_User_setLoginable
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_setLoginable]
	@authID      int
,	@isLoginable int
,	@RET         int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nAuthID FROM tUserAuth WHERE nAuthID = @authID)
		BEGIN
			UPDATE tUserAuth SET bIsLoginAble = @isLoginable WHERE nAuthID = @authID
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				SET @RET = 0
			ELSE
				SET @RET = 1
		END
	ELSE
		SET @RET = 0

	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_update
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_update]
	@userNo   int
,	@eMail    nvarchar(200) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN

			UPDATE tUserOption 
			   SET sEmail = @eMail, sPhoneNo = @phoneNo, nZipNo = @zipNo, sAddress = @address, bIsMail = @isMail
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END
			ELSE
				BEGIN
					COMMIT TRAN
					RETURN (0)
				END
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO

-- ========================================================================
-- Procedure: usp_User_updatePassword
-- ========================================================================
CREATE PROCEDURE [dbo].[usp_User_updatePassword]
	@userNo int
,	@prevPW nvarchar(32)
,	@userPW nvarchar(32)
,	@userIP nvarchar(30)
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @CURR_PW AS NVARCHAR(20)

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			SELECT @CURR_PW = sUserPW FROM tUser WHERE nUserNo = @userNo
			If @CURR_PW <> @prevPW
				RETURN (-2)

			BEGIN TRAN

			UPDATE tUser SET sUserPW = @userPW WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			INSERT INTO tPasswordLog (nUserNo, sPrevPW, sCurrPW, sUserIP)
			VALUES (@userNo, @CURR_PW, @userPW, @userIP)
			IF @@ERROR <> 0 OR @@IDENTITY = 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			COMMIT TRAN
			RETURN (0)
		END
	ELSE
		RETURN (-1)
END

<
GO

