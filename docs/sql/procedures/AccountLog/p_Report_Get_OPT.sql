CREATE PROCEDURE [p_Report_Get_OPT]
	@QueryExpr nvarchar(256)
AS
BEGIN
	SET NOCOUNT ON
	DECLARE @SQLCMD nvarchar(512)
	SET @SQLCMD = N'SELECT nRegNumber, sReportType, nReporter, sReporterCharID, sHarmer, nWorldNo, sReason, sChatLog, dReportDate, nResult, sRemark 
		FROM tReport ' + @QueryExpr
	EXEC sp_executesql @SQLCMD
	SET NOCOUNT OFF
END 

0威ć
GO
