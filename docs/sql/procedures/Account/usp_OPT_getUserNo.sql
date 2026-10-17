CREATE PROCEDURE [dbo].[usp_OPT_getUserNo]
	@userID nvarchar(60)
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo FROM tUser WHERE sUserID = @userID

	SET NOCOUNT OFF
END

0
GO
