CREATE PROCEDURE [dbo].[p_Char_Wedding_Propose]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @CharBNo,
			nStatus = 1,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 3, GETDATE())

		WHERE
			nCharNo = @CharANo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharANo, @CharBNo, 1,  0, GETDATE(), DATEADD(day, 3, GETDATE()))		

		 IF @@ERROR <> 0 GOTO ERROR_
	END

-- to Character B

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @CharANo,
			nStatus = 1,
			nSubStatus = 0,
			dWhen = GETDATE(),
			dSubWhen = DATEADD(day, 3, GETDATE())
		WHERE
			nCharNo = @CharBNo

		 IF @@ERROR <> 0 GOTO ERROR_
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharBNo, @CharANo, 1,  0, GETDATE(), DATEADD(day, 3, GETDATE()))		

		 IF @@ERROR <> 0 GOTO ERROR_
	END

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
