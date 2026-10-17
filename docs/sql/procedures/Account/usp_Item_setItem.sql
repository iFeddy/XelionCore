CREATE PROCEDURE [dbo].[usp_Item_setItem]
	@goodsNo INT,
	@category NVARCHAR(10),
	@groupNo INT,
	@name NVARCHAR(100),
	@imageName NVARCHAR(100),
	@itemType NVARCHAR(20),
	@unit INT,
	@price INT,
	@discount INT,
	@isRecommand BIT,
	@isPopular BIT,
	@isNew BIT,
	@isSell BIT,
	@sellStartDate NVARCHAR(10) = '',
	@sellEndDate NVARCHAR(10) = '',
	@discountStartDate NVARCHAR(10) = '',
	@discountEndDate NVARCHAR(10) = '',
	@RET INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT goodsNo FROM tItem WHERE goodsNo = @goodsNo)
		BEGIN
			BEGIN TRAN

			UPDATE tItem SET
				category = @category,
				groupNo = @groupNo,
				name = @name,
				imageName = @imageName,
				itemType = @itemType,
				unit = @unit,
				price = @price,
				discount = @discount,
				isRecommand = @isRecommand,
				isPopular = @isPopular,
				isNew = @isNew,
				isSell = @isSell,
				sellStartDate = @sellStartDate,
				sellEndDate = @sellEndDate,
				discountStartDate = @discountStartDate,
				discountEndDate = @discountEndDate
			WHERE goodsNo = @goodsNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
				BEGIN
					SET @RET = -1
					ROLLBACK TRAN
				END
			ELSE
				BEGIN
					SET @RET = 1
					COMMIT TRAN
				END
		END
	ELSE
		SET @RET = -2

	SET LOCK_TIMEOUT -1
	SET NOCOUNT OFF
END
GO
