CREATE PROCEDURE [dbo].[usp_GuildAcademy_getMoney]
	@nNo  int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		SELECT nMoney FROM tGuildAcademy WHERE nNo = @nNo
	END
	ELSE BEGIN
		SET @nRet = -1521
	END

	SET NOCOUNT OFF
END

0ꌁ蜩z
GO
