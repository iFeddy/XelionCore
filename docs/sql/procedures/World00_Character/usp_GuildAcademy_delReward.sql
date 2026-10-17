CREATE PROCEDURE [dbo].[usp_GuildAcademy_delReward]
	@nNo     int
,	@nLevel  smallint
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
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel AND nTypeID = @nTypeID
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END
	ELSE BEGIN
		DELETE FROM tGuildAcademyReward WHERE nNo = @nNo AND nLevel = @nLevel AND nTypeID > 1
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1596
			RETURN
		END
	END

	SET NOCOUNT OFF
END

0ᐁs
GO
