CREATE PROCEDURE [dbo].[usp_GuildAcademy_setChatBanTime]
	@nNo          int
,	@nCharNo      int
,	@nChatBanTime int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tGuildAcademyMember 
		   SET nChatBanTime = @nChatBanTime
		 WHERE nNo = @nNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			RETURN
		END
	END 
	ELSE BEGIN
		SET @nRet = -1511
		RETURN
	END

	SET NOCOUNT OFF
END

0
GO
