CREATE PROCEDURE [dbo].[usp_SpamBlock_getTime]
	@nCharNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT nCharNo, nBlockTime, nStandbyTime
		  FROM tSpamBlockTime
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SET @nRet = -2822
	END

	SET NOCOUNT ON
END

0匁䷴(
GO
