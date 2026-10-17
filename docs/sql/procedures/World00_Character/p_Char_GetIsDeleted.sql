CREATE PROCEDURE [dbo].[p_Char_GetIsDeleted]
@nCharNo int,
@nRet int OUTPUT

AS
SET NOCOUNT ON

IF EXISTS(SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo AND bDeleted > 0)

BEGIN
	SET @nRet = 1
		RETURN
END

0焁彽ō
GO
