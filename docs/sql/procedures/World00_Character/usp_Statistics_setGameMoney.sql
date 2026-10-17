CREATE PROCEDURE [dbo].[usp_Statistics_setGameMoney]
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tStatGameMoney
	SELECT T1.Date, T1.Lv030, T1.Lv060, T1.Lv110, T1.Lv120, T2.UserMoney, T3.GuildMoney
	  FROM (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date],
				SUM(CASE WHEN nLevel >= 1 AND nLevel <= 30 THEN nMoney ELSE 0 END) AS [Lv030], 
				SUM(CASE WHEN nLevel >= 31 AND nLevel <= 60 THEN nMoney ELSE 0 END) AS [Lv060], 
				SUM(CASE WHEN nLevel >= 61 AND nLevel <= 110 THEN nMoney ELSE 0 END) AS [Lv110], 
				SUM(CASE WHEN nLevel >= 111 AND nLevel <= 999 THEN nMoney ELSE 0 END) AS [Lv120]
			  FROM tCharacter WITH(NOLOCK)
			 WHERE nAdminLevel = 0 AND bDeleted = 0) AS T1, (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date], SUM(nUserMoney) AS [UserMoney]
			  FROM tUserMoney WITH(NOLOCK)) AS T2, (
			SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS [Date], SUM(nMoney) AS [GuildMoney]
			  FROM tGuild WITH(NOLOCK)) AS T3
	 WHERE T1.Date = T2.Date
	   AND T1.Date = T3.Date

	SET NOCOUNT OFF
END

0
GO
