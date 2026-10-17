CREATE PROCEDURE [dbo].[p_Item_Create_UseEmptyStorage_tmp]
/*
Item Create
2008.4 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
modified by Adam.
when @nStorage <0,  print the charno + itemkey.
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
		PRINT cast(@nOwner as varchar(20)) + ':' + cast(@nItemKey as varchar(20))
		ROLLBACK TRAN
		RETURN
	END

	-- 아이템 생성
	INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
	VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
	SET @nRet = @@ERROR
	IF @nRet = 0
		BEGIN			
		EXEC p_Item_SetOption @nItemKey, 300,1,0
		COMMIT TRAN
		END
	ELSE
		BEGIN
		PRINT cast(@nOwner as varchar(20)) + ':' + cast(@nItemKey as varchar(20))
		ROLLBACK TRAN
		RETURN
		END
END

0뤁ኻŏ
GO
