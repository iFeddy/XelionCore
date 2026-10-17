CREATE PROCEDURE [dbo].[p_GuildTournament_GetAllTournamentList]
AS
BEGIN
	SELECT nGTNo
	FROM   tGuildTournament
	GROUP BY nGTNo
END

0
GO
