CREATE PROCEDURE [dbo].[usp_OPT_getBuyList]
	@userNo   int
AS
BEGIN
	SET NOCOUNT ON

	SELECT orderNo, orderDate, T1.goodsNo, name, unit, unitDescription, toUserNo, 
		   orderType = CASE orderType WHEN 0 THEN 'Buy' ELSE 'Gift' END, orderPrice,
		   cashType  = CASE cashType  WHEN 1 THEN 'Game' ELSE 'Event' END
	  FROM tOrder AS T1, tItem AS T2
	 WHERE userNo = @userNo AND isPay = 1 AND T1.goodsNo = T2.goodsNo
	 ORDER BY orderNo DESC

	SET NOCOUNT OFF
END

0
GO
