CREATE PROCEDURE [dbo].[p_Char_GetTitleAll]

/*
Character title set
2005.10 By CJC
  Input:
    nCharNo = character number
  output: 
    featch to get data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nType, nElementNo, nElementValue 
  FROM tCharacterTitle WITH(READUNCOMMITTED)
 WHERE nCharNo = @nCharNo ORDER BY nNo  

-- end

0矬Ť
GO
