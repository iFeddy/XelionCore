CREATE PROCEDURE [dbo].[p_Guild_Set]

/*
Update Guild Data
2006.3 By blacksky
*/

@nNo				int,

@sName				nvarchar(18),
@sPassword			nvarchar(9),
@nMaxMembers		smallint,
@nMoney				bigint,
@nType				tinyint,
@nGrade				tinyint,	
@nFame				int,
@nStoneLevel		smallint,
@nExp				bigint,
@dCreateDate		datetime,

@dWarRequestDate	datetime,
@dWarStartDate		datetime,
@dWarEndDate		datetime,
@nWarEnemy			int,
@nWarWinCount		int,
@nWarLoseCount		int,
@nWarDrawCount		int,

@nDismissStatus		tinyint,
@dDismissDate		datetime,

@dNotifyDate		datetime,
@sNotifyCharID		nvarchar(16),
@sNotify			nvarchar(512),

@sIntro				nvarchar(128),

@nRet				tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE
	tGuild

SET
	sName			= @sName,
	sPassword		= @sPassword,
	nMaxMembers		= @nMaxMembers,
	nMoney			= @nMoney,
	nType			= @nType,
	nGrade			= @nGrade,
	nFame			= @nFame,
	nStoneLevel		= @nStoneLevel,
	nExp			= @nExp,
	dCreateDate		= @dCreateDate,

	dWarRequestDate	= @dWarRequestDate,
	dWarStartDate	= @dWarStartDate,
	dWarEndDate		= @dWarEndDate,
	nWarEnemy		= @nWarEnemy,
	nWarWinCount	= @nWarWinCount,
	nWarLoseCount	= @nWarLoseCount,
	nWarDrawCount	= @nWarDrawCount,
	
	nDismissStatus	= @nDismissStatus,
	dDismissDate	= @dDismissDate,

	dNotifyDate		= @dNotifyDate,
	sNotifyCharID	= @sNotifyCharID,
	sNotify			= @sNotify,
	
	sIntro			= @sIntro
	
WHERE
	@nNo = nNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END
-- end

0ଁ白Ĳ
GO
