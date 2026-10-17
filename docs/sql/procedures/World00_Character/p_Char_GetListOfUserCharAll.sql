CREATE PROCEDURE [dbo].[p_Char_GetListOfUserCharAll]
/*
Character Get list of user character(include deleted character)
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
SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nUserNo = @nUserNo
-- end

0唁『ő
GO
