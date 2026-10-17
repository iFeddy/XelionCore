CREATE PROCEDURE [dbo].[usp_OPT_getUserID]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT sUserID FROM tUser WHERE nUserNo = @userNo

	SET NOCOUNT OFF
END

0
GO
