CREATE PROCEDURE [dbo].[usp_GuildAcademy_getRanking]
	@nCount int = 1
,	@nRet   int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @MAX_COUNT AS INT

	SELECT @MAX_COUNT = ISNULL(MAX(nCount), 0)
	  FROM tGuildAcademyRanking
	
	SELECT dDate, nRanking, nGuildNo, sGuildName, nMasterNo, sMasterID, nPoint
	  FROM tGuildAcademyRanking
	 WHERE nCount = @MAX_COUNT - (@nCount - 1)

	SET NOCOUNT OFF
END

0ᔁ潲|
GO
