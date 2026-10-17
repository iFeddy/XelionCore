CREATE PROCEDURE [dbo].[p_Char_Wedding_Get_Partner_Info]
	@CharNo INT
AS

BEGIN

	SET NOCOUNT ON

	SELECT nLevel, dLastLoginDate FROM tCharacter WHERE nCharNo = @CharNo

	SET NOCOUNT OFF

END

0
GO
