CREATE PROCEDURE [dbo].[usp_Character_setChargedbuffScroll]
	@nCharNo			int
,	@nAutoChargedBuff	tinyint
,	@nRet				int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DECLARE @nAutoBuff AS tinyint
	
	SET @nRet = 0
	
	BEGIN TRAN
	
	IF EXISTS (SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		UPDATE tCharacterSingleOption SET nAutoChargedBuff = @nAutoChargedBuff WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		SET @nAutoBuff = 0
	
		INSERT INTO tCharacterSingleOption (nCharNo, nAutoBuff, nAutoChargedBuff)
		VALUES (@nCharNo, @nAutoBuff, @nAutoChargedBuff)
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

0ᰁ葜8
GO
