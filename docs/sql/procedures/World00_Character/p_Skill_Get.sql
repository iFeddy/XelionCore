CREATE PROCEDURE [dbo].[p_Skill_Get]
/*
Skill Get
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillLevel tinyint OUTPUT,
@nSkillExp int OUTPUT,
@nSkillWriteTime bigint OUTPUT,
@nSkillCoolTime int OUTPUT,
@nSkillPowerDemage tinyint OUTPUT,
@nSkillPowerSP tinyint OUTPUT,
@nSkillPowerKeepTime tinyint OUTPUT,
@nSkillPowerCoolTime tinyint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT 
	@nSkillLevel		= nSkillLevel,
	@nSkillExp			= nSkillExp,
	@nSkillWriteTime	= nSkillWriteTime,
	@nSkillCoolTime		= nSkillCoolTime,
	@nSkillPowerDemage	= nSkillPowerDemage,
	@nSkillPowerSP		= nSkillPowerSP,
	@nSkillPowerKeepTime = nSkillPowerKeepTime,
	@nSkillPowerCoolTime = nSkillPowerCoolTime
FROM tSkill WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
SET @nRet = @@ROWCOUNT
-- end

0㄁噯Đ
GO
