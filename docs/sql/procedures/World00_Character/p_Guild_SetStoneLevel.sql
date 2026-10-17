CREATE PROCEDURE [dbo].[p_Guild_SetStoneLevel]

/*
Set Guild StoneLevel
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nStoneLevel smallint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET nStoneLevel = @nStoneLevel WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0ﴁ䲨ĺ
GO
