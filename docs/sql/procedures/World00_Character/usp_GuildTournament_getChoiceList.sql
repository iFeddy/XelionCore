CREATE PROCEDURE [dbo].[usp_GuildTournament_getChoiceList]
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

	SELECT nNo FROM tGuildTournamentRegist WITH(NOLOCK) WHERE nGTNo = @nGTNo AND bChoice = 1

	SET NOCOUNT OFF
END

0쀁➥ċ
GO
