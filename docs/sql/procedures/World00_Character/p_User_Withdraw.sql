CREATE PROCEDURE [dbo].[p_User_Withdraw]
/*
Set User withdraw money
2006.5 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK
	        1 = Miss match UserNo and User's CharNo
	        2 = User Money Sub Failed(Not enough money)
	        3 = Withdraw money faield
	        4 = Get result money faield
*/
@nUserNo int,
@nCharNo int,
@nMoney bigint,
@nResultUserMoney bigint OUTPUT,
@nRet tinyint OUTPUT
AS

SET NOCOUNT ON

SET @nRet = 0

-- Check UserNo and User's CharNo
DECLARE @nDBUserNo int
SET @nDBUserNo = 0

SELECT @nDBUserNo = nUserNo FROM tCharacter WHERE nCharNo = @nCharNo
IF @nDBUserNo <> @nUserNo
BEGIN
	SET @nRet = 1
	RETURN
END

-- Deposit money
BEGIN TRAN

EXEC p_User_MoneySub @nUserNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 2
	ROLLBACK TRAN
	RETURN
END

EXEC p_Char_MoneyAdd @nCharNo, @nMoney, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 3
	ROLLBACK TRAN
	RETURN
END

EXEC p_User_MoneyGet @nUserNo, @nResultUserMoney OUTPUT, @nRet OUTPUT
IF @nRet <> 1
BEGIN
	SET @nRet = 4
	ROLLBACK TRAN
	RETURN
END

-- Done
COMMIT TRAN

SET @nRet = 0
RETURN

0蠁꒎Ğ
GO
