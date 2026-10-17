CREATE PROCEDURE [dbo].[p_Item_GetStorage]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
 output:
	1, ok nStorageType, nStorage = valid value
	0, Not found item
*/
@nItemKey bigint,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nRet tinyint = 1 OUTPUT 
AS
SET NOCOUNT ON
SELECT @nStorageType = nStorageType, @nStorage = nStorage FROM tItem WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
	SET @nRet = 0
-- end

0⬁훋Ŝ
GO
