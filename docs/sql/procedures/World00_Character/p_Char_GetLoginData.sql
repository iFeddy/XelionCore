CREATE PROCEDURE [dbo].[p_Char_GetLoginData]
/*
Character Get login data
2004.6 By CJC
 input:
	nCharNo = To be get login data
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
	C.sID, 
	C.nLevel, 
	C.nSlotNo, 
	C.sLoginZone, 
	C.bDeleted, 
	C.dDeletedDate,
	
	S.nRace, 
	S.nClass, 
	S.nGender, 
	S.nHairType, 
	S.nHairColor, 
	S.nFaceShape,
	
	C.nKQHandle,
	C.sKQMap,
	C.nKQX,
	C.nKQY,
	C.nKQD,
	C.dKQDate
FROM tCharacter AS C
INNER JOIN tCharacterShape AS S ON S.nCharNo = @nCharNo
WHERE C.nCharNo = @nCharNo
-- end

0踁␲Œ
GO
