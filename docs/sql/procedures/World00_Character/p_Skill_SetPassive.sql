CREATE PROCEDURE [dbo].[p_Skill_SetPassive]
/*
Skill Set Passive
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/

@nCharNo int,
@nSkillNewNo smallint,
@bSkillOld tinyint,
@nSkillOldNo smallint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

IF @bSkillOld = 1
BEGIN
	UPDATE tSkill 
	SET nSkillNo = @nSkillNewNo
	WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillOldNo
	
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
		SET @nRet = 1 -- ERROR, No existed old skill
	ELSE
		SET @nRet = 0
		
	RETURN
END

INSERT tSkill (  nCharNo,     nSkillNo )
VALUES        ( @nCharNo, @nSkillNewNo )

SET @nRet = @@ERROR

-- end

0㌁ꔴĕ
GO
