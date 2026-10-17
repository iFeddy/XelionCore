CREATE PROCEDURE [dbo].[p_Item_SetStorage]
/*
Get item storage
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nStorageType, nStorage = Storage new value
 output:
	1 = OK
	0 = Not found or error
*/
@nItemKey bigint,
@nStorageType tinyint,
@nStorage smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
UPDATE tItem 
SET	nStorageType = @nStorageType, nStorage = @nStorage
WHERE nItemKey = @nItemKey
SET @nRet = @@ROWCOUNT
-- end

0ꀁꌞŬ
GO
