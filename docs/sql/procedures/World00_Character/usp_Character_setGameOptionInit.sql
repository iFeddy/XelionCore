CREATE PROCEDURE [dbo].[usp_Character_setGameOptionInit]
	@nCharNo   int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	DELETE FROM tCharacterGameOption WHERE nCharNo = @nCharNo
END

0윁惈;
GO
