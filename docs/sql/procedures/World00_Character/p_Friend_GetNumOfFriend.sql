CREATE PROCEDURE [dbo].[p_Friend_GetNumOfFriend]

/*
Char Get friend all
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nNumOfFriend int OUTPUT

AS

SET NOCOUNT ON

SET @nNumOfFriend = 0

SELECT @nNumOfFriend = COUNT(*) FROM tFriend AS F WITH(READUNCOMMITTED) 
INNER JOIN tCharacter AS C WITH(READUNCOMMITTED) ON C.nCharNo = F.nFriendCharNo AND C.bDeleted = 0
WHERE F.nCharNo = @nCharNo

-- end

0稁黳Ě
GO
