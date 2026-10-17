CREATE PROCEDURE [dbo].[usp_Item_setPutOnBelonged]
/*
Get item storage
2009.9.23 By kodong
'	output Parameter
'		@nRet	 : check SQL execute(valid/invalid)
'			= 0	 : works ok
'			< 0	 : error
'				-2592 : db error(insert)
'				-2593 : db error(update)
*/
@nItemKey     bigint,
@nStorageType tinyint,
@nStorage     smallint,
@nOptionType  smallint,
@nOptionData  bigint,
@nRet         int OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	BEGIN TRAN
	
	UPDATE tItem 
	SET	nStorageType = @nStorageType, nStorage = @nStorage
	WHERE nItemKey = @nItemKey
	IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
		SET @nRet = -2593
		ROLLBACK TRAN
		RETURN
	END

	IF EXISTS(SELECT nItemKey FROM tItemOptions WHERE nItemKey = @nItemKey AND nOptionType = @nOptionType) BEGIN
		UPDATE tItemOptions
		SET nOptionType = @nOptionType,
			nOptionData = @nOptionData
		WHERE nItemKey = @nItemKey
		AND   nOptionType = @nOptionType
		IF @@ERROR <> 0 OR @@ROWCOUNT = 0 BEGIN
			SET @nRet = -2593
			ROLLBACK TRAN
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tItemOptions (nItemKey, nOptionType, nOptionData)
		VALUES (@nItemKey, @nOptionType, @nOptionData)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2592
			ROLLBACK TRAN
			RETURN
		END
	END

	COMMIT TRAN

	SET NOCOUNT OFF
END

0
GO
