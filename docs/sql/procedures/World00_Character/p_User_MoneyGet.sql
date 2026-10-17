CREATE PROCEDURE [dbo].[p_User_MoneyGet]
/*
Get user money
2005.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK, Money return
	        0 = Error
*/
@nUserNo int,
@nUserMoney bigint OUTPUT,
@nRet tinyint OUTPUT

AS

SET NOCOUNT ON

SET @nRet = 0

-- 돈 얻기
IF EXISTS(SELECT nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo )
BEGIN
	SELECT @nUserMoney = nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo
	IF @@ROWCOUNT = 1
		SET @nRet = 1
	RETURN 
END

-- First access
SET @nUserMoney = 0
INSERT INTO tUserMoney VALUES (@nUserNo, @nUserMoney)
IF @@ROWCOUNT = 1 AND @@ERROR = 0
	SET @nRet = 1
	
RETURN 

-- End

0倁槩Ě
GO
