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
