CREATE PROCEDURE [dbo].[usp_GuildAcademy_setRanking]
	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @MAX_COUNT AS INT

	BEGIN TRAN

	SELECT @MAX_COUNT = ISNULL(MAX(nCount), 0) + 1
	  FROM tGuildAcademyRanking

	INSERT INTO tGuildAcademyRanking (dDate, nRanking, nGuildNo, sGuildName, nMasterNo, sMasterID, nPoint, nCount)
		SELECT GETDATE() AS dDate, R.nRanking, R.nNo, G.sName, GA.nMasterCharNo, C.sID, R.nPoint, @MAX_COUNT
		  FROM tGuild AS G, tGuildAcademy AS GA, tCharacter AS C, (
				SELECT COUNT(T2.nPoint) + 1 AS nRanking, T1.nNo, T1.nPoint
				  FROM tGuildAcademy T1 LEFT OUTER JOIN tGuildAcademy T2 ON T1.nPoint < T2.nPoint
				 GROUP BY T1.nNo, T1.nPoint) AS R
		 WHERE G.nNo = R.nNo AND C.nCharNo = GA.nMasterCharNo AND GA.nNo = R.nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1592
		ROLLBACK TRAN
		RETURN
	END

	UPDATE tGuildAcademy SET nPoint = 0
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -1593
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END

0蠁
GO
