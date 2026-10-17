CREATE PROCEDURE [dbo].[p_Skill_ResetPower]
/*
Skill Reset Power
2005.6 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DECLARE @nFlagsTemp AS int
SET @nRet = 0 -- set error first
SELECT @nFlagsTemp = nFlags FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 0
	RETURN
BEGIN TRAN
SET @nFlagsTemp = @nFlagsTemp & 0x7fffffff -- Clear skillempower_can_reset's Bit
UPDATE tCharacter SET nFlags = @nFlagsTemp WHERE nCharNo = @nCharNo
UPDATE tSkill SET 
	nSkillPowerDemage	= 0,
	nSkillPowerSP		= 0,
	nSkillPowerKeepTime = 0,
	nSkillPowerCoolTime = 0
WHERE nCharNo = @nCharNo
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
COMMIT TRAN
SET @nRet = 1 -- no error return
-- end

0蠁죇Ē
GO
