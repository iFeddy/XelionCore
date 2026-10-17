CREATE PROCEDURE [dbo].[p_Friend_Get]

/*
Char Get friend
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nFriendNo int

AS

SET NOCOUNT ON SELECT 
F.nFriendCharNo,
F.nFlags, 

C.sID,
S.nClass,
C.nLevel,
C.sLoginZone,
C.dLastLoginDate

FROM tFriend AS F
INNER JOIN tCharacterShape AS S ON S.nCharNo =  @nFriendNo
INNER JOIN tCharacter AS C ON C.nCharNo = @nFriendNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo AND F.nFriendCharNo = @nFriendNo

-- end

0更ė
GO
