CREATE PROCEDURE [dbo].[usp_Character_SetFriendpoint]
	@nCharNo		int
,	@nFriendPoint	int
,	@nRet			int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SET @nRet = 0
	
	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		UPDATE tCharacter SET nFriendPoint = @nFriendPoint WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			IF @@TRANCOUNT > 0 COMMIT TRAN
			SET @nRet = 0
			RETURN
		END
	END
	ELSE BEGIN
		SET @nRet = -1111
	END
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
