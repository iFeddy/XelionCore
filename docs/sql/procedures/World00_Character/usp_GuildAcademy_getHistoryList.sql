CREATE PROCEDURE [dbo].[usp_GuildAcademy_getHistoryList]

@nAcademyID int,
@nAmount int,
@nRet int output

AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT TOP (@nAmount) sID, nAction, dDate FROM tGuildAcademyHistory WHERE nAcademyID = @nAcademyID

	SET NOCOUNT OFF
END
GO
