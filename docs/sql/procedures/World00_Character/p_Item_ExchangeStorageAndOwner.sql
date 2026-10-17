CREATE PROCEDURE [dbo].[p_Item_ExchangeStorageAndOwner]
/*
Exchange item storage
2004.8 By CJC
 input:
	nOwner1,nItemKey1,nStorageType1,nStoage1 = #1 Item
	nOwner2,nItemKey2,nStorageType2,nStoage2 = #2 Item
 output:
	0 = Ok
	1~3 = Error step No.
 operation:
	Set nItemKey1 To nOwner1,nStorageType1,nStoage1
	Set nItemKey2 To nOwner2,nStorageType2,nStoage2
	
	1. #1 item move to temp storage   * #1,#2 is already input that exchanged data.
	2. set #2's storage to #2
	3. set #1's storage to #1
*/
@nItemKey1 bigint,
@nOwner1 int,
@nStorageType1 tinyint,
@nStorage1 smallint,
@nItemKey2 bigint,
@nOwner2 int,
@nStorageType2 tinyint,
@nStorage2 smallint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
BEGIN TRAN
-- Step 1
UPDATE tItem SET nStorageType = 11
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 1
	ROLLBACK TRAN
	RETURN
END
	
-- Step 2
UPDATE tItem SET nStorageType = @nStorageType2, nStorage = @nStorage2, nOwner = @nOwner2
WHERE nItemKey = @nItemKey2
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END
	
-- Step 3
UPDATE tItem SET nStorageType = @nStorageType1, nStorage = @nStorage1, nOwner = @nOwner1
WHERE nItemKey = @nItemKey1
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END
COMMIT TRAN
-- end

0搁ő
GO
