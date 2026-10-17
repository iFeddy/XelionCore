CREATE PROCEDURE [dbo].[usp_Stat_sellitem]
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @start AS NVARCHAR(10), @end AS NVARCHAR(10)

	SET @start = CONVERT(NVARCHAR(10), DATEADD(dd, -1, GETDATE()), 120)
	SET @end = CONVERT(NVARCHAR(10), GETDATE(), 120)

	INSERT INTO tbSellitem (baseDate, goodsNo, amount)
	SELECT CONVERT(NVARCHAR(10), registerDate, 120) AS baseDate
	      , goodsNo AS goodsNo
	      , SUM(amount) AS amount
	  FROM tChargeItem WITH(NOLOCK)
	 WHERE registerDate >= @start AND registerDate < @end
	 GROUP BY CONVERT(NVARCHAR(10), registerDate, 120),goodsNo
	 ORDER BY CONVERT(NVARCHAR(10), registerDate, 120),goodsNo

END

0
GO
