CREATE PROCEDURE [dbo].[usp_SafeTrade_setMerge]
	@nBuyer            int,
	@nBuyerItemKey     bigint,
	@nBuyerOptionType  smallint,
	@nSeller           int,
	@nSellerItemKey    bigint,
	@nSellerOptionType smallint,
	@nMoveMoney        bigint,
	@nSplitLot         bigint,
	@bIsDelete         bit,
	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nBuyerOptionData AS bigint
	DECLARE @nSellerOptionData AS bigint
	DECLARE @nSellerLot AS bit
	SET @ret = 0
	SET @nBuyerOptionData = 0
	SET @nSellerOptionData = 0
	SET @nSellerLot = 0

	BEGIN TRAN
	
	-- 구매자의 돈을 판매자에게 이동 시킴
	EXEC p_Char_MoneyMove @nBuyer, @nSeller, @nMoveMoney, @ret output
	IF @ret = 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2623
		RETURN
	END	
	
	-- 아이템을 구매자의 인벤토리로 이동시킴
	EXEC p_Item_Split_N_Merge_New @nSellerItemKey, @nSeller, @nSellerOptionType, @nSellerOptionData output,
								  @nBuyerItemKey, @nBuyer, @nBuyerOptionType, @nBuyerOptionData output, @nSplitLot, @ret output
	IF @ret <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2624
		RETURN
	END

	-- 판매자 아이템의 Lot 개수가 0 이 아닌경우 1 로 만들어 판매자의 아이템이 삭제됨을 방지한다.
	IF @nSellerOptionData > 0 BEGIN
		SET @nSellerLot = 1
	END

	IF @bIsDelete = @nSellerLot BEGIN
		IF @bIsDelete = 0 BEGIN
			-- Lot이 0인데 삭제 요청이 오지 않은 경우
			ROLLBACK TRAN
			SET @nRet = -2625
		END
		ELSE BEGIN
			-- Lot이 0이 아닌데 삭제 요청이 온 경우
			ROLLBACK TRAN
			SET @nRet = -2626
		END
	END
	ELSE BEGIN
		IF @bIsDelete = 1 BEGIN
			-- 판매자의 아이템삭제
			EXEC p_Item_Delete @nSellerItemKey, @ret output
			IF @ret <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -2627
			END
		END

		COMMIT TRAN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0ᨁ姐'
GO
