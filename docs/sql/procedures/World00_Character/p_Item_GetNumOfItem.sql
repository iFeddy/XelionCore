CREATE PROCEDURE [dbo].[p_Item_GetNumOfItem]
/*
Get number of item
2004.8 By CJC
 input:
	nCharNo	    = Item Owner
 output:
	nItemCount,	-1 = error
	             ? = Count of item;
*/
@nCharNo int,
@nNumOfItem int OUTPUT
AS
SET NOCOUNT ON
SET @nNumOfItem = -1
SELECT @nNumOfItem = COUNT(*) FROM tItem WHERE nOwner = @nCharNo
-- end

0䜁غř
GO
