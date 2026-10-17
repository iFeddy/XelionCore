CREATE PROCEDURE [dbo].[usp_Item_getGroupPrice]
	@groupNo INT
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TODAY AS NVARCHAR(10)
	SET @TODAY = CONVERT(NVARCHAR(10), GETDATE(), 120)

	SELECT goodsNo, unit, unitDescription, price, discount, sellPrice = 
		CASE
			WHEN discountStartDate <= @TODAY AND discountEndDate >= @TODAY THEN price - ((price * discount) / 100) 
			ELSE price 
		END
	  FROM VT_ITEM_SELL
	WHERE groupNo = @groupNo
	ORDER BY goodsNo

	SET NOCOUNT OFF
END
GO
