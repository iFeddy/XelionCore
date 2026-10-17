CREATE PROCEDURE [dbo].[p_Char_MoneyGetGreaterThen]
	@nMoney		bigint
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nUserNo, nCharNo, sID, nMoney, nLevel, nSlotNo, bDeleted, nAdminLevel 
	FROM tCharacter WHERE nMoney >=  @nMoney
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0鄁ū
GO
