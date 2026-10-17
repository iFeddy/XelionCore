CREATE PROCEDURE [dbo].[usp_Character_chkKeyMap]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tCharacterKeyMap WHERE nCharNo = @nCharNo) BEGIN
		SELECT 1
	END
	ELSE BEGIN
		SELECT 0
	END
END

0㐁ᙁĪ
GO
