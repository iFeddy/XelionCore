CREATE PROCEDURE [dbo].[usp_GuildAcademy_delRewardAll]
	@nNo     int
,	@nTypeID tinyint
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
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID = @nTypeID
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END
	ELSE BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END

	SET NOCOUNT OFF
END

0䴁칐t
GO
