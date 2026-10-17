CREATE PROCEDURE [dbo].[p_Char_GetNumOfUserCharAll]
/*
Character number of user character(include deleted character)
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
    0 = none
    ? = number of user's character
*/
@nUserNo int,
@nNumOfUserChar int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfUserChar = 0
SELECT @nNumOfUserChar = COUNT(nCharNo) FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo
-- end

0⨁䆅Ŕ
GO
