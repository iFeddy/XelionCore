CREATE PROCEDURE [dbo].[usp_Cash_getCash]
	@userNo   int
,	@cashType tinyint = 1
,	@CASH	  int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @CASH = 0
	SELECT @CASH = ISNULL(cash, 0) FROM tCash 
	 WHERE userNo = @userNo AND cashType = @cashType AND status = 0

	SET NOCOUNT OFF
END

0
GO
