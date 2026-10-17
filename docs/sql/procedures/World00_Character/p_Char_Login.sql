CREATE PROCEDURE [dbo].[p_Char_Login]
	@nCharNo	int
,	@nRet		int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @dLastloginDate AS datetime
	
	SET @nRet = 0

	IF EXISTS (SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo) BEGIN
		BEGIN TRAN

		SET @dLastloginDate = (SELECT dLastLoginDate FROM tCharacter WHERE nCharNo = @nCharNo)

		UPDATE tCharacter SET nLoginCount = nLoginCount + 1, dLastLoginDate = GETDATE() WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = 1
			IF @@TRANCOUNT > 0 ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			SELECT DATEDIFF(d,@dLastLoginDate,GETDATE()) AS nLoginDateDiff
			  FROM tCharacter
			 WHERE nCharNo = @nCharNo

			IF @@TRANCOUNT > 0 COMMIT TRAN
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
