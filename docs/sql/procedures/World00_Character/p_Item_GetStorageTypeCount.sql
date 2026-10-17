CREATE PROCEDURE [dbo].[p_Item_GetStorageTypeCount]
/*
Item Get number of item count from StorageType
2008.4 By CJC
 input:
	Item Owner
	Item StorageType
 output:
	Number of item count from StorageType
 */

@nOwner int,
@nStorageType tinyint,
@nItemCount int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SELECT @nItemCount=COUNT(nOwner) FROM tItem WHERE nOwner = @nOwner AND nStorageType = @nStorageType
END

0
GO
