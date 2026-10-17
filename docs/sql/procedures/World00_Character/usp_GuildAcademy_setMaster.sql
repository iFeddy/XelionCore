CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMaster]
	@nNo       int
,	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF EXISTS(SELECT nNo FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tGuildAcademy SET nMasterCharNo = @nCharNo
		 WHERE nNo = @nNo
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

0ꐁᱯ
GO
