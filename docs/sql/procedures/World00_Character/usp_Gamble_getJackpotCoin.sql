CREATE PROCEDURE [dbo].[usp_Gamble_getJackpotCoin] 
	@sGameName    nvarchar(20)
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nBattingUnit, nJackpotCoin 
	  FROM tGambleJackpot 
	 WHERE sGameName = @sGameName 
END
GO
