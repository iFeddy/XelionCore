CREATE PROCEDURE [dbo].[usp_Character_setChatBan]
	@sID       nvarchar(40)
,	@nChatBan  int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000 
	
	DECLARE @nCharNo AS int
	
	SET @nRet = 0
	
	SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID AND bDeleted = 0
	IF @nCharNo IS NULL BEGIN
		SET @nRet = -3121
	END
	ELSE BEGIN
		IF NOT EXISTS(SELECT nCharNo FROM tCharacterSingleOption WHERE nCharNo = @nCharNo) BEGIN
			INSERT INTO tCharacterSingleOption (nCharNo, nChatBan)
			VALUES (@nCharNo, @nChatBan)
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3192
				RETURN
			END		
		END
		ELSE BEGIN
			UPDATE tCharacterSingleOption 
			   SET nChatBan = @nChatBan 
			 WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 BEGIN
				SET @nRet = -3193
				RETURN
			END		
		END
	END
END

0⬁䍶9
GO
