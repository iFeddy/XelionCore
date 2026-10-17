CREATE PROCEDURE [dbo].[p_UserGetAuthList]
AS
BEGIN
	SET NOCOUNT ON

	SELECT nAuthID, sAuthName, bIsLoginAble  FROM tUserAuth

	SET NOCOUNT OFF
END


0
GO
