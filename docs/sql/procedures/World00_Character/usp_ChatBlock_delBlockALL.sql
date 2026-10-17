CREATE PROCEDURE [dbo].[usp_ChatBlock_delBlockALL]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF NOT EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo) BEGIN
		SET @nRet = -2321
		RETURN
	END
	ELSE BEGIN
		DELETE FROM tChatBlock WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2394
		END
	END
END

0줁譔L
GO
