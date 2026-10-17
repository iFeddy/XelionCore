CREATE PROCEDURE [dbo].[usp_Gamble_setCoinToMoney]
	@nCharNo int
,	@nMoney  bigint
,	@nCoin   bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @status tinyint
	DECLARE @money bigint, @coin bigint
	
	SET @status = 0 -- ÄÚÀÎ -> °ÔÀÓžÓŽÏ
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

	BEGIN TRAN
		
	SELECT @coin = nCoin FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF @coin IS NULL BEGIN
		SET @nRet = -3321
		ROLLBACK TRAN
		RETURN
	END
	
	IF @coin < @nCoin BEGIN
		SET @nRet = -3322
		ROLLBACK TRAN
		RETURN
	END
	
	UPDATE tGambleCoin 
	   SET sys_update_date = GETDATE()
		 , nCoin = nCoin - @nCoin 
	 WHERE nCharNo = @nCharNo 
	   AND nCoin  >= @coin
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		ROLLBACK TRAN
		SET @nRet = -3323
		RETURN
	END
	
	SELECT @money = nMoney FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @money IS NULL BEGIN		
		ROLLBACK TRAN
		SET @nRet = -3324
		RETURN
	END
	ELSE BEGIN
		UPDATE tCharacter
		   SET nMoney   = nMoney + @nMoney
		 WHERE nCharNo  = @nCharNo
		   AND bDeleted = 0
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3325
			RETURN		
		END	
	END

	INSERT INTO tGambleCoinExchangeLog (nCharNo, nExchangedCoin, nStatus, nBeforeCoin, nAfterCoin, nBeforeMoney, nAfterMoney)
	VALUES (@nCharNo, @nCoin, @status, @coin, @coin - @nCoin, @money, @money + @nMoney)
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
