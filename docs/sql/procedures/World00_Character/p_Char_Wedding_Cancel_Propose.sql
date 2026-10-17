CREATE PROCEDURE [dbo].[p_Char_Wedding_Cancel_Propose]
	@CharANo INT,
	@CharBNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	SELECT @nRet = 0

-- to Character A

	BEGIN TRAN

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharANo AND nPartnerCharNo = @CharBNo)
	BEGIN
		UPDATE tMarriage
		SET
			nStatus = 0,
			nSubStatus = 0
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
			nStatus = 0,
			nSubStatus = 0
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
