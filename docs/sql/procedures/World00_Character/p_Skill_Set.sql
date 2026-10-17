CREATE PROCEDURE [dbo].[p_Skill_Set]
/*
Skill Set
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillLevel tinyint,
@nSkillExp int,
@nSkillWriteTime bigint,
@nSkillCoolTime int,
@nSkillPowerDemage tinyint,
@nSkillPowerSP tinyint,
@nSkillPowerKeepTime tinyint,
@nSkillPowerCoolTime tinyint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tSkill SET
	nSkillLevel			= @nSkillLevel,
	nSkillExp			= @nSkillExp,
	nSkillWriteTime		= @nSkillWriteTime,
	nSkillCoolTime		= @nSkillCoolTime,
	nSkillPowerDemage	= @nSkillPowerDemage,
	nSkillPowerSP		= @nSkillPowerSP,
	nSkillPowerKeepTime = @nSkillPowerKeepTime,
	nSkillPowerCoolTime = @nSkillPowerCoolTime
WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tSkill (  nCharNo,  nSkillNo,  nSkillLevel,  nSkillExp,  nSkillWriteTime,  nSkillCoolTime,  nSkillPowerDemage,  nSkillPowerSP,  nSkillPowerKeepTime,  nSkillPowerCoolTime )
	VALUES        ( @nCharNo, @nSkillNo, @nSkillLevel, @nSkillExp, @nSkillWriteTime, @nSkillCoolTime, @nSkillPowerDemage, @nSkillPowerSP, @nSkillPowerKeepTime, @nSkillPowerCoolTime )
END
SET @nRet = @@ERROR
-- end

0섁볫ē
GO
