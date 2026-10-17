CREATE PROCEDURE [dbo].[usp_User_insert_Old]
	@userID   nvarchar(60)
,	@userPW   nvarchar(32)
,	@userName nvarchar(10)
,	@userIP   nvarchar(30)
,	@eMail    nvarchar(70) = NULL
,	@juminNo  nvarchar(20) = NULL
,	@phoneNo  nvarchar(20) = NULL
,	@zipNo    int = NULL
,	@address  nvarchar(90) = NULL
,	@isMail   bit = NULL
,	@channel  nvarchar(20) = 'ONSON'
,	@isOption bit
,	@userNo   int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET AS INT

	EXEC @RET = usp_User_chkUserID @userID
	IF (@RET <> 0)
		RETURN (-2)

	BEGIN TRAN

	INSERT INTO tUser (sUserID, sUserPW, sUserName, sUserIP)
	VALUES (@userID, @userPW, @userName, @userIP)
	IF (@@ERROR <> 0 OR @@IDENTITY = 0)
		BEGIN
			ROLLBACK TRAN
			RETURN (-9)
		END
	ELSE
		BEGIN			
			SET @userNo = @@IDENTITY

			IF (@isOption = 1)
				BEGIN
					EXEC @RET = usp_User_insertOption @userNo, @eMail, @juminNo, @phoneNo, @zipNo, @address, @isMail, @channel
					IF (@RET <> 0)
						BEGIN
							ROLLBACK TRAN
							RETURN (-3)
						END
				END

			COMMIT TRAN
			RETURN (0)
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT ON
END
GO
