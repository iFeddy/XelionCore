CREATE PROCEDURE [dbo].[p_Char_Wipe]
/*

Character Wipe(Erase)
2006.5 By CJC

 input:
	nCharNo

 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,

@nRet int OUTPUT

AS

SET NOCOUNT ON

DELETE FROM tFriend WHERE nCharNo = @nCharNo			-- 관계에 의해 삭제됨
DELETE FROM tFriend WHERE nFriendCharNo = @nCharNo	-- 친구 테이블 순환 참고 없앰용

DELETE FROM tCharacter WHERE nCharNo = @nCharNo

SET @nRet = @@ERROR

-- end

0눁﷑Ē
GO
