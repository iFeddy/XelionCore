CREATE PROCEDURE [dbo].[usp_Pay_chkLimit]
	@userNo  int
,	@chkType int
,	@amount  int
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @TOTAL_AMOUNT INT

	IF EXISTS(SELECT orderNo FROM VT_PAY WHERE userNo = @userNo)
		BEGIN
			IF @chkType = 1
				BEGIN
					SELECT @TOTAL_AMOUNT = ISNULL(SUM(cash), 0) 
					  FROM VT_PAY 
					 WHERE userNo = @userNo AND resultDate > CONVERT(NVARCHAR(10), GETDATE(), 120) AND resultDate <= GETDATE()
					IF @TOTAL_AMOUNT >= @amount
						RETURN (1)
					ELSE
						RETURN (0)
				END
			ELSE
				RETURN(0)
		END
	ELSE
		RETURN (0)

	SET NOCOUNT OFF
END

0
GO
