CREATE PROCEDURE [dbo].[p_Char_GetOptWindowsPos]

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
@sData varbinary(392) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sWindowsPos, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sWindowsPos IS NOT NULL

-- end

0輁른Ś
GO
