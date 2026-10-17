CREATE PROCEDURE [dbo].[p_Skill_Del]
/*
Skill Delete
2004.10 By CJC
 input:
	nCharNo, nSkillNo = To be delete info.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nSkillNo smallint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tSkill WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
SET @nRet = @@ERROR
-- end

0ꐁĎ
GO
