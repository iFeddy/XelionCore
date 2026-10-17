CREATE PROCEDURE [dbo].[usp_Charge_ItemSelectList]
	@userNo INT
AS
BEGIN
	SET NOCOUNT ON

	SELECT rowNo, goodsNo, amount, registerDate
	FROM tChargeItem WITH (NOLOCK)
	WHERE userNo = @userNo AND isDraw = 0
    ORDER BY rowNo ASC

	SET NOCOUNT OFF
END

0
GO
