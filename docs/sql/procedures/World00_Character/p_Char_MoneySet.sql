CREATE PROCEDURE [dbo].[p_Char_MoneySet]
/*
Set money
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nCharNo int,
@nSetMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0
UPDATE tCharacter SET nMoney = @nSetMoney WHERE nCharNo = @nCharNo
IF @@ERROR = 0 AND @@ROWCOUNT = 1
	SET @nRet = 1
-- end
RETURN

0쨁Ŭ
GO
