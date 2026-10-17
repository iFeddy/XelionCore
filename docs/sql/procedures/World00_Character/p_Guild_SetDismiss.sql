CREATE PROCEDURE [dbo].[p_Guild_SetDismiss]

/*
Set Guild Type
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nDismissStatus tinyint,
@dDismissDate datetime,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nDismissStatus = @nDismissStatus, dDismissDate = @dDismissDate WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0渁龫ĳ
GO
