CREATE PROCEDURE [dbo].[p_Item_GetOptionAll]
/*
Item Get Option All
2004.8 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/
@nItemKey bigint
AS
SET NOCOUNT ON
SELECT nOptionType, nOptionData FROM tItemOptions WHERE nItemKey = @nItemKey
-- end

0攁葮Ś
GO
