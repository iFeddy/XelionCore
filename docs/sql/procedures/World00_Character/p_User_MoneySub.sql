CREATE PROCEDURE [dbo].[p_User_MoneySub]
/*
Sub user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error, or Minus value
*/
@nUserNo int,
@nSubMoney bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

IF @nSubMoney < 0
BEGIN
	SET @nRet = 0
	RETURN
END

DECLARE @nUserMoney bigint
EXEC p_User_MoneyGet @nUserNo, @nUserMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Sub money
SET @nUserMoney = @nUserMoney - @nSubMoney
-- Check minus
IF @nUserMoney < 0
BEGIN
	SET @nRet = 0
	RETURN	
END
EXEC p_User_MoneySet @nUserNo, @nUserMoney, @nRet OUTPUT
-- End

0企끪ĝ
GO
