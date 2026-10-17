CREATE PROCEDURE [dbo].[p_Item_SellAll]
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

0攁怵Ŧ
GO
