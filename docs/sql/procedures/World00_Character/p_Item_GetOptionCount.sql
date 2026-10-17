CREATE PROCEDURE [dbo].[p_Item_GetOptionCount]
/*
Item Get Option Count
2004.8 By CJC
 input:
	see SQL
 output:
	@nOptionCount = Item option count
*/
@nItemKey bigint,
@nOptionCount int OUTPUT
AS
SET NOCOUNT ON
SELECT @nOptionCount = COUNT(nItemKey) FROM tItemOptions WHERE nItemKey = @nItemKey
-- end

0鸁碒ś
GO
