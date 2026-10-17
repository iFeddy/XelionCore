CREATE PROCEDURE [dbo].[usp_Pay_setOrder]
	@userNo int
,	@payNo  int
,	@cash   int
,	@userIP nvarchar(30)
,	@ORDER_NO int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tPay (userNo, payNo, cash, userIP)
	VALUES (@userNo, @payNo, @cash, @userIP)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		SET @ORDER_NO = 0
	ELSE
		SET @ORDER_NO = @@IDENTITY

	SET NOCOUNT OFF
END
GO
