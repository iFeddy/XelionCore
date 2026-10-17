CREATE PROCEDURE [dbo].[usp_User_getAuthList]
AS
BEGIN
	SET NOCOUNT ON

	SELECT nAuthID, sAuthName, bIsLoginAble  FROM tUserAuth

	SET NOCOUNT OFF
END

0
GO
