CREATE PROCEDURE [dbo].[usp_Gamble_subCoin]
	@nCharNo  int
,	@nCoin    bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @coin bigint

	SET @status = 3 -- ??	-> ???
	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nCoin, 0) < 0 BEGIN
		SET @nRet = -3313
		RETURN
	END
	
	-- ??? ??? 0 ?? ??? ?? ??? ????? ??.
	IF @nCoin = 0 BEGIN
		SET @nRet = 0
		RETURN
	END	
		
	SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF @coin < @nCoin BEGIN
		SET @nRet = -3322
		RETURN		
	END
	ELSE BEGIN
		BEGIN TRAN
		
		UPDATE tGambleCoin 
		   SET sys_update_date = GETDATE()
			 , nCoin = nCoin - @nCoin 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3323
			RETURN
		END

		INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin)
		VALUES (@nCharNo, @nCoin, @status, @coin, @coin - @nCoin)
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
END

0
GO
