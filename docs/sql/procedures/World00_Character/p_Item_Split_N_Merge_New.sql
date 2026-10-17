CREATE PROCEDURE [dbo].[p_Item_Split_N_Merge_New]
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
                     5 = From lot read failed
                     6 = From lot less more nSplitLot
                     7 = To lot read failed
                     8 = From lot set failed
                     9 = To lot set failed
*/

@nItemKeyFrom bigint,
@nOwnerFrom int,
@nOptionTypeFrom smallint,
@nOptionDataFrom bigint OUTPUT,

@nItemKeyTo bigint,
@nOwnerTo int,
@nOptionTypeTo smallint,
@nOptionDataTo bigint OUTPUT,

@nSplitLot bigint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- ¼ö·® ÁÙ¾îµé ¾ÆÀÌÅÛ(From) ¼ÒÀ¯ÀÚ È®ÀÎ
DECLARE @nTinyintRet tinyint
EXEC p_Item_IsValidOwner @nItemKeyFrom, @nOwnerFrom, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet
	RETURN
END

-- ¼ö·® ´Ã¾î³¯ ¾ÆÀÌÅÛ(To) ¼ÒÀ¯ÀÚ È®ÀÎ
EXEC p_Item_IsValidOwner @nItemKeyTo, @nOwnerTo, @nTinyintRet OUTPUT
IF @nTinyintRet <> 0
BEGIN
	SET @nRet = @nTinyintRet + 2
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ù ¾ò±â
EXEC p_Item_GetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 5
	RETURN
END

-- »©Áú ¾ÆÀÌÅÛ ¶ùÀÌ ÃæºÐÇÑÁö È®ÀÎ
IF @nOptionDataFrom < @nSplitLot
BEGIN
	SET @nRet = 6
	RETURN
END

-- ´õÇØÁú ¾ÆÀÌÅÛ ¶ù ¾ò±â
DECLARE @nOptionDataLotAdd bigint
EXEC p_Item_GetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 7
	RETURN
END

-- °á°ú ¶ù °è»ê
SET @nOptionDataFrom = @nOptionDataFrom - @nSplitLot
SET @nOptionDataTo = @nOptionDataTo + @nSplitLot

BEGIN TRAN
-- ¼ö·® ÁÙ¾îµé ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
EXEC p_Item_SetOption @nItemKeyFrom, @nOptionTypeFrom, @nOptionDataFrom, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 8
	RETURN
END

-- ¼ö·® ´Ã¾î³¯ ¾ÆÀÌÅÛ ¶ù º¯°æÀúÀåÇÁ·Î½ÃÁ® È£Ãâ
EXEC p_Item_SetOption @nItemKeyTo, @nOptionTypeTo, @nOptionDataTo, @nRet OUTPUT
IF @nRet <> 0
BEGIN
	ROLLBACK TRAN
	SET @nRet = 9
	RETURN
END

-- ¿Ï·á
COMMIT TRAN
-- end

0
GO
