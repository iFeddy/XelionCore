CREATE PROCEDURE [dbo].[usp_Pay_getAgreement]
	@userNo int
AS
BEGIN
	SET NOCOUNT ON

	SELECT TOP 1 status
	  FROM tPayAgreement WITH(NOLOCK) 
	 WHERE userNo = @userNo AND vaildDate >= GETDATE()
	 ORDER BY rowNo DESC

	SET NOCOUNT OFF
END

0
GO
