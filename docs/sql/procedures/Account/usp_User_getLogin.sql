CREATE PROCEDURE [dbo].[usp_User_getLogin]
	@userID   NVARCHAR(60)
,	@userPW   NVARCHAR(32)
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.nUserNo, sUserID, sUserPW, sUserName, bIsBlock, nAuthID, sEmail, sJuminNo, sPhoneNo
	  FROM tUser AS T1 LEFT JOIN tUserOption AS T2 ON T1.nUserNo = T2.nUserNo
	 WHERE sUserID = @userID AND sUserPW = @userPW AND bIsDelete = 0

	SET NOCOUNT OFF
END

0
GO
