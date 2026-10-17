CREATE PROCEDURE [dbo].[p_Char_Wedding_Divorce_by_Force]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A : applicant

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 3,
			dSubWhen = DATEADD(Day, 3, GETDATE())
--			dSubWhen = DATEADD(minute, 10, GETDATE())
		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE GOTO ERROR_

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo AND nPartnerCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nSubStatus = 4,
			dSubWhen = DATEADD(Day, 3, GETDATE())
--			dSubWhen = DATEADD(minute, 10, GETDATE())
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
