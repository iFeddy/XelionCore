CREATE PROCEDURE [dbo].[p_MasPup_DelMaster]
/*
Delete My Master
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,   -- As Pupils
@nRet int OUTPUT
AS
SET NOCOUNT ON
DELETE FROM tMasterPupils WHERE nPupils = @nCharNo
SET @nRet = @@ERROR
-- end

0찁ʴž
GO
