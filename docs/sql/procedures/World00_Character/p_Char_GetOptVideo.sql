CREATE PROCEDURE [dbo].[p_Char_GetOptVideo]

/*
Character Get options
2005.8 By CJC

 input:
   nCharNo = To be get option data
 output:
	@sData = Option binary data
*/

@nCharNo int,
@nRet tinyint output,
@sData varbinary(60) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sVideo, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sVideo IS NOT NULL

-- end

0Ⰱ遊ř
GO
