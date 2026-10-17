CREATE PROCEDURE [dbo].[usp_Character_setItemActionCoolTime]
	@nCharNo      int
,	@nActionGroup smallint
,	@dEndTime     datetime
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tItemActionCoolTime WHERE nCharNo = @nCharNo AND nActionGroup = @nActionGroup) BEGIN
		BEGIN TRAN
		
		UPDATE tItemActionCoolTime
		   SET dEndTime = @dEndTime 
		 WHERE nCharNo = @nCharNo 
		   AND nActionGroup = @nActionGroup
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -3291
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		INSERT INTO tItemActionCoolTime (nCharNo, nActionGroup, dEndTime)
		VALUES (@nCharNo, @nActionGroup, @dEndTime)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3291
			RETURN
		END
	END	
	
	SET NOCOUNT OFF
END

0༁ᐇ=
GO
