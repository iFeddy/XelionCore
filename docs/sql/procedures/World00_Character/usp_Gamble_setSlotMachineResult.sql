CREATE PROCEDURE [dbo].[usp_Gamble_setSlotMachineResult] 
	@nCharNo      int
,	@nBattingUnit bigint
,	@nBatting     tinyint
,	@sGameResult  nvarchar(10)
,	@bIsWin       bit
,	@nPrizeCoin   bigint
,	@bIsJackpot   bit
,	@nJackpotRate int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @gameName nvarchar(20), @charID nvarchar(80)
	DECLARE @gameNo bigint, @battingCoin bigint, @coin bigint, @resultCoin bigint, @jackpotCoin bigint
	
	SET @nRet = 0
	SET @gameName = 'SlotMachine'
	SET @battingCoin = 0
	SET @resultCoin  = 0
	SET @battingCoin = @nBattingUnit * @nBatting 
		
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -3311
		RETURN
	END
	IF ISNULL(@nBattingUnit, 0) = 0 BEGIN
		SET @nRet = -3314
		RETURN		
	END
	IF ISNULL(@nBatting, 0) = 0 BEGIN
		SET @nRet = -3315
		RETURN		
	END	
	IF ISNULL(@sGameResult, '') = '' BEGIN
		SET @nRet = -3316
		RETURN		
	END
	
	SELECT @charID = sID FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0
	IF @charID IS NOT NULL BEGIN
		BEGIN TRAN
		
		-- 1. ???? ??
		INSERT INTO tGambleSlotMachineLog (nCharNo, nBattingUnit, nBatting, sGameResult, bIsWin, nPrizeCoin)
		VALUES (@nCharNo, @nBattingUnit, @nBatting, @sGameResult, @bIsWin, @nPrizeCoin)
		IF @@ERROR <> 0 OR @@IDENTITY = 0 BEGIN
			SET @nRet = -3333
			ROLLBACK TRAN
			RETURN
		END
		SET @gameNo = @@IDENTITY
		
		-- 2. ??( ????(-) / ????(+)
		SELECT @coin = ISNULL(nCoin, 0) FROM tGambleCoin WITH(UPDLOCK, HOLDLOCK) WHERE nCharNo = @nCharNo
		IF (@coin - @battingCoin + @nPrizeCoin) > 0 BEGIN
			SET @resultCoin = @coin - @battingCoin + @nPrizeCoin
		END
		ELSE BEGIN
			SET @resultCoin = 0
		END
		
		UPDATE tGambleCoin 
		   SET nCoin = @resultCoin
			 , sys_update_date = GETDATE()
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -3334
			ROLLBACK TRAN
			RETURN
		END		 
		
		-- 3. ??? ??
		IF @bIsJackpot = 1 BEGIN
			-- 3-1. ?? ???? ??
			SELECT @jackpotCoin = ISNULL(nJackpotCoin, 0) 
			  FROM tGambleJackpot WITH(UPDLOCK, HOLDLOCK)
			 WHERE sGameName = @gameName
			   AND nBattingUnit = @nBattingUnit
			IF @jackpotCoin <> 0 BEGIN
				SET @jackpotCoin = @jackpotCoin * @nJackpotRate / 1000
			END
			
			UPDATE tGambleCoin 
			   SET nCoin = nCoin + @jackpotCoin 
			     , sys_update_date = GETDATE()
			 WHERE nCharNo = @nCharNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3335
				ROLLBACK TRAN
				RETURN
			END
			
			-- 3-2. ?? ???? ???
			UPDATE tGambleJackpot 
			   SET nJackpotCoin = 0
			     , sys_update_date = GETDATE()
			 WHERE sGameName = @gameName 
			   AND nBattingUnit = @nBattingUnit 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3336
				ROLLBACK TRAN
				RETURN
			END
			
			-- 3-3. ???? ??
			INSERT INTO tGambleJackpotLog (sGameName, nBattingUnit, sID, nJackpotCoin, nGameNo)
			VALUES (@gameName, @nBattingUnit, @charID, @jackpotCoin, @gameNo)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3337
				ROLLBACK TRAN
				RETURN
			END						   				  
		END
		-- 4. ??? ????
		ELSE BEGIN
			-- 4-1. ???? ??
			IF EXISTS(SELECT sGameName FROM tGambleJackpot WHERE sGameName = @gameName AND nBattingUnit = @nBattingUnit) BEGIN
				UPDATE tGambleJackpot 
				   SET nJackpotCoin = nJackpotCoin + @battingCoin 
					 , sys_update_date = GETDATE()
				 WHERE sGameName = @gameName
				   AND nBattingUnit = @nBattingUnit 
				IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
					SET @nRet = -3338
					ROLLBACK TRAN
					RETURN
				END	
			END
			ELSE BEGIN
				INSERT INTO tGambleJackpot (sGameName, nBattingUnit, nJackpotCoin)
				VALUES (@gameName, @nBattingUnit, @battingCoin)
				IF @@ERROR <> 0 BEGIN
					SET @nRet = -3338
					ROLLBACK TRAN
					RETURN
				END				
			END		   
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		SET @nRet = -3339
	END
END
GO
