CREATE PROCEDURE [dbo].[usp_GuildEmblem_GetInfoAll]

AS
BEGIN
	SET NOCOUNT ON
	SELECT nGuildNo, nColor, nIcon, dCreated, dUpdated, nUnk FROM tGuildEmblem
	SET NOCOUNT OFF
END0
GO
