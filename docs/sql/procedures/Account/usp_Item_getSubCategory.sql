CREATE PROCEDURE [dbo].[usp_Item_getSubCategory]
	@parent NVARCHAR(10)
AS
BEGIN
	SET NOCOUNT ON

	SELECT * FROM tItemCategory WITH(NOLOCK) WHERE parent = @parent

	SET NOCOUNT OFF
END

0
GO
