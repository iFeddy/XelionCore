CREATE PROCEDURE [dbo].[p_Skill_SetPower]
/*
Skill Set Power
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nSkillNo smallint,
@nSkillPowerSlot tinyint,
@nSkillPowerValue tinyint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	IF @nSkillPowerValue > 5 BEGIN
		SET @nRet = 0
		RETURN
	END
	
	-- nSkillPowerDemage
	IF @nSkillPowerSlot = 0
	BEGIN
		UPDATE tSkill SET nSkillPowerDemage = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerSP
	IF @nSkillPowerSlot = 1
	BEGIN
		UPDATE tSkill SET nSkillPowerSP = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerKeepTime
	IF @nSkillPowerSlot = 2
	BEGIN
		UPDATE tSkill SET nSkillPowerKeepTime = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- nSkillPowerCoolTime
	IF @nSkillPowerSlot = 3
	BEGIN
		UPDATE tSkill SET nSkillPowerCoolTime = @nSkillPowerValue WHERE nCharNo = @nCharNo AND nSkillNo = @nSkillNo
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0
			SET @nRet = 0
		ELSE
			SET @nRet = 1
		RETURN
	END
	-- ERROR
	SET @nRet = 0;
	-- end
END

0
GO
