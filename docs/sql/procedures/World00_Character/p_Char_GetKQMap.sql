CREATE PROCEDURE [dbo].[p_Char_GetKQMap]
/*
Character Get Kingdom quest map
2006.4 By CJC
 input:
   nCharNo = To be get KQ Map data
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

nKQHandle,
sKQMap,
nKQX,
nKQY,
nKQD,
dKQDate

FROM tCharacter

WHERE nCharNo = @nCharNo

-- end

0䟅ŏ
GO
