CREATE PROCEDURE [dbo].[p_Char_GetOptSound]

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
@sData varbinary(1) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sSound, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sSound IS NOT NULL

-- end

0鰥Ř
GO
