CREATE PROCEDURE [dbo].[usp_Item_setAutoMon]
	@nCharNo int
,	@nType	 int
,	@nSlot	 int
,	@nItemID int
,	@nRet    int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	BEGIN TRAN
				
	IF EXISTS (SELECT nCharNo, nType, nSlot FROM tAutoMon WHERE nCharNo = @nCharNo AND nType = @nType AND nSlot = @nSlot) BEGIN
		UPDATE tAutoMon SET nItemID = @nItemID WHERE nCharNo = @nCharNo AND nType = @nType AND nSlot = @nSlot
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tAutoMon (nCharNo, nType, nSlot, nItemID)
		VALUES (@nCharNo, @ntype, @nSlot, @nItemID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
			
	COMMIT TRAN
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF

END

0渁
GO
