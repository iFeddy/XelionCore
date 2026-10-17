CREATE PROCEDURE [dbo].[usp_SafeTrade_setLot]
	@nBuyer            int,
	@nBuyerItemKey     bigint,
	@nBuyerStorageType tinyint,
	@nBuyerStorage     smallint,
	@nBuyerItemID      int,
	@nBuyerFlags       int,
	@nBuyerOptionType  smallint,
	@nSeller           int,
	@nSellerItemKey    bigint,
	@nSellerOptionType smallint,
	@nMoveMoney        bigint,
	@nSplitLot         bigint,
	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nBuyerOptionData AS bigint
	DECLARE @nSellerOptionData AS bigint
	DECLARE @bIsDelete AS bit
	SET @ret = 0
	SET @nBuyerOptionData = 0
	SET @nSellerOptionData = 0
	SET @bIsDelete = 0

	BEGIN TRAN
	
	-- 구매자에게 판매자의 아이템을 구입할 수 있도록 아이템을 생성시킨다. 단, Lot 은 0개로 초기화
	INSERT INTO tItem (nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags)
	VALUES (@nBuyerItemKey, @nBuyer, @nBuyerStorageType, @nBuyerStorage, @nBuyerItemID, @nBuyerFlags)
	IF @@ERROR <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = -2628
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nBuyerItemKey, @nBuyerOptionType, @nBuyerOptionData)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2629
			RETURN
		END	
	END
	
	-- 아이템을 구매자의 인벤토리로 이동시킴
	EXEC usp_SafeTrade_setMerge @nBuyer, @nBuyerItemKey, @nBuyerOptionType, @nSeller, @nSellerItemKey,
								@nSellerOptionType, @nMoveMoney, @nSplitLot, @bIsDelete, @ret output
	IF @ret <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = @ret
		RETURN
	END
	ELSE BEGIN
		COMMIT TRAN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END
GO
