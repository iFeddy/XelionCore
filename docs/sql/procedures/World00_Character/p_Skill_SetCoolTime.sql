CREATE PROCEDURE [dbo].[p_Skill_SetCoolTime]
/*
Skill Set Cool Time
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillWriteTime bigint,
@nSkillCoolTime int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tSkill SET
	nSkillWriteTime = @nSkillWriteTime,
	nSkillCoolTime	= @nSkillCoolTime
WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
	SET @nRet = 0
ELSE
	SET @nRet = 1
-- end

0度넏Ĕ
GO
