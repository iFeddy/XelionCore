CREATE PROCEDURE [dbo].[p_Guild_Member_GetAll]

/*
Get Guild Member All
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int

AS

SET NOCOUNT ON

SELECT G.nCharNo, G.dJoinDate, G.nGrade, G.nFlags, G.nExpRatio, G.sIntro, C.sID, S.nClass, C.nLevel, C.sLoginZone, C.dLastLoginDate
FROM tGuildMember    AS G INNER JOIN 
     tCharacter      AS C ON C.nCharNo =  G.nCharNo AND C.bDeleted = 0 INNER JOIN 
     tCharacterShape AS S ON S.nCharNo =  G.nCharNo
WHERE G.nNo = @nNo

-- end

0弁ӊħ
GO
