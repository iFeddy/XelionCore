CREATE PROCEDURE [dbo].[usp_Item_Get_EnchantAll]
	@nEquipmentItemKey BIGINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nSlot, nGemItemID, nRestCount FROM tItemEnchant WHERE nItemKey = @nEquipmentItemKey

	IF @@ERROR <> 0
	BEGIN
		SET @nRet = -101
		RETURN
	END 
END

0
GO
