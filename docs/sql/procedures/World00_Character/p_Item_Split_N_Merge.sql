CREATE PROCEDURE [dbo].[p_Item_Split_N_Merge]
/*
Item Split and Merge - change lot
2004.9 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
	        ? = Error code(Lot change/Delete item error)
*/
@nItemKeyFrom bigint,
@nOwnerFrom int,
@nOptionTypeFrom smallint,
@nOptionDataFrom bigint,
@nItemKeyTo bigint,
@nOwnerTo int,
@nOptionTypeTo smallint,
@nOptionDataTo bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 수량 줄어들 아이템(From) 소유자 확인
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END
-- 수량 늘어날 아이템(To) 소유자 확인
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END
BEGIN TRAN
-- 수량 줄어들 아이템 랏 변경저장프로시져 호출
EXEC p_Item_SetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 수량 늘어날 아이템 랏 변경저장프로시져 호출
EXEC p_Item_SetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	RETURN
END
-- 완료
COMMIT TRAN
-- end

0봁染ű
GO
