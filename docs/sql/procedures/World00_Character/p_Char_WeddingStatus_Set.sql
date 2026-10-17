CREATE PROCEDURE [dbo].[p_Char_WeddingStatus_Set]
	@CharNo INT,
	@PartnerCharNo INT,
	@Status SMALLINT,
	@SubStatus SMALLINT,
	@When DATETIME,
	@SubWhen DATETIME,
	@MarriageCount TINYINT,
	@DivorceCancel TINYINT
AS

BEGIN

	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @CharNo)
	BEGIN
		UPDATE tMarriage
		SET
			nPartnerCharNo = @PartnerCharNo,
			nStatus = @Status,
			nSubStatus = @SubStatus,
			dWhen = @When,
			dSubWhen = @SubWhen,
			nMarriageCount = @MarriageCount,
			nDivorceCancel = @DivorceCancel
		WHERE
			nCharNo = @CharNo
	END
	ELSE
	BEGIN
		INSERT tMarriage(nCharNo, nPartnerCharNo, nStatus, nSubStatus, dWhen, dSubWhen)
		VALUES(@CharNo, @PartnerCharNo, @Status, @SubStatus, @When, @SubWhen)
		
	END

	SET NOCOUNT OFF

END

0
GO
