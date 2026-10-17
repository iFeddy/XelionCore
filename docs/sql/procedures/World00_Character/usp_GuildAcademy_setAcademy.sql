CREATE PROCEDURE [dbo].[usp_GuildAcademy_setAcademy]
	@nNo           int
,	@nMasterCharNo int
,	@nMaxMember    smallint
,	@nRet          int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		SET @nRet = -1521
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tGuildAcademy (nNo, nMasterCharNo, nMaxMember)
		VALUES (@nNo, @nMasterCharNo, @nMaxMember)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1592
			RETURN
		END
	END
END

0更䀂
GO
