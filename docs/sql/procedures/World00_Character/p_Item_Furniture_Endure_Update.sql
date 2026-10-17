CREATE PROCEDURE [dbo].[p_Item_Furniture_Endure_Update]
/*
Item Furniture Endure Update
2008.10 By Lunar
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Not enough money to sub
	        4 = Upgrade value setting, upgrade fail count error(option set error)
	        5 = Raw delete fail
	        6 = RawLeft delete fail, if valid item
	        7 = RawRight delete fail, if valid item
	        8 = RawMiddle delete fail, if valid item
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nEndureKitKey bigint,
@nEndureKitOptionType smallint,
@nEndureKitLot smallint,
@nFurnitureKey bigint,
@nFurnitureEndureOptionType smallint,
@nFurnitureEndureOptionData bigint,
@nRet int OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTinyintRet tinyint
	DECLARE @nRetTemp int

	-- 1.1 내구도 키트 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nEndureKitKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END
	
	-- 1.2 가구 아이템 소유자 확인
	EXEC p_Item_IsValidOwner @nFurnitureKey, @nOwner, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 0 BEGIN
		SET @nRet = @nTinyintRet
		RETURN
	END

	BEGIN TRAN
	
	-- 2.1 사용된 내구도 키트 아이템 지우기 (사용후 Lot == 0 인 경우)
	IF @nEndureKitLot = 0 BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_Delete @nEndureKitKey, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	-- 2.2 사용된 내구도 키트 아이템 차감하기 (사용후 Lot > 0 인 경우)
	ELSE BEGIN
		SET @nRetTemp = 1
		EXEC p_Item_SetOption @nEndureKitKey, @nEndureKitOptionType, @nEndureKitLot, @nRetTemp OUTPUT
		IF @nRetTemp <> 0 BEGIN
			ROLLBACK TRAN
			SET @nRet = 5
			RETURN
		END
	END
	
	-- 3 가구 아이템 내구도 상승
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nFurnitureKey, @nFurnitureEndureOptionType, @nFurnitureEndureOptionData, @nRetTemp OUTPUT
	IF @nRetTemp <> 0 BEGIN
		ROLLBACK TRAN
		SET @nRet = 5
		RETURN
	END
	
	-- 업글 성공
	COMMIT TRAN
	SET @nRet = 0

	SET NOCOUNT OFF
END

0鴁Œ
GO
