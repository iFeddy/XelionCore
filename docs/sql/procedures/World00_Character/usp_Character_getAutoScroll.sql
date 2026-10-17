CREATE PROCEDURE [dbo].[usp_Character_getAutoScroll]
	@nCharNo		int
,	@nRet			int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	SELECT nAutoBuff, nAutoChargedbuff
	  FROM tCharacterSingleOption
	 WHERE nCharNo= @nCharNo
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0ᤁꁑĪ
GO
