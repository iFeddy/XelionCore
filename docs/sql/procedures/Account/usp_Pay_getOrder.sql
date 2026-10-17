CREATE PROCEDURE [dbo].[usp_Pay_getOrder]
	@orderNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT T1.orderNo, T1.userNo, T1.payNo, T1.cash, T1.userIP, T1.registerDate,
		   T2.result, T2.resultDate, T2.pgTID
	  FROM tPay T1 LEFT JOIN tPayResult T2 ON T1.orderNo = T2.orderNo
	 WHERE T1.orderNo = @orderNo

	SET NOCOUNT OFF
END

0
GO
