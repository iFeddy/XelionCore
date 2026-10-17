CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMember]
	@nNo       int
,	@nCharNo   int
,	@dJoinDate datetime
,	@nGrade    tinyint = 0
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @nLevel AS int

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -1522
		RETURN
	END

	IF EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0) BEGIN
		SELECT @nLevel = nLevel FROM tCharacter WHERE nCharNo = @nCharNo
		IF @nLevel <= 0 BEGIN
			SET @nRet = -1523
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -1524
		RETURN
	END
	
	INSERT tGuildAcademyMember (nNo, nCharNo, nJoinLevel, dJoinDate, nGrade)
	VALUES ( @nNo, @nCharNo, @nLevel, @dJoinDate, @nGrade )
	IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END
GO
