CREATE PROCEDURE [dbo].[usp_GuildTournament_getRegistList]
	@nGTNo int
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	SELECT nNo FROM tGuildTournamentRegist WITH(NOLOCK) WHERE nGTNo = @nGTNo

	SET NOCOUNT OFF
END

0ꔁ놵
GO
