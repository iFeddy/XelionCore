CREATE PROCEDURE [dbo].[p_Skill_GetNumOfSkill]
/*
Skill Get Num of skill
2005.5 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int,
@nNumOfSkill int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfSkill = 0
SELECT @nNumOfSkill = COUNT(*) FROM tSkill WHERE nCharNo = @nCharNo
-- end

0企풣đ
GO
