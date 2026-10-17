CREATE PROCEDURE [dbo].[usp_Gamble_getCoin]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tGambleCoin WHERE nCharNo = @nCharNo) BEGIN
		SELECT nCoin FROM tGambleCoin WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SELECT 0 AS nCoin
	END
END

0
GO
