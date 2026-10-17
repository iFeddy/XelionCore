CREATE PROCEDURE [dbo].[p_Item_GetStorageEmpty]
	@owner       int = 0
,	@storageType tinyint = 0
,	@ret         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @err AS bit
	DECLARE @emptyStorage AS smallint
	SET @err = 0

	IF @owner = 0
		SET @err = 1
	IF @storageType = 0
		SET @err = 1

	IF @err <> 1
		BEGIN	
			IF EXISTS(SELECT nItemKey FROM tItem WHERE nOwner = @owner AND nStorageType = @storageType AND nStorage = 0)
				BEGIN
					SELECT @emptyStorage = ISNULL(MIN(nStorage + 1), 0)
					  FROM tItem 
					 WHERE nOwner = @owner 
					   AND nStorageType = @storageType 
					   AND nStorage + 1 NOT IN (
							SELECT nStorage
							  FROM tItem 
							 WHERE nOwner = @owner AND nStorageType = @storageType)
				END
			ELSE
				BEGIN
					SET @emptyStorage = 0
				END

			SET @ret = @emptyStorage
			PRINT @ret
		END
	ELSE
		SET @ret = -1

	SET NOCOUNT OFF
END

0Ёꨅ^
GO
