CREATE PROCEDURE [dbo].[usp_GuildAcademy_delMember]
	@nNo     int
,	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	DELETE FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1592
	END

	SET NOCOUNT OFF
END

0ꈁq
GO
