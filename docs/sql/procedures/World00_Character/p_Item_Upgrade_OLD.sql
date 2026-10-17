CREATE PROCEDURE [dbo].[p_Item_Upgrade_OLD]
/*
Item upgrade
2005.6 By CJC
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
	        
 operation:
	1. money sub
	2. upgrade item level
	3. delete raws
*/
@nOwner int,
@nItemKey bigint,
@nItemDestroy int,
@nUpgradeOptionType smallint,
@nUpgradeOptionValue bigint,
@nUpgradeFailCountType smallint,
@nUpgradeFailCountValue bigint,
@nRawKey bigint,
@nRawLeftKey bigint,
@nRawRightKey bigint,
@nGiftMoney bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKey, @nOwner, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 아이템 소유자의 돈 감소
BEGIN TRAN
IF @nGiftMoney > 0 -- 수고비가 있는 경우, 금액을 빼준다.
BEGIN
	EXEC p_Char_MoneySub @nOwner, @nGiftMoney, @nTinyintRet OUTPUT
	IF @nTinyintRet <> 1
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 3
		RETURN
	END
END	
-- 아이템 업그레이드 값 조정
DECLARE @nRetTemp int
IF @nItemDestroy = 1 -- 업글 실패로 아이템 파괴해야 되는경우
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nItemKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
ELSE
BEGIN
	-- 업글 실패로 파괴가 아닌 경우, 새로운 업글값, 실패 카운트 기록	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeOptionType, @nUpgradeOptionValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
	
	SET @nRetTemp = 1
	EXEC p_Item_SetOption @nItemKey, @nUpgradeFailCountType, @nUpgradeFailCountValue, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 4
		RETURN
	END
END
-- 사용된 자원 지우기
SET @nRetTemp = 1
EXEC p_Item_Delete @nRawKey, @nRetTemp OUTPUT
IF @nRetTemp <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 5
	RETURN
END
-- 사용된 자원왼족 지우기(입력된 값이 있으면)
IF @nRawLeftKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawLeftKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 6
		RETURN
	END
END
-- 사용된 자원오른쪽 지우기(입력된 값이 있으면)
IF @nRawRightKey <> 0
BEGIN
	SET @nRetTemp = 1
	EXEC p_Item_Delete @nRawRightKey, @nRetTemp OUTPUT
	IF @nRetTemp <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = 7
		RETURN
	END
END
-- 업글 성공
COMMIT TRAN
SET @nRet = 0
-- end

0锁膎v
GO
