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
