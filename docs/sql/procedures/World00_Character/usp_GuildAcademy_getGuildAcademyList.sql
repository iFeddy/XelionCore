CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGuildAcademyList]
	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT GA.nNo, nMasterCharNo, C1.sID, nMaxMember, GA.nMoney, nPoint
		 , T2.nRanking, ISNULL(dBuffEndTime, GETDATE()), ISNULL(dNotifyDate, GETDATE()), ISNULL(sNotify, ''), ISNULL(sIntro, '')
		 , T1.nCharNo, T1.sID, ISNULL(T3.nGraduateCount, 0) AS nGraduateCount
	  FROM tGuildAcademy AS GA INNER JOIN 
		   tCharacter AS C1 ON GA.nMasterCharNo = C1.nCharNo INNER JOIN (
				SELECT G.nNo, GM.nCharNo, C.sID
				  FROM tGuild AS G, tGuildMember AS GM, tCharacter AS C
				 WHERE G.nNo = GM.nNo AND GM.nCharNo = C.nCharNo
				   AND GM.nGrade = 0 AND C.bDeleted = 0
				) AS T1 ON GA.nNo = T1.nNo INNER JOIN (
				-- ranking check
				SELECT T1.nNo, COUNT(T2.nPoint) + 1 AS nRanking
				  FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint < T2.nPoint
				 GROUP BY T1.nNo, T1.nPoint
				) AS T2 ON GA.nNo = T2.nNo LEFT OUTER JOIN (
				-- graduate counting
				SELECT nNo, ISNULL(COUNT(nCharNo), 0) AS nGraduateCount
				  FROM tGuildAcademyGraduate
				 GROUP BY nNo
				) AS T3 ON GA.nNo = T3.nNo
	 ORDER BY T2.nRanking
	IF @@ERROR <> 0
		BEGIN
			SET @nRet = -1594
		END

	SET NOCOUNT OFF
END

0㄁黡x
GO
