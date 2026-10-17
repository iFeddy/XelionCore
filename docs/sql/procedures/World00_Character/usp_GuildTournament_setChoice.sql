CREATE PROCEDURE [dbo].[usp_GuildTournament_setChoice]
	@nGTNo        int
,	@nLimitLevel  int = 1
,	@nLimitMember int = 1
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	BEGIN TRAN

	UPDATE tGuildTournamentRegist
	   SET bChoice = 1
	 WHERE nGTNo = @nGTNo AND nNo IN (
			SELECT TOP 16 T4.nNo
			  FROM tGuildTournamentRegist AS T4 INNER JOIN (
					SELECT T1.nNo, ISNULL(COUNT(T2.nCharNo), 0) AS nMember
					  FROM tGuildMember AS T1, tCharacter AS T2 
					 WHERE T1.nCharNo = T2.nCharNo
					   AND T2.nLevel >= @nLimitLevel
					 GROUP BY T1.nNo
					HAVING COUNT(T2.nCharNo) >= @nLimitMember) AS T3 ON T4.nNo = T3.nNo
			 WHERE nGTNo = @nGTNo
			 ORDER BY nSeed DESC)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1493
		ROLLBACK TRAN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END

0ꗙ
GO
