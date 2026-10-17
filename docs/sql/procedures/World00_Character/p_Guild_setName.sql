CREATE PROCEDURE [dbo].[p_Guild_setName]
	@nNo   int
,	@sName nvarchar(18)
,	@nRet  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	IF @nNo <= 0 BEGIN
		SET @nRet = -1311
		RETURN
	END

	IF LEN(ISNULL(@sName, '')) = 0 BEGIN
		SET @nRet = -1312
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuild WHERE nNo = @nNo) BEGIN
		IF EXISTS(SELECT sName FROM tGuild WHERE sName = @sName) BEGIN
			SET @nRet = -1322
		END
		ELSE BEGIN

			BEGIN TRAN

			UPDATE tGuild SET sName = @sName WHERE nNo = @nNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				ROLLBACK TRAN
				SET @nRet = -1393
			END
			ELSE BEGIN
				COMMIT TRAN
				SET @nRet = 0
			END
		END
	END
	ELSE BEGIN
		SET @nRet = -1321
	END

	SET NOCOUNT OFF
END

0逷7
GO
