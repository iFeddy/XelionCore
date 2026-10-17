CREATE PROCEDURE [dbo].[usp_GuildAcademy_getRewardAll]
	@nNo     int
,	@nTypeID int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	IF @nTypeID = 1 BEGIN
		SELECT nLevel, nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID = 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1591
		END	
	END 
	ELSE BEGIN
		SELECT nLevel, nTypeID, nValue, nItemID FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1591
		END	
	END

	SET NOCOUNT OFF
END

0蜁垺~
GO
