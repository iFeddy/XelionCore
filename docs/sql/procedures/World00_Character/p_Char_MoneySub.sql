CREATE PROCEDURE [dbo].[p_Char_MoneySub]
/*
Sub money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nCharNo int,
@nSubMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharMoney bigint
EXEC p_Char_MoneyGet @nCharNo, @nCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub money
SET @nCharMoney = @nCharMoney - @nSubMoney
-- Check minus
IF @nCharMoney < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_Char_MoneySet @nCharNo, @nCharMoney, @nRet OUTPUT
-- End

0ሁ譧Ů
GO
