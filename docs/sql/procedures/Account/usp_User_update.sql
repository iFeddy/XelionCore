CREATE PROCEDURE [dbo].[usp_User_update]
	@userNo   int
,	@eMail    nvarchar(200) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			BEGIN TRAN

			UPDATE tUserOption 
			   SET sEmail = @eMail, sPhoneNo = @phoneNo, nZipNo = @zipNo, sAddress = @address, bIsMail = @isMail
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
