CREATE PROCEDURE [dbo].[p_Char_Wedding_Escape_Divorce_State]
	@CharNo INT,
	@nRet INT OUTPUT
AS

BEGIN

	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		SELECT @nRet = 0
		UPDATE tMarriage
		SET
			nStatus = 0,
			nSubStatus = 0,
			nDivorceCancel = 0
		WHERE
			nCharNo = @CharNo
	END
	ELSE
	BEGIN
		SELECT @nRet = 1
	END

	SET NOCOUNT OFF

END

0
GO
