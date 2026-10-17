CREATE PROCEDURE [dbo].[p_Char_DelOptions]
	@nCharNo		int,
	@nRet			int OUTPUT
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	DELETE FROM tCharacterOptions WHERE nCharNo = @nCharNo 
	SET @nRet = @@ERROR
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0⋔ł
GO
