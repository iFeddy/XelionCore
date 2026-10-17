CREATE PROCEDURE [dbo].[usp_Item_getGroup]
	@groupNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT category, name, imageName, itemType, description 
	  FROM VT_ITEM_GROUP WITH (NOLOCK) 
	 WHERE groupNo = @groupNo

	SET NOCOUNT OFF
END

0
GO
