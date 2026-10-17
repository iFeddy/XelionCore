CREATE PROCEDURE [dbo].[p_MasPup_GetAll]
/*
Get My Master and Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT nMaster, nPupils, nFlags FROM tMasterPupils WHERE (nMaster = @nCharNo) OR (nPupils = @nCharNo)  
-- end

0㸁ſ
GO
