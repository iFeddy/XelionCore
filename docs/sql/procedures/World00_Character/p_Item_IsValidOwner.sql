CREATE PROCEDURE [dbo].[p_Item_IsValidOwner]
/*
Is valid item owner ?
2004.8 By CJC
 input:
	nItemKey	 = Item Key
	nOwner	     = Item Owner
 output:
	0 = OK
	1 = Not found item
	2 = Missmatch item owner
*/
@nItemKey bigint,
@nOwner int,
@nRet tinyint OUTPUT 
AS
SET NOCOUNT ON
DECLARE @nDBOwner int 
SET @nDBOwner = 0
SET @nRet = 0
SELECT @nDBOwner = nOwner FROM tItem WITH (READUNCOMMITTED) WHERE nItemKey = @nItemKey
-- Check not found item
IF @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 1
	RETURN
END
-- Check missmatch item owner
IF @nDBOwner <> @nOwner
BEGIN
	SET @nRet = 2
	RETURN
END
-- end

0昁ᦴţ
GO
