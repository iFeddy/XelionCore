CREATE PROCEDURE [dbo].[usp_Collection_setFavorites]
	@nCharNo      int
,	@nSlotNo      tinyint	
,	@nCardID      int
,	@bRegist      bit
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	IF ISNULL(@nSlotNo, -1) < 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END

	BEGIN TRAN

	IF @bRegist = 1 BEGIN
		IF ISNULL(@nCardID, 0) = 0 BEGIN
			SET @nRet = -1111
			ROLLBACK TRAN
			RETURN	
		END		
		
		-- 1. 마이컬렉션 카드 등록
		IF EXISTS(SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo) BEGIN
			-- 1-1. 이미 존재하는 경우
			UPDATE tCollectionFavorites 
			   SET nCardID = @nCardID 
			 WHERE nCharNo = @nCharNo
			   AND nSlotNo = @nSlotNo 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			-- 1-2. 새로 등록하는 경우
			INSERT INTO tCollectionFavorites (nCharNo, nSlotNo, nCardID)
			VALUES (@nCharNo, @nSlotNo, @nCardID)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
	END
	ELSE BEGIN
		IF EXISTS(SELECT nCharNo FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo) BEGIN
			DELETE FROM tCollectionFavorites WHERE nCharNo = @nCharNo AND nSlotNo = @nSlotNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			SET @nRet = -5555
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
END

0́ᒿV
GO
