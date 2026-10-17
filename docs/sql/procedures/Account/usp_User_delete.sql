CREATE PROCEDURE [dbo].[usp_User_delete]
	@userNo  int
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN

			UPDATE tUser
			   SET sUserID = '~' + sUserID + '_' + CONVERT(NVARCHAR(8), GETDATE(), 112),
			       bIsDelete = 1
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			UPDATE tUserOption 
			   SET sJuminNo = '~' + sJuminNo + '_' + CONVERT(NVARCHAR(8), GETDATE(), 112)
			 WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END
			ELSE
				BEGIN
					COMMIT TRAN
					RETURN (0)
				END	
		END
	ELSE
		RETURN (-1)

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
