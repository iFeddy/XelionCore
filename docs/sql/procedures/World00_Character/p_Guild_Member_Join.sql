CREATE PROCEDURE [dbo].[p_Guild_Member_Join]

/*
Get Guild Member Join
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0, Error
	nRet = 1, OK
*/

@nNo int,
@nCharNo int,
@nGrade tinyint,

@nRet tinyint output
AS

SET NOCOUNT ON

SET @nRet = 0

-- 만들기
INSERT tGuildMember (  nNo,  nCharNo,  nGrade )
VALUES              ( @nNo, @nCharNo, @nGrade )

-- 오류 확인
IF @@ERROR <> 0 
BEGIN
	RETURN
END

SET @nRet = 1

-- end

RETURN

0䌁핛Ī
GO
