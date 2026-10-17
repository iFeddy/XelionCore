CREATE PROCEDURE [dbo].[usp_GuildAcademy_setGraduate]
	@nNo            int
,	@nCharNo        int
,	@sID            nvarchar(40)
,	@nClass         tinyint
,	@nGraduateLevel smallint
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @nJoinLevel smallint
	DECLARE @dJoinDate  datetime

	SET @nRet = 0
	SET @nJoinLevel = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF LEN(ISNULL(@sID, '')) = 0 BEGIN
		SET @nRet = -1515
		RETURN
	END

	IF @nClass <= 0 BEGIN
		SET @nRet = -1516
		RETURN
	END

	IF @nGraduateLevel <= 0 BEGIN
		SET @nRet = -1517
		RETURN
	END

	SELECT @nJoinLevel = nJoinLevel, @dJoinDate = dJoinDate
	  FROM tGuildAcademyMember
	 WHERE nNo = @nNo AND nCharNo = @nCharNo
	IF @nJoinLevel <= 0 BEGIN
		SET @nRet = -1528
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tGuildAcademyGraduate (nNo, nCharNo, sID, nClass, nJoinLevel, dJoinDate, nGraduateLevel)
		VALUES (@nNo, @nCharNo, @sID, @nClass, @nJoinLevel, @dJoinDate, @nGraduateLevel)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1592
			RETURN
		END
	END

	SET NOCOUNT OFF
END
GO
