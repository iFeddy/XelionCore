CREATE PROCEDURE [dbo].[usp_Gamble_addCoin]
	@nCharNo  int
,	@nItemKey bigint
,	@nCoin    bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint, @coin bigint
	
	SET @status = 2 -- À¯·áŸÆÀÌÅÛ -> ÄÚÀÎ
	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nCoin, 0) <= 0 BEGIN
		SET @nRet = -3313
		RETURN
	END
		
	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo) BEGIN
		BEGIN TRAN

		DELETE FROM tItem WHERE nItemKey = @nItemKey AND nOwner = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3332
			RETURN		
		END
		
		SELECT @coin = nCoin FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
		IF @coin IS NULL BEGIN
			INSERT INTO tGambleCoin (nCharNo, nCoin)
			VALUES (@nCharNo, @nCoin)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3329
				RETURN
			END	
			
			SET @coin = 0	
		END
		ELSE BEGIN
			UPDATE tGambleCoin 
			   SET sys_update_date = GETDATE()
				 , nCoin = nCoin + @nCoin 
			 WHERE nCharNo = @nCharNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				ROLLBACK TRAN
				SET @nRet = -3323
				RETURN
			END
		END
		
		INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin)
		VALUES (@nCharNo, @nCoin, @status, @coin, @coin + @nCoin)
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3326
			RETURN
		END
		ELSE BEGIN
			COMMIT TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -3331
		RETURN	
	END
END
GO
