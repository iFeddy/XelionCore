CREATE PROCEDURE [dbo].[usp_User_chkUserID]
	@userID nvarchar(60) = ''
AS
BEGIN
	SET NOCOUNT ON

	IF @userID = '' OR LEN(@userID) < 5 OR LEN(@userID) > 16
		RETURN -1

	IF EXISTS (SELECT sUserID FROM tUser WITH(NOLOCK) WHERE sUserID = @userID AND bIsDelete = 0)
		RETURN -2
	ELSE
		RETURN 0

	SET NOCOUNT OFF
END
GO
