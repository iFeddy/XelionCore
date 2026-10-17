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
