CREATE PROCEDURE [dbo].[usp_Gamble_getExchangedCoin]
	@nCharNo int
,	@dDate   datetime
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT ISNULL(SUM(nExchangedCoin), 0) AS nExchangedCoin
	  FROM tGambleCoinExchangeLog 
	 WHERE nCharNo = @nCharNo 
	   AND dDate   >= @dDate 
	   AND nStatus = 1
END

0
GO
