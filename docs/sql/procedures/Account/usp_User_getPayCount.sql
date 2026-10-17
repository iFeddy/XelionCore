CREATE PROCEDURE [dbo].[usp_User_getPayCount]
	@userNo INT = 0,
	@CNT    INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CNT = ISNULL(COUNT(orderNo), 0)
	  FROM VT_PAY WITH (NOLOCK)
	 WHERE userNo = @userNo

	SET NOCOUNT OFF
END

0
GO
