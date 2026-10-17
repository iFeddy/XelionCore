CREATE PROCEDURE [dbo].[usp_GuildTournament_setRegist]
	@nGTNo        int
,	@nNo          int
,	@nCharNo      int
,	@nMoney       bigint = 0
,	@nLimitLevel  int = 1
,	@nLimitMember int = 1
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @prevGTNo int, @prevWinnerNo int
	DECLARE @cnt AS int

	SET @nRet = 0
	SET @cnt  = 0

	IF ISNULL(@nGTNo, 0) <= 0 BEGIN
		SET @nRet = -1411
		RETURN
	END

	IF ISNULL(@nNo, 0) <= 0 BEGIN
		SET @nRet = -1412
		RETURN
	END

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -1413
		RETURN
	END

	IF EXISTS(SELECT nGTNo FROM tGuildTournamentRegist WHERE nGTNo = @nGTNo AND nNo = @nNo) BEGIN
		SET @nRet = -1421
		RETURN
	END

	-- �������ʸ�Ʈ ��û���� üũ(����, �ο�)
	SELECT @cnt = ISNULL(COUNT(T1.nCharNo), 0)
	  FROM tGuildMember AS T1, tCharacter AS T2 
	 WHERE T1.nCharNo = T2.nCharNo
	   AND T1.nNo = @nNo
	   AND T2.nLevel >= @nLimitLevel
	IF @cnt < @nLimitMember BEGIN
		SET @nRet = -1422
		RETURN
	END

	-- 2011.8.22 by kodong
	-- �������ʸ�Ʈ ��û�� ���� �������ʸ�Ʈ ���±����� ��û�� �� ����.	
	SELECT TOP 1 @prevGTNo = nGTNo 
	  FROM tGuildTournament 
	 WHERE nGTNo < @nGTNo 
	 ORDER BY nGTNo DESC
	IF @prevGTNo IS NOT NULL BEGIN
		SELECT @prevWinnerNo = nGuildNo
		  FROM tGuildTournament 
		 WHERE nGTNo = @prevGTNo
		   AND nType >= 3028 AND nType <= 3029
		   AND nData = 1
		IF @prevWinnerNo = @nNo BEGIN
			SET @nRet = -1424
			RETURN
		END
	END
	
	BEGIN TRAN

	-- �������ʸ�Ʈ ��û�� �
	EXEC p_Char_MoneySub @nCharNo, @nMoney, @nRet OUTPUT
	IF @nRet <> 1 BEGIN
		SET @nRet = -1423
		ROLLBACK TRAN
		RETURN
	END

	INSERT INTO tGuildTournamentRegist (nGTNo, nNo, nMoney, nSeed)
	VALUES (@nGTNo, @nNo, @nMoney, CAST(RAND() * 100000 AS INT))
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1492
		ROLLBACK TRAN
	END
	ELSE BEGIN
		SET @nRet = 0
		COMMIT TRAN
	END

	SET NOCOUNT OFF
END
GO
