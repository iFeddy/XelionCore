CREATE PROCEDURE [dbo].[p_Item_GetData]
/*
Item Get Data
2004.10 By CJC
 input:
	nItemKey = Item Key
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nItemKey bigint,
@nOwner int OUTPUT,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nItemID int OUTPUT,
@nFlags int OUTPUT,
@dDate datetime OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
SELECT 
	@nOwner = nOwner,
	@nStorageType = nStorageType,
	@nStorage = nStorage,
	@nItemID = nItemID,
	@nFlags = nFlags,
	@dDate = dDate
FROM tItem
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0혁흰œ
GO
