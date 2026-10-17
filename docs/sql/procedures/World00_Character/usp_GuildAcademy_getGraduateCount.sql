CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGraduateCount]
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

	SELECT ISNULL(COUNT(nCharNo), 0) AS cnt
	  FROM tGuildAcademyGraduate
	 WHERE nNo = @nNo
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END

	SET NOCOUNT OFF
END

0뼁뚘v
GO
