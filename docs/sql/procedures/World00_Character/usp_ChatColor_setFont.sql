CREATE PROCEDURE [dbo].[usp_ChatColor_setFont]
	@nCharNo      int
,	@nFontColorID tinyint
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nBalloonColorID AS tinyint

	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tChatColor WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tChatColor 
		   SET nFontColorID = @nFontColorID 
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
		SET @nBalloonColorID = 0
			
		INSERT INTO tChatColor (nCharNo, nFontColorID, nBalloonColorID)
		VALUES (@nCharNo, @nFontColorID, @nBalloonColorID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
		END
	END
END

0倉Q
GO
