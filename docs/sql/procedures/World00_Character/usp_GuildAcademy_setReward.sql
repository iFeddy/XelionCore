CREATE PROCEDURE [dbo].[usp_GuildAcademy_setReward]
	@nNo     int
,	@nLevel  smallint
,	@nTypeID tinyint
,	@nValue  bigint
,	@nItemID int = NULL
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nLevel <= 0 BEGIN
		SET @nRet = -1513
		RETURN
	END

	INSERT INTO tGuildAcademyReward (nNo, nLevel, nTypeID, nValue, nItemID)
	VALUES (@nNo, @nLevel, @nTypeID, @nValue, @nItemID)
	IF @@ERROR <> 0 BEGIN
		SET @nRet = -1595
		RETURN
	END

	SET NOCOUNT OFF
END

0섁
GO
