CREATE PROCEDURE [dbo].[usp_User_setLoginable]
	@authID      int
,	@isLoginable int
,	@RET         int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF EXISTS(SELECT nAuthID FROM tUserAuth WHERE nAuthID = @authID)
		BEGIN
			UPDATE tUserAuth SET bIsLoginAble = @isLoginable WHERE nAuthID = @authID
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				SET @RET = 0
			ELSE
				SET @RET = 1
		END
	ELSE
		SET @RET = 0

	SET NOCOUNT OFF
END

0
GO
