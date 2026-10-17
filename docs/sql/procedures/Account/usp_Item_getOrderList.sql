CREATE PROCEDURE [dbo].[usp_Item_getOrderList]
	@userNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderDate, orderType, T1.goodsNo, name, imageName, orderPrice
	  FROM tOrder AS T1, tItem AS T2
	WHERE userNo = @userNo AND isPay = 1 AND T1.goodsNo = T2.goodsNo

	SET NOCOUNT OFF
END

0
GO
