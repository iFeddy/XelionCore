CREATE PROCEDURE [dbo].[p_Char_GetOptKeyMapping]

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
@sData varbinary(308) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sKeyMapping, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sKeyMapping IS NOT NULL

-- end

0䠁뾹ŕ
GO
