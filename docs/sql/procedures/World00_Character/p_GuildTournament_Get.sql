CREATE PROCEDURE [dbo].[p_GuildTournament_Get]
/*
Get Tournament Guild Data
2007.5 By CJC
*/
@nGTNo int
AS
SET NOCOUNT ON
SELECT nType, nData FROM tGuildTournament WHERE @nGTNo = nGTNo
-- end

0
GO
