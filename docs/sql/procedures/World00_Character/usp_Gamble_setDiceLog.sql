CREATE PROCEDURE [dbo].[usp_Gamble_setDiceLog]
	@nGameNo       int 
,	@nCharNo       int
,	@nSlotNo       int
,	@nBattingMoney bigint
,	@nBattingRate  int
,	@bStatus       bit
,	@nTaxRate      int
,	@nDividend     bigint
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	INSERT INTO tGambleDiceLog (nGameNo, nCharNo, nSlotNo, nBattingMoney, nBattingRate, bStatus, nTaxRate, nDividend)
	VALUES (@nGameNo, @nCharNo, @nSlotNo, @nBattingMoney, @nBattingRate, @bStatus, @nTaxRate, @nDividend)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -3392
		RETURN
	END

	SET NOCOUNT OFF
END

0
GO
