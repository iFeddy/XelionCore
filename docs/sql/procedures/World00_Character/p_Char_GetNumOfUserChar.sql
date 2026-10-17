CREATE PROCEDURE [dbo].[p_Char_GetNumOfUserChar]
/*
Character number of user character
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
    0 = none
    ? = number of user's character, execpt deleted.
*/
@nUserNo int,
@nNumOfUserChar int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfUserChar = 0
SELECT @nNumOfUserChar = COUNT(nCharNo) FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo AND bDeleted = 0
-- end

0윁ᡖœ
GO
