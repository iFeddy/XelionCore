CREATE PROCEDURE [dbo].[p_Guild_Member_GetCount]

/*
Get Guild Member Count
2006.2 By CJC

Input:
	See SQL
	
Output:
	nCount = Number of member
*/

@nNo int,
@nCount int output
AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tGuildMember AS G WITH(READUNCOMMITTED) 
INNER JOIN tCharacter AS C WITH(READUNCOMMITTED) ON C.nCharNo = G.nCharNo AND C.bDeleted = 0
WHERE G.nNo = @nNo 

-- end

0갬ĩ
GO
