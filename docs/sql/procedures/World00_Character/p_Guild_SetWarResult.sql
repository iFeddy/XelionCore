CREATE PROCEDURE [dbo].[p_Guild_SetWarResult]

/*
Set Guild War Result
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nWarWinCount int,
@nWarLoseCount int,
@nWarDrawCount int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET 
	nWarWinCount = @nWarWinCount,
	nWarLoseCount = @nWarLoseCount,
	nWarDrawCount = @nWarDrawCount
WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0ఁ௃Ļ
GO
