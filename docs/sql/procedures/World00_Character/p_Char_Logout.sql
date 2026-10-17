CREATE PROCEDURE [dbo].[p_Char_Logout]
	@nCharNo	int
,	@nPlayMin	int
,	@nRet		int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @nLastFriendPoint int

	SET @nRet = 1
	SET @nLastFriendPoint = (SELECT nFriendPoint FROM tCharacter WHERE nCharNo = @nCharNo)
	

	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tCharacter SET nPlayMin = nPlayMin + @nPlayMin WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 0
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			UPDATE tCharacter SET nLastFriendPoint = @nLastFriendPoint WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = 0
				IF @@TRANCOUNT > 0 ROLLBACK TRAN
				RETURN
			END
			ELSE BEGIN
				IF @@TRANCOUNT > 0 COMMIT TRAN
				RETURN
			END
		END
	END
	ELSE BEGIN
		SET @nRet = 0
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
