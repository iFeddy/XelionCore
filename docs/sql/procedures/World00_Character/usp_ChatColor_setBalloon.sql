CREATE PROCEDURE [dbo].[usp_ChatColor_setBalloon]
	@nCharNo         int
,	@nBalloonColorID tinyint
,	@nRet            int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nFontColorID AS tinyint

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tChatColor 
		   SET nBalloonColorID = @nBalloonColorID 
		 WHERE nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -1
		END
		ELSE BEGIN
			COMMIT TRAN
		END	
	END
	ELSE BEGIN
		SET @nFontColorID = 0
			
		INSERT INTO tChatColor (nCharNo, nFontColorID, nBalloonColorID)
		VALUES (@nCharNo, @nFontColorID, @nBalloonColorID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
		END
	END
END

0괁寥P
GO
