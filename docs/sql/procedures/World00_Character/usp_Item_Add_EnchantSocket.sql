CREATE PROCEDURE [dbo].[usp_Item_Add_EnchantSocket]
	@nCharNo INT,
	@nEquipmentItemKey BIGINT,
	@nDrillItemKey BIGINT,
	@nAllSocketCount TINYINT OUTPUT,
	@nDrillLotType SMALLINT,
	@nDrillLotValue BIGINT OUTPUT,
	@nRet INT = 0 OUTPUT
AS
BEGIN
	-- SET NOCOUNT ON added to prevent extra result sets from
	-- interfering with SELECT statements.
	SET NOCOUNT ON
	
	SET @nRet = 0

    BEGIN TRAN

	EXEC p_Item_Delete @nDrillItemKey, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -101
		RETURN
	END
	
	SET @nDrillLotValue = 0
	
	DECLARE @nCreatedSocketCount TINYINT
	
	EXEC p_Item_GetOption @nEquipmentItemKey, 523, @nAllSocketCount OUTPUT, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -201
		RETURN
	END

	EXEC p_Item_GetOption @nEquipmentItemKey, 524, @nCreatedSocketCount OUTPUT, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet = 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -202
		RETURN
	END
	
	SET @nAllSocketCount = @nAllSocketCount + 1
	SET @nCreatedSocketCount = @nCreatedSocketCount + 1
	
	EXEC p_Item_SetOption @nEquipmentItemKey, 523, @nAllSocketCount, @nRet OUTPUT

	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -301
		RETURN
	END
	
	EXEC p_Item_SetOption @nEquipmentItemKey, 524, @nCreatedSocketCount, @nRet OUTPUT
	
	IF @@ERROR <> 0 OR @nRet <> 0
	BEGIN
		ROLLBACK TRAN
		SET @nRet = -302
		RETURN
	END
	
    COMMIT TRAN  
END

0
GO
