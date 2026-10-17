CREATE PROCEDURE [dbo].[p_Item_GetStorageListType]

/*
Item Get storage list from StorageType
2004.8 By CJC
 input:
	nCharNo	     = Item Owner
	nStorageType = Item Storage Type
 output:
	none
 desc:
	fetch success, get data
	fetch fail, not data
*/

@nCharNo int,
@nStorageType tinyint

AS

SET NOCOUNT ON

SELECT nStorage FROM tItem WHERE nOwner = @nCharNo AND nStorageType = @nStorageType

-- end

0༁ꝜŠ
GO
