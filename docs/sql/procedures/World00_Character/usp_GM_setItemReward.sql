CREATE PROCEDURE [dbo].[usp_GM_setItemReward]
	@nGroupNo     bigint
,	@nCharNo      int
,	@nStorageType tinyint
,	@nItemID      int
,	@nItemKey     bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nStorage AS smallint

	SET @nStorage = 0 -- nStoragType = 2 번의 nStorage = 0 번에 보상인벤 제일 앞에 인벤토리임.
	SET @nRet = 0
	
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SET @nRet = -1111
		RETURN
	END
		
	IF EXISTS(SELECT nGroupNo FROM tGMItemRewardLog WHERE nGroupNo = @nGroupNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2222
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		UPDATE tItem 
		   SET nStorage = nStorage + 1 
		 WHERE nOwner = @nCharNo AND nStorageType = @nStorageType 
		IF @@ERROR <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = -3333
			RETURN
		END
		ELSE BEGIN
			INSERT INTO tItem (nItemKey, nStorageType, nOwner, nStorage, nItemID)
			VALUES (@nItemKey, @nStorageType, @nCharNo, @nStorage, @nItemID)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -4444
				RETURN
			END	
			
			INSERT INTO tGMItemRewardLog (nGroupNo, nCharNo, nStorageType, nItemID, nItemKey)
			VALUES (@nGroupNo, @nCharNo, @nStorageType, @nItemID, @nItemKey)
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -5555
				RETURN
			END	
			ELSE BEGIN
				COMMIT TRAN
			END			
		END
	END
END

0鐁m
GO
