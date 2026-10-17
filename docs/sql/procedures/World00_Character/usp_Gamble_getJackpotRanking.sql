CREATE PROCEDURE [dbo].[usp_Gamble_getJackpotRanking] 
	@sGameName    nvarchar(20)
,	@nBattingUnit bigint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT TOP 50 sID, nJackpotCoin  
	  FROM tGambleJackpotLog 
	 WHERE sGameName = @sGameName 
	   AND nBattingUnit = @nBattingUnit 
END

0
GO
