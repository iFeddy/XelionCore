CREATE PROCEDURE [dbo].[p_Char_WeddingStatus_Get]
	@CharNo INT
AS
BEGIN
	SET NOCOUNT ON

	IF NOT EXISTS(SELECT * FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, dWhen, dSubWhen) VALUES(@CharNo, 0, GETDATE(), GETDATE())
	END

	SELECT nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen, nMarriageCount, nDivorceCancel FROM tMarriage WHERE nCharNo = @CharNo

	SET NOCOUNT OFF
END

0
GO
