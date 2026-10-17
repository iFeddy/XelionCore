CREATE PROCEDURE [dbo].[p_Guild_Member_SetExpRatio]

/*
Set Guild Member Exp Ratio
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error
	nRet = 1,  OK
*/

@nNo int,
@nCharNo int,
@nExpRatio smallint,

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuildMember SET nExpRatio = @nExpRatio WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

0謁袙Ĭ
GO
