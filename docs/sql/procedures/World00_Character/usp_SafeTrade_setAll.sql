CREATE PROCEDURE [dbo].[usp_SafeTrade_setAll]
	@nBuyer       int,
	@nSeller      int,
	@nMoveMoney   bigint,
	@nItemKey     bigint,
	@nStorageType tinyint,
	@nStorage     smallint,	
	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @ret AS int
	DECLARE @nOwner AS int
	DECLARE @nOwnerStorageType AS tinyint
	SET @ret = 0
	SET @nOwner = 0
	SET @nOwnerStorageType = 0

	-- 아이템 소유권 확인(판매자 소유의 아이템인지 확인)
	SELECT @nOwner = nOwner, @nOwnerStorageType = nStorageType FROM tItem WHERE nItemKey = @nItemKey
	IF @nOwner <> @nSeller BEGIN
		SET @nRet = -2621
		RETURN
	END
	IF @nOwnerStorageType <> 9 BEGIN
		SET @nRet = -2622
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		-- 구매자의 돈을 판매자에게 이동 시킴
		EXEC p_Char_MoneyMove @nBuyer, @nSeller, @nMoveMoney, @ret output
		IF @ret = 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2623
			RETURN
		END	
		
		-- 아이템을 구매자의 인벤토리로 이동시킴
		EXEC p_Item_SetStorageAndOwner @nItemKey, @nBuyer, @nStorageType, @nStorage, @ret output
		IF @ret = 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -2624
			RETURN
		END
		ELSE BEGIN
			COMMIT TRAN
			SET @nRet = 0
		END
	END

	SET NOCOUNT OFF
END
GO
