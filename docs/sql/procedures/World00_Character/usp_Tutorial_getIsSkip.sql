CREATE PROCEDURE [dbo].[usp_Tutorial_getIsSkip]
	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT nState FROM tTutorial WHERE nUserNo = @nUserNo AND nState = 1
END

0吁0
GO
