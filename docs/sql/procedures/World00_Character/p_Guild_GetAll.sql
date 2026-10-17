CREATE PROCEDURE [dbo].[p_Guild_GetAll]

/*
Char Get All Guild
2006.2 By CJC
*/

AS

SET NOCOUNT ON

SELECT

	nNo,
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
	
FROM tGuild

-- end

0蠁꒎Ğ
GO
