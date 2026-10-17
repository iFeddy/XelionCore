CREATE PROCEDURE [dbo].[p_Friend_Del]
/*
Friend Delete
2004.12 By CJC
 input:
	nCharNo, nFriendCharNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nFriendCharNo int,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- cross delete
DELETE FROM tFriend WHERE nCharNo = @nCharNo AND nFriendCharNo = @nFriendCharNo
SET @nRet = @@ERROR

DELETE FROM tFriend WHERE nCharNo = @nFriendCharNo AND nFriendCharNo = @nCharNo 
SET @nRet = @nRet + @@ERROR

-- end

0蜁཈Ė
GO
