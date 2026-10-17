CREATE PROCEDURE [dbo].[usp_GuildEmblem_SetInfo]
/*
Set Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = -1,  Error(Not found guild)
	nRet = -2,  Error(Not found guild)
	nRet = 1,  OK
*/
@nGuildNo int,
@nColor int,
@nIcon int,
@dCreated datetime,
@dUpdated datetime,
@nUnk int,
@nRet int output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nGuildNo FROM tGuildEmblem WHERE nGuildNo = @nGuildNo) BEGIN
		UPDATE tGuildEmblem
		SET nColor = @nColor, nIcon = @nIcon, dUpdated = @dUpdated, nUnk = @nUnk
		WHERE nGuildNo = @nGuildNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tGuildEmblem (nGuildNo, nColor, nIcon, dCreated, dUpdated, nUnk)
		VALUES (@nGuildNo, @nColor, @nIcon, @dCreated, @dUpdated, @nUnk)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2
			RETURN
		END
	END			
	
	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		RETURN
	END
	RETURN @@ERROR	

	SET NOCOUNT OFF
END
-- end
GO
