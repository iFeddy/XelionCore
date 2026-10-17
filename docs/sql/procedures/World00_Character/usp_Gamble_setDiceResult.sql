CREATE PROCEDURE [dbo].[usp_Gamble_setDiceResult] 
	@nGameNo        int
,	@nCharNo        int
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @coin bigint, @resultCoin bigint

	SET @nRet  = 0
	SET @coin = 0
	SET @resultCoin = 0

	IF ISNULL(@nGameNo, 0) <= 0 BEGIN
		SET @nRet = -3315
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -3311
		RETURN
	END

	BEGIN TRAN

	SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
	IF ( @coin + @nTotalDividend ) > 0 BEGIN
		SET @resultCoin = @coin + @nTotalDividend
	END
	ELSE BEGIN
		SET @resultCoin = 0
	END
	
	UPDATE tGambleCoin 
	   SET nCoin = @resultCoin
	     , sys_update_date = GETDATE()
	 WHERE nCharNo = @nCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -3323
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tGambleDiceResult (nGameNo, nCharNo, nBeforeMoney, nTotalDividend, nAfterMoney)
	VALUES (@nGameNo, @nCharNo, @coin, @nTotalDividend, @resultCoin)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = - 3330
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SELECT nCharNo, nCoin FROM tGambleCoin WHERE nCharNo = @nCharNo
END
GO
