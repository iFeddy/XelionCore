CREATE PROCEDURE [dbo].[usp_Char_GetAutoPick]
/*
2009.6 By temlost
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	-- 삭제되지 않은 실존하는 케릭터인지 검색
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo and bDeleted = 0)
	BEGIN
		SET @nRet = -4444
		RETURN
	END

	IF EXISTS(SELECT nCharNo, IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
--		SELECT nCharNo, getdate() FROM tAutoPick WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
END

0ꠁ熇%
GO
