CREATE PROCEDURE [dbo].[usp_Character_initKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN
		
		DELETE FROM tCharacterKeyMap WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2723
			ROLLBACK TRAN
			RETURN
		END
		
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, sys_create_date, nExtendedKey, nASCIICode)
		SELECT @nCharNo, nFunctionNo, GETDATE(), nExtendedKey, nASCIICode
		  FROM tKeyMapInit
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2724
			ROLLBACK TRAN
			RETURN
		END
		
		COMMIT TRAN
	END
	ELSE BEGIN
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, sys_create_date, nExtendedKey, nASCIICode)
		SELECT @nCharNo, nFunctionNo, GETDATE(), nExtendedKey, nASCIICode
		  FROM tKeyMapInit
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2724
			RETURN
		END
	END
END

0영쮂2
GO
