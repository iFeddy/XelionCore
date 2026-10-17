CREATE PROCEDURE [dbo].[usp_Character_getItemActionCoolTime]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND dEndTime < GETDATE()) BEGIN
		BEGIN TRAN
		
		DELETE FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND dEndTime < GETDATE()
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -3294
			ROLLBACK TRAN
			RETURN
		END	
		
		SELECT nCharNo, nActionGroup, dEndTime FROM tItemActionCoolTime WHERE nCharNo = @nCharNo
		
		COMMIT TRAN
	END
	ELSE BEGIN
		SELECT nCharNo, nActionGroup, dEndTime FROM tItemActionCoolTime WHERE nCharNo = @nCharNo 
	END
	
	SET NOCOUNT OFF
END

0먋/
GO
