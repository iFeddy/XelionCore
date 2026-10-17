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
