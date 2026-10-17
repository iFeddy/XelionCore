CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGraduateList]
	@nNo int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	SELECT sID, nClass, nJoinLevel, dJoinDate, nGraduateLevel, dGraduateDate
	  FROM tGuildAcademyGraduate
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END

0ꪼw
GO
