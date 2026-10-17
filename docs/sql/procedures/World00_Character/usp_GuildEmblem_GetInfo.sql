CREATE PROCEDURE [dbo].[usp_GuildEmblem_GetInfo]

@nGuildNo INT

AS
BEGIN
	SET NOCOUNT ON
	SELECT nColor, nIcon, dCreated, dUpdated, nUnk FROM tGuildEmblem
	WHERE nGuildNo = @nGuildNo
	SET NOCOUNT OFF
END
GO
