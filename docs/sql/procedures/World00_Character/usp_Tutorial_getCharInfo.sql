CREATE PROCEDURE [dbo].[usp_Tutorial_getCharInfo]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState, nStep FROM tTutorial WHERE nCharNo = @nCharNo 
END

0ᬁ/
GO
