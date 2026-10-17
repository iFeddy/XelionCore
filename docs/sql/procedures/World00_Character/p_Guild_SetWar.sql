CREATE PROCEDURE [dbo].[p_Guild_SetWar]

/*
Set Guild War
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@dWarRequestDate datetime,
@dWarStartDate datetime,
@dWarEndDate datetime,
@nWarEnemy int,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET 
	dWarRequestDate = @dWarRequestDate,
	dWarStartDate = @dWarStartDate,
	dWarEndDate = @dWarEndDate,
	nWarEnemy = @nWarEnemy
WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0态痗Ļ
GO
