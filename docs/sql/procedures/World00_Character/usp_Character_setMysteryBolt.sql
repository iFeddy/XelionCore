CREATE PROCEDURE [dbo].[usp_Character_setMysteryBolt]
	@nCharNo      int
,	@bMysteryBolt bit
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tCharacterSingleOption (nCharNo, bMysteryBolt)
		VALUES (@nCharNo, @bMysteryBolt)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3192
			RETURN
		END		
	END
	ELSE BEGIN
		UPDATE tCharacterSingleOption 
		   SET bMysteryBolt = @bMysteryBolt 
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3193
			RETURN
		END		
	END
END

0爁㴵>
GO
