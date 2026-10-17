CREATE PROCEDURE [dbo].[p_Char_GetTitleCount]

/*
Character number of title
2005.10 By CJC
  Input:
    nCharNo = character number
  output: 
    0 = none
    ? = number of title on this character
*/

@nCharNo int,
@nCount int OUTPUT

AS

SET NOCOUNT ON

SET @nCount = 0

SELECT @nCount = COUNT(*) FROM tCharacterTitle(NOLOCK) WHERE nCharNo = @nCharNo

-- end

0줁䋢Ť
GO
