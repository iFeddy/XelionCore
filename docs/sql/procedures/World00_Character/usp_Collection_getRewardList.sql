CREATE PROCEDURE [dbo].[usp_Collection_getRewardList]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	SELECT nReward 
	  FROM tCollectionReward 
	 WHERE nCharNo = @nCharNo 
END

0༁쮕Ŕ
GO
