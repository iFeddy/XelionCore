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
