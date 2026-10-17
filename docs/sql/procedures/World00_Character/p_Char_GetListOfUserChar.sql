CREATE PROCEDURE [dbo].[p_Char_GetListOfUserChar]
/*
Character Get list of user character
2004.6 By CJC
  Input:
    nUserNo = Search user's character
  output: 
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nUserNo int
AS
SET NOCOUNT ON
SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo AND bDeleted = 0
-- end

0ᰁ㯪Ő
GO
