CREATE PROCEDURE [dbo].[p_Skill_GetAll]
/*
Skill Get All
2004.10 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT	nSkillNo, nSkillLevel, nSkillExp, nSkillWriteTime, nSkillCoolTime, 
		nSkillPowerDemage, nSkillPowerSP, nSkillPowerKeepTime, nSkillPowerCoolTime
FROM tSkill WHERE nCharNo = @nCharNo
-- end

0樁䪓đ
GO
