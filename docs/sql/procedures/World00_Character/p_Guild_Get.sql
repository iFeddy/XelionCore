CREATE PROCEDURE [dbo].[p_Guild_Get]

/*
Get Guild
2006.2 By CJC
*/

@nNo int

AS

SET NOCOUNT ON

SELECT 

	sName,
	sPassword,
	nMaxMembers,
	nMoney,
	nType,
	nGrade,
	nFame,
	nStoneLevel,
	nExp,
	dCreateDate,

	dWarRequestDate,
	dWarStartDate,
	dWarEndDate,
	nWarEnemy,
	nWarWinCount,
	nWarLoseCount,
	nWarDrawCount,
	
	nDismissStatus,
	dDismissDate,

	dNotifyDate,
	sNotifyCharID,
	sNotify,
	
	sIntro
	
FROM tGuild WHERE @nNo = nNo

-- end

0企끪ĝ
GO
