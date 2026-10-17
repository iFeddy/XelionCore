CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_Cancel]
	@CharANo INT,
	@CharBNo INT,
	@IncreseCount TINYINT,
	@SubStatusA TINYINT OUTPUT,
	@SubStatusB TINYINT OUTPUT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		SELECT @SubStatusA = nSubStatus FROM tMarriage WHERE nCharNo = @CharANo
		UPDATE tMarriage
		SET
			nSubStatus =
			CASE nSubStatus
				WHEN 1 THEN 1 + 4
				WHEN 3 THEN 3 + 4
				ELSE nSubStatus
			END,
			nDivorceCancel = nDivorceCancel + @IncreseCount,
			dSubWhen = dWhen
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		SELECT @SubStatusB = nSubStatus FROM tMarriage WHERE nCharNo = @CharBNo
		UPDATE tMarriage
		SET
			nSubStatus = 0,
			dSubWhen = dWhen
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

	COMMIT TRAN
	GOTO END_

ERROR_:
	ROLLBACK TRAN
	SELECT @nRet = 1

END_:
	SET NOCOUNT OFF

END

0
GO
