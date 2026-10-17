CREATE PROCEDURE [dbo].[usp_Item_getItem]
	@goodsNo INT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TODAY AS NVARCHAR(10)
	SET @TODAY = CONVERT(NVARCHAR(10), GETDATE(), 120)

	SELECT goodsNo, billingNo, category, groupNo, name, imageName, itemType, unit, unitDescription, price, discount, sellPrice = 
		CASE
			WHEN discountStartDate <= @TODAY AND discountEndDate >= @TODAY THEN price - ((price * discount) / 100) 
			ELSE price 
		END, isSell, description 
	  FROM tItem WITH (NOLOCK)
	WHERE goodsNo = @goodsNo

	SET NOCOUNT OFF
END

0
GO
