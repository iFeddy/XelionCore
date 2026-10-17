CREATE  PROCEDURE [dbo].[usp_Guild_GetToken]
/*
Get Guild Token
2008.12 By Lunar
Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/
@nNo int,
@nToken bigint output,
@nRet int output
AS
SET NOCOUNT ON
SET @nRet = 0
SELECT @nToken = ISNULL(MAX(nToken), 0) FROM tGuildToken WHERE nNo = @nNo
IF @@ERROR <> 0 --OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
-- end

0
GO
