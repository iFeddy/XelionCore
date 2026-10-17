CREATE PROCEDURE [dbo].[p_Char_CheckOptionsExistance]
	@nCharNo		int,
	@nRet			int OUTPUT
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	IF EXISTS(SELECT * FROM  tCharacterOptions WHERE nCharNo = @nCharNo )
		SET @nRet = 1
	ELSE
		SET @nRet = 0
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0漁㓱ļ
GO
