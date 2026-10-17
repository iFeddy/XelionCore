CREATE PROCEDURE [dbo].[usp_User_USloginGame]
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

	/*-------------------------------------------------------------------------------------
		Project: FiestaHeroes - NA2016
		Author: xkl
		
		Using the "r_" parameter, automatically create the account if it doesn't exist.
	------------------------------------------------------------------------------------*/
	IF EXISTS (SELECT Enabled FROM Options..Auto_Account_Creation WHERE Enabled = 1)
	BEGIN
		IF NOT EXISTS (SELECT sUserID FROM tUser WHERE sUserID = @userID)
		BEGIN
			-- If the username starts with "r_", trim it.
			IF @userID LIKE 'r[_]%'
			BEGIN
				SET @userID = SUBSTRING(@userID, 3, LEN(@userID) - 2);
	
				-- If the trimmed username doesn't already exist, insert it into the table.
				IF NOT EXISTS (SELECT sUserID FROM tUser WHERE sUserID = @userID)
				BEGIN
					INSERT INTO tUser(sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP)
					VALUES (@userID, @userPW, @userID, 0, 0, 1, '127.0.0.1')
				END
			END
			ELSE
			BEGIN
				-- The username doesn't start with "r_", so exit.
				RETURN;
			END
		END
	END
	/*-------------------------------------------------------------------------------------
		Continue as normal ~ 
	------------------------------------------------------------------------------------*/

	
	SELECT @userNo = nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
	  FROM tUser AS U(NOLOCK) INNER JOIN tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0
	IF @userNo = 0 BEGIN
		SELECT @userNo = U.nUserNo, @authID = U.nAuthID, @block = bIsBlock, @isLoginable = bIsLoginAble
		  FROM tUser AS U(NOLOCK) 
		 INNER JOIN tUserAuth AS UA(NOLOCK) ON U.nAuthID = UA.nAuthID
		 INNER JOIN tUserOption AS UO(NOLOCK) ON U.nUserNo = UO.nUserNo 
		 WHERE UO.sEmail = @userID AND sUserPW = @userPW AND bIsDelete = 0
		IF @userNo = 0 BEGIN
			SET @userNo = -1
			SET @authID = 0
			SET @block  = 0
			RETURN
		END
	END	
	
	IF @block = 1 BEGIN
		EXEC usp_User_chkBlock @userNo, @RET OUTPUT
		IF @RET = 1 BEGIN
			SELECT @block = bIsBlock FROM tUser(NOLOCK) WHERE nUserNo = @userNo
		END
	END

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
GO
