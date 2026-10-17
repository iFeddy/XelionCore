CREATE PROCEDURE [dbo].[usp_Minihouse_setGate]
	@nCharNo  int
,	@nItemID  int
,	@bType    tinyint
,	@dEndDate datetime = NULL
,	@nItemKey bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		-- 기존에 등록된 게이트가 있을 경우 사용기간이 이미 지났으면 삭제처리를 먼저한다.
		DELETE FROM tMinihouseGate 
		 WHERE nCharNo = @nCharNo 
		   AND nItemID = @nItemID 
		   AND dEndDate < GETDATE()
		   AND dEndDate IS NOT NULL
	END

	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		-- 상위에서 삭제했는데도 있는 경우는 기간이 남았거나 영구아이템일 경우임.
		SET @nRet = -1621
	END
	ELSE BEGIN
		IF @bType = 0 BEGIN
			SET @dEndDate = NULL
		END

		BEGIN TRAN
		-- 게이트 스크롤 아이템을 삭제함.
		EXEC p_Item_Delete @nItemKey, @nRet output
		IF @nRet = 0 BEGIN
			INSERT INTO tMinihouseGate (nCharNo, nItemID, bType, dEndDate) 
			VALUES (@nCharNo, @nItemID, @bType, @dEndDate)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -1692
				ROLLBACK TRAN
			END
			ELSE BEGIN
				COMMIT TRAN
			END
		END
		ELSE BEGIN
			SET @nRet = -1622
			ROLLBACK TRAN
		END		
	END

	SET NOCOUNT OFF
END

0㘁褿#
GO
