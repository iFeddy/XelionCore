CREATE PROCEDURE [dbo].[usp_Dice_setDicePay]
	@nCharNo        int
,	@nItemKey       bigint
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF ISNULL(@nTotalDividend, 0) <= 0 BEGIN
		SET @nRet = -2416
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo) BEGIN
		BEGIN TRAN

		-- žÓŽÏÄ«µå ŸÆÀÌÅÛ »èÁŠ
		DELETE FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2494
			ROLLBACK TRAN
			RETURN
		END

		-- µ· Áö±Þ
		UPDATE tCharacter SET nMoney = nMoney + @nTotalDividend WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
			RETURN
		END

		-- µ· Áö±Þ ·Î±×
		INSERT INTO tDicePayLog (nCharNo, nItemKey, nTotalDividend)
		VALUES (@nCharNo, @nItemKey, @nTotalDividend)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2492
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			SET @nRet = 0
			COMMIT TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -2423
		RETURN	
	END
END
GO
