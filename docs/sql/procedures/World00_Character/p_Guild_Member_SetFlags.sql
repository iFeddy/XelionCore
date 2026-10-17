CREATE PROCEDURE [dbo].[p_Guild_Member_SetFlags]

/*
Set Guild Member Flags
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error
	nRet = 1,  OK
*/

@nNo int,
@nCharNo int,
@nFlags int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuildMember SET nFlags = @nFlags WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0뇇ĭ
GO
