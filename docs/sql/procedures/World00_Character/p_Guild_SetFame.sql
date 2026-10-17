CREATE PROCEDURE [dbo].[p_Guild_SetFame]

/*
Set Guild Fame
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nFame int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nFame = @nFame WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0蟳ĵ
GO
