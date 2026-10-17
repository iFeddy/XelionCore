CREATE PROCEDURE [dbo].[p_Item_SellAll_LotCheck]
/*
Item sell all
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code
	        
 operation:
	1. money add
	2. delete sell item
	
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nLot bigint,
@nPrice bigint,
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

-- 1보다 큰 경우에만 Lot 을 체크해 본다.
-- 1일 경우는 어차피 tItem 테이블에 존재하는 것만으로도 옳은 것으로 취급한다.
IF @nLot > 1
BEGIN
	DECLARE @nGetOptionDataLot bigint
	DECLARE @nGetOptionRet int
	EXEC p_Item_GetOption @nItemKey, @nOptionType, @nGetOptionDataLot OUTPUT, @nGetOptionRet OUTPUT
	IF @nGetOptionRet <> 1
	BEGIN
		SET @nRet = 4
		RETURN
	END

	IF @nGetOptionDataLot < @nLot
	BEGIN
		SET @nRet = 5
		RETURN
	END
END

-- 아이템 소유자의 돈 증가
BEGIN TRAN
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nTinyintRet OUTPUT
IF @nTinyintRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 삭제 저장프로시져 호출
EXEC p_Item_Delete @nItemKey, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0䨁Ŧ
GO
