CREATE PROCEDURE [dbo].[usp_Item_Del_EnchantGem] 
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nRemoverItemKey BIGINT,
	@nResetSocketIndex TINYINT,
	@bNeedItem TINYINT,
	@nRemoverLotType SMALLINT,
	@nRemoverLotValue BIGINT = 0 OUTPUT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0
	SET @nRemoverLotValue = 0
	
    BEGIN TRAN

	IF @bNeedItem = 1
	BEGIN
		EXEC p_Item_Delete @nRemoverItemKey, @nRet OUTPUT

		IF @@ERROR <> 0 OR @nRet <> 0
		BEGIN
			ROLLBACK TRAN
			SET @nRet = -101
			RETURN
		END
	END
	
    DELETE FROM tItemEnchant
		WHERE nItemKey = @nEquipmentItemKey
			AND nSlot = @nResetSocketIndex

	IF @@ERROR <> 0 OR @@ROWCOUNT = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -201
		RETURN
	END

    COMMIT TRAN
END

0
GO
