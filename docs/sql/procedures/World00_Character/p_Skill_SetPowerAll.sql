CREATE PROCEDURE [dbo].[p_Skill_SetPowerAll]
/*
Skill Set Power All
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillPowerDemage tinyint,
@nSkillPowerSP tinyint,
@nSkillPowerKeepTime tinyint,
@nSkillPowerCoolTime tinyint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	IF @nSkillPowerDemage > 5 OR @nSkillPowerSP > 5 OR @nSkillPowerKeepTime > 5 OR @nSkillPowerCoolTime > 5 BEGIN
		SET @nRet = 0
		RETURN
	END
	
	UPDATE tSkill SET 
		nSkillPowerDemage	= @nSkillPowerDemage,
		nSkillPowerSP		= @nSkillPowerSP,
		nSkillPowerKeepTime	= @nSkillPowerKeepTime,
		nSkillPowerCoolTime	= @nSkillPowerCoolTime
	WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
		SET @nRet = 0
	ELSE
		SET @nRet = 1
	-- end
END

0
GO
