CREATE PROCEDURE [dbo].[usp_ChargeItem_GetList]
	@USER_NO INT
AS

BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	EXEC Account.dbo.usp_Charge_ItemSelectList @USER_NO

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1

END
GO
