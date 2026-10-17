CREATE PROCEDURE [dbo].[usp_Item_setCategory]
	@category NVARCHAR(10),
	@name NVARCHAR(20),
	@parent NVARCHAR(10),
	@status BIT,
	@RET INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	IF EXISTS(SELECT category FROM tItemCategory WITH(NOLOCK) WHERE category = @parent)
		BEGIN
			IF EXISTS(SELECT category FROM tItemCategory WITH(NOLOCK) WHERE category = @category)
				BEGIN
					UPDATE tItemCategory 
						SET name = @name,
							parent = @parent,
							status = @status
					WHERE category = @category
					IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
						SET @RET = -2
					ELSE
						SET @RET = 1
				END
			ELSE
				BEGIN
					INSERT INTO tItemCategory (category, name, parent)
					VALUES (@category, @name, @parent)
					IF @@ERROR <> 0
						SET @RET = -1
					ELSE
						SET @RET = 1
				END
		END
	ELSE
		SET @RET = -3

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1
END

0
GO
