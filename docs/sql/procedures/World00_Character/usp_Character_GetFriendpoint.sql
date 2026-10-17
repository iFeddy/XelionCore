CREATE PROCEDURE [dbo].[usp_Character_GetFriendpoint]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		SELECT ISNULL(nFriendPoint, 0) AS nFriendPoint
		  FROM tCharacter
		 WHERE @nCharNo = nCharNo
	END
	ELSE BEGIN
		SET @nRet = -1111
	END

	SET LOCK_TIMEOUT -1
    SET NOCOUNT OFF    
END

0
GO
