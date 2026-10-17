CREATE PROCEDURE [dbo].[p_Char_GetOptGame]

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
@sData varbinary(64) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sGame, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sGame IS NOT NULL

-- end

0౻Ŕ
GO
