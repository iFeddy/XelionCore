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
