CREATE PROCEDURE [dbo].[p_Char_GetShape]
/*
Character Get shape
2004.6 By CJC
 input:
   nCharNo = To be get shape
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT 
	nRace, 
	nClass, 
	nGender, 
	nHairType, 
	nHairColor, 
	nFaceShape 
FROM tCharacterShape WHERE nCharNo = @nCharNo
-- end

0적궜ś
GO
