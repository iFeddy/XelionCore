CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelongedExchange]
/*
Exchange item storage
2009.9.23 By kodong
 input:
	nItemKey1,nStorageType1,nStoage1, nOptionType1, nOptionData1 = #1 Item
	nItemKey2,nStorageType2,nStoage2 = #2 Item
'	output Parameter
'		@nRet	 : check SQL execute(valid/invalid)
'			= 0	 : works ok
'			< 0	 : error
'				-2592 : db error(insert)
'				-2593 : db error(update)
 operation:
	Set nItemKey1 To nStorageType1,nStoage1
	Set nItemKey2 To nStorageType2,nStoage2
	
	1. #1 item move to temp storage   * #1,#2 is already input that exchanged data.
	2. set #2's storage to #2
	3. set #1's storage to #1
	4. set #1's option #1
*/
@nItemKey1     bigint,
@nStorageType1 tinyint,
@nStorage1     smallint,
@nItemKey2     bigint,
@nStorageType2 tinyint,
@nStorage2     smallint,
@nOptionType1  smallint,
@nOptionData1  bigint,
@nRet          int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	BEGIN TRAN

	-- Step 1 (nStorageType = 11 은 swap 용도로 사용하는 슬롯)
	UPDATE tItem SET nStorageType = 11
	WHERE nItemKey = @nItemKey1
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END
		
	-- Step 2
	UPDATE tItem SET nStorageType = @nStorageType2, nStorage = @nStorage2
	WHERE nItemKey = @nItemKey2
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END
		
	-- Step 3
	UPDATE tItem SET nStorageType = @nStorageType1, nStorage = @nStorage1
	WHERE nItemKey = @nItemKey1
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END

	-- Step 4 (착용시 귀속 옵션을 찾아 있으면 업데이트, 없으면 새로 생성)
	IF EXISTS(SELECT nItemKey FROM tItemOptions WHERE nItemKey = @nItemKey1 AND nOptionType = @nOptionType1) BEGIN
		UPDATE tItemOptions
		SET nOptionType = @nOptionType1,
			nOptionData = @nOptionData1
		WHERE nItemKey = @nItemKey1
		AND nOptionType = @nOptionType1
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -2593
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nItemKey1, @nOptionType1, @nOptionData1)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2592
			ROLLBACK TRAN
			RETURN
		END
	END
	
	COMMIT TRAN

	SET NOCOUNT OFF
END

0퀁埌Ġ
GO
