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
