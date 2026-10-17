CREATE PROCEDURE [dbo].[p_GuildTournament_GetLastMatch]
AS
BEGIN
	SELECT   nGTNo, nType, nData
	FROM      tGuildTournament
	WHERE   (nGTNo IN
                    (SELECT   MAX(nGTNo)
                     FROM      tGuildTournament))
             ORDER BY nType
END

0
GO
