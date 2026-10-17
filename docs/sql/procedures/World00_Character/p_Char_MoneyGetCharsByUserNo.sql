CREATE PROCEDURE [dbo].[p_Char_MoneyGetCharsByUserNo]
	@nUserNo		int
AS

BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	SELECT nCharNo, sID, nMoney, nLevel, nSlotNo, bDeleted, nAdminLevel 
	FROM tCharacter WHERE nUserNo =  @nUserNo
	
	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF	
END

0⸁뫖Ū
GO
