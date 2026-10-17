CREATE PROCEDURE [dbo].[usp_SpamFilter_get]
	@nUserNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nUserNo FROM tSpamFilter WHERE nUserNo = @nUserNo)
	BEGIN
		SELECT nUserNo, bIsSpamSpread, bIsSpamBlock, nSpamFilterRate
		  FROM tSpamFilter
		 WHERE nUserNo = @nUserNo
	END
	ELSE BEGIN
		SET @nRet = -2821
	END

	SET NOCOUNT ON
END

0씁㘼*
GO
