CREATE PROCEDURE [dbo].[p_Item_SetOption]
/*
Item Set Option
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nOptionType smallint,
@nOptionData bigint,
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemOptions SET	nOptionData = @nOptionData
WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemOptions (  nItemKey,  nOptionType,   nOptionData )
	VALUES				( @nItemKey, @nOptionType,  @nOptionData )
END
SET @nRet = @@ERROR
-- end

RETURN

0朁껺ū
GO
