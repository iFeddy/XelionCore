CREATE PROCEDURE [dbo].[p_Item_Create_UseEmptyStorage]
/*
Item Create
2008.4 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint OUTPUT,
@nItemID int,
@nFlags int,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	BEGIN TRAN

	-- 빈인벤 위치찾기
	EXEC  p_Item_GetStorageEmpty @nOwner, @nStorageType, @nStorage OUTPUT
	IF @nStorage < 0
	BEGIN
		SET @nRet = -1
		ROLLBACK TRAN
		RETURN
	END

	-- 아이템 생성
	INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
	VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
	SET @nRet = @@ERROR
	IF @nRet = 0
		COMMIT TRAN
	ELSE
		ROLLBACK TRAN
END

0퐁被Ŏ
GO
