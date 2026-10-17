CREATE PROCEDURE [dbo].[usp_GuildAcademy_withdrawMoney]
	@nNo     int
,	@nCharNo int
,	@nMoney  bigint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @money AS bigint
	SET @nRet = 0
	SET @money = 0

	IF @nNo <= 0
		BEGIN
			SET @nRet = -1511
			RETURN
		END

	IF @nCharNo <= 0
		BEGIN
			SET @nRet = -1512
			RETURN
		END

	IF @nMoney <= 0
		BEGIN
			SET @nRet = -1514
			RETURN
		END

	BEGIN TRAN
	
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo AND bDeleted = 0)
		BEGIN
			SET @nRet = -1524
			ROLLBACK TRAN
			RETURN
		END

	IF NOT EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo)
		BEGIN
			SET @nRet = -1525
			ROLLBACK TRAN
			RETURN
		END

	SELECT @money = ISNULL(nMoney, 0) FROM tGuildAcademy WITH(UPDLOCK, HOLDLOCK) WHERE nNo = @nNo
	IF @money < @nMoney 
		BEGIN		
			SET @nRet = -1527
			ROLLBACK TRAN
			RETURN
		END
	ELSE
		BEGIN
			UPDATE tGuildAcademy SET nMoney = nMoney - @nMoney WHERE nNo = @nNo AND nMoney >= @nMoney
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @nRet = -1593
					ROLLBACK TRAN
					RETURN
				END

			UPDATE tCharacter SET nMoney = nMoney + @nMoney WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @nRet = -1593
					ROLLBACK TRAN
					RETURN
				END
			ELSE
				BEGIN
					SET @nRet = 0
					COMMIT TRAN
				END
		END
			
	SET NOCOUNT OFF
END
GO
