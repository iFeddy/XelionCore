CREATE PROCEDURE [dbo].[usp_Gamble_setMoneyToCoin]
	@nCharNo int
,	@nMoney  bigint
,	@nCoin   bigint
,	@nLimit  bigint
,	@dDate   datetime
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @money bigint, @coin bigint, @exchangedCoin bigint
	
	SET @status = 1 -- °ÔÀÓžÓŽÏ -> ÄÚÀÎ
	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	IF ISNULL(@nMoney, 0) <= 0 BEGIN
		SET @nRet = -3312
		RETURN
	END	
	
	IF ISNULL(@nCoin, 0) <= 0 BEGIN
		SET @nRet = -3313
		RETURN
	END	
	
	-- ÀÏÁ€œÃ°£µ¿ŸÈ ±³È¯ÇÑ ÄÚÀÎÀÌ ³ÑÀ» °æ¿ì ±³È¯ÇÒ Œö ŸøÀœ
	SELECT @exchangedCoin = ISNULL(SUM(nExchangedCoin), 0) 
	  FROM tGambleCoinExchangeLog 
	 WHERE nCharNo = @nCharNo 
	   AND dDate  >= @dDate
	   AND nStatus = @status
	IF @nLimit < @exchangedCoin + @nCoin BEGIN
		SET @nRet = -3327
		RETURN	
	END

	BEGIN TRAN
		
	SELECT @money = ISNULL(nMoney, 0) FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @money IS NULL BEGIN
		SET @nRet = -3324
		ROLLBACK TRAN
		RETURN
	END
	
	IF @money < @nMoney BEGIN
		SET @nRet = -3328
		ROLLBACK TRAN
		RETURN
	END
		
	UPDATE tCharacter
	   SET nMoney   = nMoney - @nMoney
	 WHERE nCharNo  = @nCharNo
	   AND bDeleted = 0
	   AND nMoney  >= @nMoney
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3325
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

	INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin, nBeforeMoney, nAfterMoney)
	VALUES (@nCharNo, @nCoin, @status, @coin, @coin + @nCoin, @money, @money - @nMoney)
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
GO
