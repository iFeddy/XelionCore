CREATE PROCEDURE [dbo].[p_Char_SetAdminLevel]
@nCharNo	int,
@nAdmin		Tinyint,

@nRet		int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET 
	nAdminLevel			= @nAdmin
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0鈁祊Ŵ
GO
