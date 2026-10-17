CREATE PROCEDURE [dbo].[p_Item_BuyAll]
	@nItemKey bigint,
	@nOwner int,
	@nStorageType tinyint,
	@nStorage smallint,
	@nItemID int,
	@nFlags int,
	@nPrice bigint,
	@nFame int,
	@nCoin bigint,
	@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @CreateItemRet int

	IF @nPrice <= 0 AND @nFame <= 0 AND @nCoin <= 0 BEGIN
		SET @nRet = 2
		RETURN
	END

	BEGIN TRAN
	
	-- ??? ???? ??? ?? ??? ??? ?? ? ? ??
	IF @nPrice > 0 BEGIN
		EXEC p_Char_MoneySub @nOwner, @nPrice, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nFame > 0 BEGIN
		EXEC p_Char_FameSub @nOwner, @nFame, @nRet OUTPUT
		IF @nRet <> 1
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ???? ??? ??? ??? ??? ?? ? ?? ??
	IF @nCoin > 0 BEGIN
		EXEC usp_Gamble_subCoin @nOwner, @nCoin, @nRet OUTPUT
		IF @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = 1
			RETURN
		END
	END
	
	-- ??? ?? ?????? ??
	EXEC p_Item_Create @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags, @CreateItemRet OUTPUT
	IF @CreateItemRet <> 0 OR @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = @CreateItemRet
		RETURN
	END
	-- ??
	COMMIT TRAN
	SET @nRet = 0
	-- end
END

<
GO
