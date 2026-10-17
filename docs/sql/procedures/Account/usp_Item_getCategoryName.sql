CREATE PROCEDURE [dbo].[usp_Item_getCategoryName]
	@category NVARCHAR(10) = '',
	@CATEGORY_NAME NVARCHAR(10) OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CATEGORY_NAME = name FROM tItemCategory WITH(NOLOCK) WHERE category = @category

	SET NOCOUNT OFF
END

0
GO
