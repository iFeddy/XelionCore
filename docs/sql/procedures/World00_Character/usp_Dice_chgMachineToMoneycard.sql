CREATE PROCEDURE [dbo].[usp_Dice_chgMachineToMoneycard]
	@nItemKey       bigint
,	@nItemIDOld     int
,	@nItemIDNew     int
,	@nOptionType    smallint
,	@nTotalDividend bigint
,	@nRet              int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @nIsPlace AS smallint 
	DECLARE @nIsPlaceValue AS bigint

	SET @nRet          = 0
	SET @nIsPlace      = 904 -- 가구아이템이 거치된 아이템인지 판단하는 OptionType
	SET @nIsPlaceValue = 0   -- 거치된 아이템을 회수하기 위함 Value

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	IF ISNULL(@nItemIDOLD, 0) <= 0 BEGIN
		SET @nRet = -2417
		RETURN
	END

	IF ISNULL(@nItemIDNEW, 0) <= 0 BEGIN
		SET @nRet = -2418
		RETURN
	END

	IF ISNULL(@nOptionType, 0) <= 0 BEGIN
		SET @nRet = -2419
		RETURN
	END

	IF ISNULL(@nTotalDividend, 0) <= 0 BEGIN
		SET @nRet = -2416
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey AND nItemID = @nItemIDOLD) BEGIN
		BEGIN TRAN

		-- 게임기의 아이템아이디를 머니카드의 아이템아이디로 교체
		UPDATE tItem SET nItemID = @nItemIDNEW
		 WHERE nItemKey = @nItemKey
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
			RETURN
		END

		-- 거치된 게임기의 경우 인벤토리로 회수하는 작업을 함.
		UPDATE tItemOptions SET nOptionData = @nIsPlaceValue
		 WHERE nItemKey = @nItemKey AND nOptionType = @nIsPlace
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
		END

		-- 머니카드에 배당금을 확정함
		UPDATE tItemOptions SET nOptionData = @nTotalDividend
		 WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2493
			ROLLBACK TRAN
		END
		ELSE BEGIN
			SET @nRet = 0
			COMMIT TRAN
		END
	END
	ELSE BEGIN
		SET @nRet = -2422
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0긁X
GO
