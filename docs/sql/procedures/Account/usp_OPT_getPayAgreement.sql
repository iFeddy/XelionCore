CREATE PROCEDURE [dbo].[usp_OPT_getPayAgreement]
	@userNo  int
,	@getType int = 0
AS
BEGIN
	SET NOCOUNT ON

	IF @getType = 0
		SELECT rowNo, userNo, name, juminNo, phoneNo, status, issueDate, vaildDate
          FROM tPayAgreement WITH(NOLOCK)
         WHERE userNo = @userNo
		 ORDER BY rowNo DESC
	ELSE
		SELECT TOP 1 rowNo, userNo, name, juminNo, phoneNo, status, issueDate, vaildDate
          FROM tPayAgreement WITH(NOLOCK)
         WHERE userNo = @userNo
		 ORDER BY rowNo DESC

	SET NOCOUNT OFF
END

0
GO
