CREATE PROCEDURE [dbo].[usp_Collection_setCard]
	@nCharNo      int
,	@nCardID      int
,	@nSerialNo    bigint
,	@nStarNumber  tinyint
,	@nItemKey     bigint
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
	IF ISNULL(@nCardID, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END
	IF ISNULL(@nSerialNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END
	IF ISNULL(@nStarNumber, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN	
	END

	IF EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nItemKey) BEGIN
		BEGIN TRAN
		
		-- 1. 컬렉션 카드 등록
		IF EXISTS(SELECT nCharNo FROM tCollectionCard WHERE nCharNo = @nCharNo AND nCardID = @nCardID) BEGIN
			-- 1-1. 이미 존재하는 경우
			UPDATE tCollectionCard 
			   SET nSerialNo   = @nSerialNo 
				 , nStarNumber = @nStarNumber 
			 WHERE nCharNo = @nCharNo 
			   AND nCardID = @nCardID 
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		ELSE BEGIN
			-- 1-2. 새로 등록하는 경우
			INSERT INTO tCollectionCard (nCharNo, nCardID, nSerialNo, nStarNumber)
			VALUES (@nCharNo, @nCardID, @nSerialNo, @nStarNumber)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -9999
				ROLLBACK TRAN
				RETURN
			END
		END
		
		-- 2. 컬렉션에 등록된 카드의 경우 tItem 테이블에서 해당 아이템을 삭제
		DELETE FROM tItem WHERE nItemKey = @nItemKey
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -9999
			ROLLBACK TRAN
			RETURN
		END	
		
		COMMIT TRAN		
	END
	ELSE BEGIN
		SET @nRet = -9999
		RETURN
	END
END

0쨁ₚU
GO
