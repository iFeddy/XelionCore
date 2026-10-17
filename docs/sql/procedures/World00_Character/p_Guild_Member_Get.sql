CREATE PROCEDURE [dbo].[p_Guild_Member_Get]

/*
Get Guild Member
2006.2 By CJC

Input:
	See SQL
	
Output:
	fetch to take data
*/

@nNo int,
@nCharNo int

AS

SET NOCOUNT ON

SELECT G.dJoinDate, G.nGrade, G.nFlags, G.nExpRatio, G.sIntro, C.sID, S.nClass, C.nLevel, C.sLoginZone, C.dLastLoginDate
FROM tGuildMember    AS G INNER JOIN 
     tCharacter      AS C ON C.nCharNo =  G.nCharNo AND C.bDeleted = 0 INNER JOIN 
     tCharacterShape AS S ON S.nCharNo =  G.nCharNo
WHERE G.nCharNo = @nCharNo AND G.nNo = @nNo

-- end

0숁ⷸĨ
GO
