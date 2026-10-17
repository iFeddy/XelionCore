CREATE PROCEDURE [dbo].[usp_User_getBuyCount]
	@userNo INT = 0,
	@CNT    INT = 0 OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SELECT @CNT = ISNULL(COUNT(orderNo), 0)
	  FROM tOrder WITH (NOLOCK)
	 WHERE userNo = @userNo AND isPay = 1

	SET NOCOUNT OFF
END

0
GO
