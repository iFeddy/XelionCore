CREATE PROCEDURE [dbo].[usp_Char_SetAutoPick]
/*
get cen reward
2009.4 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@nIsAutoPick		tinyint,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	-- 삭제되지 않은 실존하는 케릭터인지 검색
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo and bDeleted = 0)
	BEGIN
		SET @nRet = -5654;
		RETURN
	END

	IF EXISTS(SELECT nCharNo, IsAutoPick FROM tAutoPick WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tAutoPick
		SET IsAutoPick = @nIsAutoPick
		WHERE nCharNo = @nCharNo

		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
		INSERT INTO tAutoPick (nCharNo, IsAutoPick)
		VALUES (@nCharNo, @nIsAutoPick)

		SET @nRet = 0
		RETURN
	END
END

0숁ⷸĨ
GO
