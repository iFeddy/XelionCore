CREATE PROCEDURE [dbo].[usp_Dice_setDiceResult] 
	@nGameNo        int
,	@nCustomCharNo  int
,	@nTotalDividend bigint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @money bigint, @resultMoney bigint

	SET @nRet  = 0
	SET @money = 0

	IF ISNULL(@nGameNo, 0) <= 0 BEGIN
		SET @nRet = -2414
		RETURN
	END

	IF ISNULL(@nCustomCharNo, 0) <= 0 BEGIN
		SET @nRet = -2415
		RETURN
	END

	BEGIN TRAN

	SELECT @money = ISNULL(nMoney, 0) FROM tCharacter WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCustomCharNo
	IF ( @money + @nTotalDividend ) > 0 BEGIN
		SET @resultMoney = @money + @nTotalDividend
	END
	ELSE BEGIN
		SET @resultMoney = 0
	END
	
	UPDATE tCharacter SET nMoney = @resultMoney
	WHERE nCharNo = @nCustomCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT > 1 BEGIN
		SET @nRet = - 2493
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tDiceResult (nGameNo, nCustomCharNo, nBeforeMoney, nTotalDividend, nAfterMoney)
	VALUES (@nGameNo, @nCustomCharNo, @money, @nTotalDividend, @resultMoney)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = - 2492
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SELECT nCharNo AS nCustomCharNo, nMoney FROM tCharacter WHERE nCharNo = @nCustomCharNo
END
GO
