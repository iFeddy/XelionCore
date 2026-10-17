CREATE PROCEDURE [dbo].[p_Guild_SetNotify]

/*
Set Guild Notify
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@sNotifyCharID nvarchar(16),
@sNotify nvarchar(512),

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET dNotifyDate = getdate(), sNotifyCharID = @sNotifyCharID, sNotify = @sNotify WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0愁⽖ĸ
GO
