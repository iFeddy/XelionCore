CREATE PROCEDURE [dbo].[usp_Character_getKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		SELECT nFunctionNo, nExtendedKey, nASCIICode
		  FROM tCharacterKeyMap 
		 WHERE nCharNo = @nCharNo 
	END
	ELSE BEGIN
		SET @nRet = -2725
		RETURN
	END
END

0漁太İ
GO
