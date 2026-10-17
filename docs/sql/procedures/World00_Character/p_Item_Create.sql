CREATE PROCEDURE [dbo].[p_Item_Create]
/*
Item Create
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 생성
INSERT tItem (  nItemKey,  nOwner,  nStorageType,  nStorage,  nItemID,  nFlags )
VALUES       ( @nItemKey, @nOwner, @nStorageType, @nStorage, @nItemID, @nFlags )
SET @nRet = @@ERROR
-- end
RETURN

0䜁⩳ō
GO
