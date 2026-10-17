CREATE PROCEDURE [dbo].[usp_Character_setRestExpExecTime]
	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo)
		VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET dRestExpLastExecuteTime = GETDATE() 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0╽@
GO
