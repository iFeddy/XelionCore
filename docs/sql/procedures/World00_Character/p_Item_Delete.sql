CREATE PROCEDURE [dbo].[p_Item_Delete]
/*
Item Delete
2004.8 By CJC
 input:
	nItemKey = To be delete.
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 아이템 삭제 옵션
DELETE FROM tItem WHERE nItemKey = @nItemKey
SET @nRet = @@ERROR
-- end

0۟Ő
GO
