CREATE PROCEDURE [dbo].[usp_ChatBlock_delBlock]
	@nUserNo int
,	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2312
		RETURN
	END

	IF NOT EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2323
		RETURN
	END
	ELSE BEGIN
		DELETE FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2394
		END
	END
END

0送霰K
GO
