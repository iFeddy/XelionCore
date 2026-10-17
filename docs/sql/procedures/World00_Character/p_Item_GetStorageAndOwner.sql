CREATE PROCEDURE [dbo].[p_Item_GetStorageAndOwner]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
 output:
	1, ok nOwner, nStorageType, nStorage = valid value
	0, Not found item
*/
@nItemKey bigint,
@nOwner int OUTPUT,
@nStorageType tinyint OUTPUT,
@nStorage smallint OUTPUT,
@nRet tinyint OUTPUT 
AS
SET NOCOUNT ON
SELECT @nOwner = nOwner, @nStorageType = nStorageType, @nStorage = nStorage FROM tItem WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
	SET @nRet = 0
ELSE
	SET @nRet = 1

-- end

0ခ惛ŝ
GO
