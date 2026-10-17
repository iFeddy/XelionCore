CREATE PROCEDURE [dbo].[usp_Item_getAutoMonList]
	@nCharNo int
,	@nType	 int
,	@nRet    int = 0 output

AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
	
	SET @nRet = 0
	
	SELECT nSlot, nItemID
	  FROM tAutoMon
	 WHERE nCharNo = @nCharNo
	   AND nType = @nType
	   
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END

0訁ឌ
GO
