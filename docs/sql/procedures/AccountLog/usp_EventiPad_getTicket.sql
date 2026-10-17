CREATE PROCEDURE dbo.usp_EventiPad_getTicket
	@nUserNo int,
	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -1
		RETURN
	END
		
	SELECT nUserNo, nTicket, sys_create_date AS [dCreateDate], sys_update_date AS [dLastUpdate]
	  FROM tEventiPad
	 WHERE nUserNo = @nUserNo
	
	SET NOCOUNT OFF
END
GO
