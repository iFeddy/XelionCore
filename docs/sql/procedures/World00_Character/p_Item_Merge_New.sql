CREATE PROCEDURE [dbo].[p_Item_Merge_New]
/*
Item Merge - change lot and delete item
2007.11 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Not found from item
	        2 = Missmatch from item owner
	        3 = Not found to item
	        4 = Missmatch to item owner
                     5 = Get Sub Lot Failed
                     6 = Get Existed Lot Failed
                     7 = Set new lot failed
                     8 = Delete item failed
*/
@nOwnerFrom int,
@nOwnerTo int,
@nItemKeyFrom bigint,
@nItemKeyTo bigint,
@nOptionType smallint,
@nOptionDataLotSub bigint OUTPUT,
@nOptionDataLotResult bigint OUTPUT,
@nRet int OUTPUT
AS
SET NOCOUNT ON

-- »èÁ¦µÉ ¾ÆÀÌÅÛ(From) ¼ÒÀ¯ÀÚ È®ÀÎ
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

-- ¶ù º¯°æµÉ ¾ÆÀÌÅÛ(To) ¼ÒÀ¯ÀÚ È®ÀÎ
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ù ¾ò±â
EXEC p_Item_GetOption @nItemKeyFrom, @nOptionType, @nOptionDataLotSub OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 5
	RETURN
END

-- ´õÇØÁú ¾ÆÀÌÅÛ ¶ù ¾ò±â
DECLARE @nOptionDataLotAdd bigint
EXEC p_Item_GetOption @nItemKeyTo, @nOptionType, @nOptionDataLotAdd OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 6
	RETURN
END

-- ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
BEGIN TRAN
SET @nOptionDataLotResult = @nOptionDataLotSub + @nOptionDataLotAdd
EXEC p_Item_SetOption @nItemKeyTo, @nOptionType, @nOptionDataLotResult, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 7
	RETURN
END

-- ¾ÆÀÌÅÛ »èÁ¦
EXEC p_Item_Delete @nItemKeyFrom, @nRet OUTPUT
IF @@ERROR <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 8
	RETURN
END

-- ¿Ï·á
COMMIT TRAN

-- end

0
GO
