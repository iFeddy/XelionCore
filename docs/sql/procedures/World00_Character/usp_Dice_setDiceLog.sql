CREATE PROCEDURE [dbo].[usp_Dice_setDiceLog]
	@nGameNo       int 
,	@nCustomCharNo int
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

	INSERT INTO tDiceLog (nGameNo, nCustomCharNo, nSlotNo, nBattingMoney, nBattingRate, bStatus, nTaxRate, nDividend)
	VALUES (@nGameNo, @nCustomCharNo, @nSlotNo, @nBattingMoney, @nBattingRate, @bStatus, @nTaxRate, @nDividend)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -2492
		RETURN
	END

	SET NOCOUNT OFF
END

0鈁솼\
GO
