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
