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
