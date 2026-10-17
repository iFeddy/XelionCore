CREATE PROCEDURE [dbo].[p_Friend_GetAll2]

/*
Char Get friend all
2007.2.14 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT 
F.nFriendCharNo,
F.nFlags, 

C.sID,
S.nClass,
C.nLevel,
C.sLoginZone,
C.dLastLoginDate

FROM tFriend AS F
INNER JOIN tCharacterShape AS S ON S.nCharNo = F.nFriendCharNo
INNER JOIN tCharacter AS C ON C.nCharNo = F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0켁슆ė
GO
