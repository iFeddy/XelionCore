CREATE PROCEDURE [dbo].[usp_Guild_GetTokenAll]
AS
BEGIN
	SET NOCOUNT ON
	SELECT nNo, nToken FROM tGuildToken
	SET NOCOUNT OFF
END

0、ছp
GO
