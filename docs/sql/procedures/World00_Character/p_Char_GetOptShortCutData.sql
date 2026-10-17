CREATE PROCEDURE [dbo].[p_Char_GetOptShortCutData]

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
@sData varbinary(1024) output

AS

SET NOCOUNT ON

SET @nRet = 0
SELECT @sData = sShortCutData, @nRet = 1 FROM tCharacterOptions WHERE nCharNo = @nCharNo AND sShortCutData IS NOT NULL

-- end

0脁돝Ŗ
GO
