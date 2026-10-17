CREATE PROCEDURE [dbo].[usp_ChatColor_get]
	@nCharNo         int
,	@nRet            int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		SELECT nFontColorID, nBalloonColorID 
		  FROM tChatColor 
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SELECT 0 AS nFontColorID, 0 AS nBalloonColorID
	END
END

0㬁玝N
GO
