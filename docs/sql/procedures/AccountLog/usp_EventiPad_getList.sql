CREATE PROCEDURE dbo.usp_EventiPad_getList
AS
BEGIN
	SET NOCOUNT ON

	SELECT nUserNo, nTicket, sys_create_date AS [dCreateDate], sys_update_date AS [dLastUpdate]
	  FROM tEventiPad
	 ORDER BY nTicket DESC
	
	SET NOCOUNT OFF
END
0騁ａD
GO
