CREATE PROCEDURE [dbo].[usp_GuildAcademy_setNotify]
	@nNo     int
,	@sNotify nvarchar(1024)
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	UPDATE tGuildAcademy 
	   SET dNotifyDate = getdate(), sNotify = @sNotify 
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1593
	END
	ELSE BEGIN
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0企
GO
