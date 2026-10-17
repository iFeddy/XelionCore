CREATE PROCEDURE [dbo].[usp_GuildAcademy_raiseMember]
	@nNo     int
,	@nCharNo int
,	@nGrade  tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nMax AS INT
	DECLARE @nCnt AS INT

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		SELECT @nMax = nMaxMembers, @nCnt = nCnt
		  FROM tGuild AS T1 INNER JOIN (
				SELECT nNo, ISNULL(COUNT(*), 0) AS nCnt
				  FROM tCharacter AS C, tGuildMember AS G
				 WHERE C.nCharNo = G.nCharNo 
				   AND G.nNo = @nNo
				   AND C.bDeleted = 0
				 GROUP BY nNo) AS T2 ON (T1.nNo = T2.nNo)
		IF @nCnt < @nMax BEGIN
			BEGIN TRAN
	
			-- 길드아카데미 탈퇴(0:ok)
			EXEC usp_GuildAcademy_delMember @nNo, @nCharNo, @nRet output
			IF @nRet <> 0 BEGIN
				ROLLBACK TRAN
				RETURN
			END
	
			-- 길드 가입(1:ok)
			EXEC p_Guild_Member_Join @nNo, @nCharNo, @nGrade, @nRet output
			IF @nRet <> 1 BEGIN
				SET @nRet = -1531
				ROLLBACK TRAN
				RETURN
			END
			ELSE BEGIN
				SET @nRet = 0
				COMMIT TRAN
			END
		END
		ELSE BEGIN
			SET @nRet = -1530
		END
	END
	ELSE BEGIN
		SET @nRet = -1528
	END
END

0쀁䯞
GO
