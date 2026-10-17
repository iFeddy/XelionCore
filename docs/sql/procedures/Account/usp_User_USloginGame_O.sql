CREATE PROCEDURE [dbo].[usp_User_USloginGame_O]
	@userID      nvarchar(60)
,	@userPW      nvarchar(32)
,	@userNo      int = 0 OUTPUT
,	@authID      tinyint = 0 OUTPUT
,	@block     	 int = 0 OUTPUT
,	@isLoginable int = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	DECLARE @RET INT	
	SET @userNo = 0
	
	SELECT @userNo = nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
	  FROM dbo.tUser AS U(NOLOCK) INNER JOIN dbo.tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0
	IF @userNo = 0
		BEGIN
			SET @userNo = -1
			SET @authID = 0
			SET @block  = 0
		END
	ELSE
		BEGIN
			IF @block = 1
			BEGIN
				EXEC usp_User_chkBlock @userNo, @RET OUTPUT
				IF @RET = 1
					SELECT @block = bIsBlock FROM tUser(NOLOCK) WHERE nUserNo = @userNo
			END
		END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0
GO
