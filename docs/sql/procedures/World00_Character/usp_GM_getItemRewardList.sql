CREATE PROCEDURE [dbo].[usp_GM_getItemRewardList]
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
		
	SELECT nGroupNo, sys_create_date AS dDate, sCause 
	  FROM tGMItemReward 
	 ORDER BY nGroupNo DESC
	
END

0嬁l
GO
