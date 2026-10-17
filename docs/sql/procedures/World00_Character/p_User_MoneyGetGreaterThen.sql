CREATE PROCEDURE [dbo].[p_User_MoneyGetGreaterThen]
	@nMoney		bigint
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nUserNo, nUserMoney FROM tUserMoney WHERE nUserMoney >=  @nMoney
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0褁帍ě
GO
