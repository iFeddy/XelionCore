CREATE PROCEDURE [dbo].[p_Item_SetStorageAndOwner]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nOwner, nStorageType, nStorage = Storage new value
 output:
	1 = OK
	0 = Not found or error
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tItem 
SET	nOwner = @nOwner, nStorageType = @nStorageType, nStorage = @nStorage
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ⴁŗŮ
GO
