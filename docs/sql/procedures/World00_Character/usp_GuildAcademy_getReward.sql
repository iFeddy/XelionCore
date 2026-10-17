CREATE PROCEDURE [dbo].[usp_GuildAcademy_getReward]
	@nNo    int
,	@nLevel tinyint
,	@nRet   int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nLevel <= 0 BEGIN
		SET @nRet = -1513
		RETURN
	END

	SELECT nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1591
	END	

	SET NOCOUNT OFF
END

0丁掖}
GO
