CREATE PROCEDURE [dbo].[p_Friend_GetAll]

/*
Char Get friend all
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT F.nFriendCharNo, F.nFlags, C.sID
FROM tFriend AS F
INNER JOIN tCharacter AS C ON C.nCharNo =  F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0㈁Ę
GO
