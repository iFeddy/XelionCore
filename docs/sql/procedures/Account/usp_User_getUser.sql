CREATE PROCEDURE [dbo].[usp_User_getUser]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo, sUserID, sUserPW, sUserName, bIsBlock, bIsDelete, nAuthID, sUserIP, dDate
	  FROM tUser WITH(NOLOCK)
	 WHERE nUserNo = @userNo

	SET NOCOUNT OFF
END
GO
