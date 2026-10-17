CREATE PROCEDURE [dbo].[p_Char_SetLevelExp]
@nCharNo	int,
@nLevel		smallint,
@nExp		bigint,

@nRet		int OUTPUT

AS

SET NOCOUNT ON

UPDATE tCharacter SET 
	nLevel			= @nLevel,
	nExp			= @nExp
WHERE nCharNo = @nCharNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	SET @nRet = 0
	RETURN
END
SET @nRet = 1

0朁諁ŷ
GO
