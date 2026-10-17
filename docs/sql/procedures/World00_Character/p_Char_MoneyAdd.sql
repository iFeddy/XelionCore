CREATE PROCEDURE [dbo].[p_Char_MoneyAdd]
/*
Add money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nAddMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
DECLARE @nCharMoney bigint
EXEC p_Char_MoneyGet @nCharNo, @nCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add money
SET @nCharMoney = @nCharMoney + @nAddMoney
EXEC p_Char_MoneySet @nCharNo, @nCharMoney, @nRet OUTPUT
-- End

0ခ㲢ũ
GO
