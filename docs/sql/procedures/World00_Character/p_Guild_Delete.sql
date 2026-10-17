CREATE PROCEDURE [dbo].[p_Guild_Delete]

/*
Char Delete Guild
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@nRet tinyint output

AS

SET NOCOUNT ON

SET @nRet = 0

DELETE FROM tGuild WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

-- end

0ᘁ뱆Ĝ
GO
