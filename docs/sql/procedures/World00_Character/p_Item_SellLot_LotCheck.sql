CREATE PROCEDURE [dbo].[p_Item_SellLot_LotCheck]
/*
Item sell - change lot and add money
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found item
	        2 = Missmatch item owner
	        3 = Money add error
	        ? = Error code(Lot change error)
*/
@nItemKey bigint,
@nOwner int,
@nOptionType smallint,
@nOptionData bigint,
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

DECLARE @nGetOptionDataLot bigint
DECLARE @nGetOptionRet int
EXEC p_Item_GetOption @nItemKey, @nOptionType, @nGetOptionDataLot OUTPUT, @nGetOptionRet OUTPUT
IF @nGetOptionRet <> 1
BEGIN
	SET @nRet = 4
	RETURN
END

IF @nGetOptionDataLot < @nOptionData
BEGIN
	SET @nRet = 5
	RETURN
END

BEGIN TRAN
-- 아이템 소유자가 돈 증가
EXEC p_Char_MoneyAdd @nOwner, @nPrice, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	ROLLBACK TRAN
	SET @nRet = 3
	RETURN
END
-- 아이템 랏 저장프로시져 호출
EXEC p_Item_SetOption @nItemKey, @nOptionType, @nOptionData, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0ခ㲢ũ
GO
