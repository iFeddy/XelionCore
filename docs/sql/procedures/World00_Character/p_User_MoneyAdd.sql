CREATE PROCEDURE [dbo].[p_User_MoneyAdd]
/*
Add user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nUserNo int,
@nAddMoney bigint,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

IF @nAddMoney < 0
BEGIN
	SET @nRet = 0
	RETURN
END

DECLARE @nUserMoney bigint
EXEC p_User_MoneyGet @nUserNo, @nUserMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0
	RETURN
-- Add money
SET @nUserMoney = @nUserMoney + @nAddMoney
EXEC p_User_MoneySet @nUserNo, @nUserMoney, @nRet OUTPUT
-- End

0ᜁ病ę
GO
