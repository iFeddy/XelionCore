CREATE PROCEDURE [dbo].[usp_User_updatePassword]
	@userNo int
,	@prevPW nvarchar(32)
,	@userPW nvarchar(32)
,	@userIP nvarchar(30)
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @CURR_PW AS NVARCHAR(20)

	IF EXISTS(SELECT nUserNo FROM tUser WITH(NOLOCK) WHERE nUserNo = @userNo AND bIsDelete = 0)
		BEGIN
			SELECT @CURR_PW = sUserPW FROM tUser WHERE nUserNo = @userNo
			If @CURR_PW <> @prevPW
				RETURN (-2)

			BEGIN TRAN

			UPDATE tUser SET sUserPW = @userPW WHERE nUserNo = @userNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			INSERT INTO tPasswordLog (nUserNo, sPrevPW, sCurrPW, sUserIP)
			VALUES (@userNo, @CURR_PW, @userPW, @userIP)
			IF @@ERROR <> 0 OR @@IDENTITY = 0
				BEGIN
					ROLLBACK TRAN
					RETURN (-9)
				END

			COMMIT TRAN
			RETURN (0)
		END
	ELSE
		RETURN (-1)
END

<
GO
