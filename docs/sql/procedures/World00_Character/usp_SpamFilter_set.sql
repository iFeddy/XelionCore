CREATE PROCEDURE [dbo].[usp_SpamFilter_set]
	@nUserNo int
,	@bIsSpamSpread bit
,	@bIsSpamBlock  bit
,	@nSpamFilterRate tinyint
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nUserNo FROM tSpamFilter WHERE nUserNo = @nUserNo)
	BEGIN
		UPDATE tSpamFilter
		   SET sys_update_date = GETDATE()
		     , bIsSpamSpread = @bIsSpamSpread
		     , bIsSpamBlock = @bIsSpamBlock
		     , nSpamFilterRate = @nSpamFilterRate
		 WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
		BEGIN
			SET @nRet = -2893
		END
	END
	ELSE BEGIN
		INSERT INTO tSpamFilter (nUserNo, bIsSpamSpread, bIsSpamBlock, nSpamFilterRate)
		VALUES (@nUserNo, @bIsSpamSpread, @bIsSpamBlock, @nSpamFilterRate)
		IF @@ERROR <> 0
		BEGIN
			SET @nRet = -2892
		END
	END

	SET NOCOUNT ON
END

0︁⩠+
GO
