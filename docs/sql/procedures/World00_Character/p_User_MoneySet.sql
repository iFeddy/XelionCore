CREATE PROCEDURE [dbo].[p_User_MoneySet]
/*
Set user money
2006.5 By CJC
 input:
	see SQL
 output:
	@nRet	1 = OK
	        0 = Error
*/
@nUserNo int,
@nSetMoney bigint,
@nRet tinyint OUTPUT
AS
SET NOCOUNT ON
SET @nRet = 0

IF EXISTS( SELECT nUserMoney FROM tUserMoney WHERE nUserNo = @nUserNo )
BEGIN
	UPDATE tUserMoney SET nUserMoney = @nSetMoney WHERE nUserNo = @nUserNo
	IF @@ERROR <> 0
		RETURN
END
ELSE
BEGIN
	INSERT INTO tUserMoney VALUES (@nUserNo, @nSetMoney)
	IF @@ERROR <> 0
		RETURN
END

-- DONE
SET @nRet = 1
RETURN

0숁刱Ĝ
GO
