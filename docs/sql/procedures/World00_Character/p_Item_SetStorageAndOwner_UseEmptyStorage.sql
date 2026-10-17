CREATE PROCEDURE [dbo].[p_Item_SetStorageAndOwner_UseEmptyStorage]
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
@nStorage smallint OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON


	-- 빈인벤 위치찾기
	EXEC  p_Item_GetStorageEmpty @nOwner, @nStorageType, @nStorage OUTPUT
	IF @nStorage < 0
	BEGIN
		SET @nRet = -1
		ROLLBACK TRAN
		RETURN
	END

	UPDATE tItem 
	SET	nOwner = @nOwner, nStorageType = @nStorageType, nStorage = @nStorage
	WHERE nItemKey = @nItemKey

	SET @nRet = @@ROWCOUNT
-- end

0昁Ů
GO
