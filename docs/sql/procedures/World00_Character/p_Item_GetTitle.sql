CREATE PROCEDURE [dbo].[p_Item_GetTitle]
/*
Item Get Title
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data fount
	        ? = Number data here(must to be 1)
	        
	@nTitleType,@nTitleCharNo,@sTitle = title data
		        
*/
@nItemKey bigint,
-- Output var
@nTitleType smallint OUTPUT,
@nTitleCharNo int OUTPUT,
@sTitle nvarchar(24) OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT 
@nTitleType = nTitleType,
@nTitleCharNo = nTitleCharNo,
@sTitle = sTitle
FROM tItemTitle WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ⴁ▐Ţ
GO
