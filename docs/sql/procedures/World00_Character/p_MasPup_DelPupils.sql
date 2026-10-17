CREATE PROCEDURE [dbo].[p_MasPup_DelPupils]
/*
Delete My Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
			? = Error code
*/
@nCharNo int,   -- As Master
@nPupils int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
DELETE FROM tMasterPupils WHERE (nPupils = @nPupils) AND (nMaster = @nCharNo)
SET @nRet = @@ERROR
-- end

0ꈁ춪Ž
GO
