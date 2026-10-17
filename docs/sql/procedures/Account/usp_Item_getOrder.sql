CREATE PROCEDURE [dbo].[usp_Item_getOrder]
	@orderNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderDate, orderType, userNo, goodsNo, orderPrice, cashType, isPay, userIP
	  FROM tOrder WITH (NOLOCK)
	WHERE orderNo = @orderNo

	SET NOCOUNT OFF
END

0
GO
