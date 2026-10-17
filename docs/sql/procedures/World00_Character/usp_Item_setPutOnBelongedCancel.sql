CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelongedCancel]
	@nCharNo               int
,	@nCancelItemKey        bigint
,	@nCancelOptionType     smallint
,	@nCancelOptionData     bigint
,	@nTargetItemKey        bigint
,	@nTargetOptionType     smallint
,	@nTargetOptionData     bigint
,	@nRet                  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cancelOptionData bigint
	DECLARE @targetOptionData bigint
	
	SET @nRet = 0
	SET @cancelOptionData = 0
	SET @targetOptionData = 0
	
	-- invalid nCharNo	
	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2511
		RETURN
	END
		
	IF NOT EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nCancelItemKey AND nOwner = @nCharNo) BEGIN
		SET @nRet = -2512
		RETURN
	END

	IF NOT EXISTS(SELECT nItemKey FROM tItem WHERE nItemKey = @nTargetItemKey AND nOwner = @nCharNo) BEGIN
		SET @nRet = -2513
		RETURN
	END

	BEGIN TRAN

	-- ÇØÁŠ À¯·áŸÆÀÌÅÛ Â÷°š ¶ÇŽÂ »èÁŠ	
	SELECT @cancelOptionData = nOptionData
	  FROM tItemOptions WITH (UPDLOCK)
	 WHERE nItemKey    = @nCancelItemKey 
	   AND nOptionType = @nCancelOptionType 
	IF @nCancelOptionData <> @cancelOptionData BEGIN
		SET @nRet = -2521
		ROLLBACK TRAN
		RETURN
	END
	IF @cancelOptionData > 1 BEGIN
		SET @cancelOptionData = @cancelOptionData - 1
		EXEC p_Item_SetOption @nCancelItemKey, @nCancelOptionType, @cancelOptionData, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2522
			ROLLBACK TRAN
			RETURN						
		END					
	END
	ELSE BEGIN
		EXEC p_Item_Delete @nCancelItemKey, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2523
			ROLLBACK TRAN
			RETURN						
		END					
	END
	
	-- Âø¿ëœÃ ±ÍŒÓ ¿ÉŒÇº¯°æ(1 -> 0)
	SELECT @targetOptionData = nOptionData
	  FROM tItemOptions WITH (UPDLOCK)
	 WHERE nItemKey    = @nTargetItemKey 
	   AND nOptionType = @nTargetOptionType 
	IF @targetOptionData <> 1 BEGIN
		SET @nRet = -2524
		ROLLBACK TRAN
		RETURN
	END
	ELSE BEGIN
		EXEC p_Item_SetOption @nTargetItemKey, @nTargetOptionType, @nTargetOptionData, @nRet output
		IF @nRet <> 0 BEGIN
			SET @nRet = -2525
			ROLLBACK TRAN
			RETURN						
		END			
	END
	
	COMMIT TRAN
END

0ᤁ쒊
GO
