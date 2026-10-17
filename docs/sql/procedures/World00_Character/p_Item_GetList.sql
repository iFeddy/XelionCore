CREATE PROCEDURE [dbo].[p_Item_GetList]
/*
Item Get list
2004.8 By CJC
 input:
	nCharNo	     = Item Owner
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/
@nCharNo int
AS
SET NOCOUNT ON
SELECT nItemKey, nStorageType, nStorage, nItemID, nFlags, dDate
FROM tItem
WHERE nOwner = @nCharNo
-- end

0持㖩ŕ
GO
