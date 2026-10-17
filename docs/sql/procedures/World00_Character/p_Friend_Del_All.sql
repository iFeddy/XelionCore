CREATE PROCEDURE [dbo].[p_Friend_Del_All]
/*
Friend Delete
2007.12 By CJC
 input:
	nCharNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- cross delete
DELETE FROM tFriend WHERE nCharNo = @nCharNo OR nFriendCharNo = @nCharNo
SET @nRet = @@ERROR

-- end

0
GO
