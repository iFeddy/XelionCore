CREATE PROCEDURE [dbo].[p_Item_GetOption]
/*
Item Get Option
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = No data found
	        ? = Number data here(must to be 1)
	        
	@nOptionData = option data
		        
*/
@nItemKey bigint,
@nOptionType smallint,
-- Output var
@nOptionData bigint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SELECT @nOptionData = nOptionData FROM tItemOptions WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
SET @nRet = @@ROWCOUNT
-- end

0輁른Ś
GO
