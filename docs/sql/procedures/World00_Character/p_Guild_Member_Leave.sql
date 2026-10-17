CREATE PROCEDURE [dbo].[p_Guild_Member_Leave]

/*
Get Guild Member Leave
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0, Error
	nRet = 1, OK
*/

@nNo int,
@nCharNo int,

@nRet tinyint output
AS

SET NOCOUNT ON

SET @nRet = 0

-- 삭제
DELETE FROM tGuildMember WHERE nNo = @nNo AND nCharNo = @nCharNo
SET @nRet = @@ROWCOUNT

-- end

0刁鑵ī
GO
