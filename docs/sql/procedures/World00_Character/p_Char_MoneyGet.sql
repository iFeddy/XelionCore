CREATE PROCEDURE [dbo].[p_Char_MoneyGet]
/*
Get money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Money return
	        0 = Error
*/
@nCharNo int,
@nCharMoney bigint OUTPUT,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
-- 돈 얻기
SELECT @nCharMoney = nMoney FROM tCharacter WHERE nCharNo = @nCharNo
IF @@ROWCOUNT = 1
	SET @nRet = 1
	
-- End
RETURN

0괁፳Ũ
GO
