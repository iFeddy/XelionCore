CREATE PROCEDURE [dbo].[usp_Character_setKeyMap]
	@nCharNo      int
,	@nFunctionNo  smallint
,	@nExtendedKey tinyint
,	@nASCIICode   tinyint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	BEGIN TRAN
	
	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo AND nFunctionNo = @nFunctionNo) BEGIN
		UPDATE tCharacterKeyMap 
		   SET nExtendedKey = @nExtendedKey 
		     , nASCIICode = @nASCIICode 
		 WHERE nCharNo = @nCharNo 
		   AND nFunctionNo = @nFunctionNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -2793
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tCharacterKeyMap (nCharNo, nFunctionNo, nExtendedKey, nASCIICode)
		VALUES (@nCharNo, @nFunctionNo, @nExtendedKey, @nASCIICode)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2792
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN
END

0㤁䤑=
GO
