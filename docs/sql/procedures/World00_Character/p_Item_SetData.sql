CREATE PROCEDURE [dbo].[p_Item_SetData]
/*
Item Set data
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, 
	        0 = Error
*/
@nItemKey bigint,
@nOwner int,
@nStorageType tinyint,
@nStorage smallint,
@nItemID int,
@nFlags int,
@dDate datetime,
@nRet int OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tItem SET 
nOwner = @nOwner,
nStorageType = @nStorageType,
nStorage = @nStorage,
nItemID = @nItemID,
nFlags = @nFlags,
dDate = @dDate
WHERE nItemKey = @nItemKey
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end

0省➓i
GO
