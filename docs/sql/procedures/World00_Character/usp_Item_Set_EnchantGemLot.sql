CREATE PROCEDURE [dbo].[usp_Item_Set_EnchantGemLot]
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nSlot TINYINT,
	@nRestGemCount TINYINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	UPDATE tItemEnchant
		SET nRestCount = @nRestGemCount
		FROM tItemEnchant
		WHERE nItemKey = @nEquipmentItemKey
			AND nSlot = @nSlot

	IF @@ROWCOUNT <> 1 OR @@ERROR <> 0
	BEGIN
		SET @nRet = -101
		RETURN
	END 
END

0
GO
