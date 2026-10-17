CREATE PROCEDURE [dbo].[usp_Item_Add_EnchantGem] 
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nGemItemKey BIGINT,
	@nGemSocketIndex TINYINT,
	@nGemID INT,
	@nRestGemCount TINYINT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	
    BEGIN TRAN
    
    INSERT tItemEnchant (nItemKey, nSlot, nGemItemID, nRestCount)
		VALUES (@nEquipmentItemKey, @nGemSocketIndex, @nGemID, @nRestGemCount)

	IF @@ERROR <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -101
		RETURN
	END

	EXEC p_Item_Delete @nGemItemKey, @nRet OUTPUT
	
	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -202
		RETURN
	END

    COMMIT TRAN 
END

0
GO
