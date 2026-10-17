CREATE PROCEDURE [dbo].[usp_InGameGuide_getViewList]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tInGameGuide WITH(NOLOCK) WHERE nCharNo = @nCharNo) BEGIN
		SELECT nGuideNo
		  FROM tInGameGuide 
		 WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		SET @nRet = -1111
	END
END

0㐁庳
GO
