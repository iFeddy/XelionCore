CREATE PROCEDURE [dbo].[p_MasPup_DelAll]
/*
Delete My Master and Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tMasterPupils WHERE (nMaster = @nCharNo) OR (nPupils = @nCharNo)  
SET @nRet = @@ERROR
-- end

0ź
GO
