CREATE PROCEDURE [p_Report_Add] 
@nReporter int,
@sReporterCharID nvarchar(16),
@sHarmer nvarchar(16),
@nWorldNo tinyint,
@sReportType nvarchar(16),
@sReason nvarchar(512),
@sChatLog text,
@nReportRegNo int OUTPUT 
AS
SET NOCOUNT ON
DECLARE	@ret	int
SET 		@ret = 0
--1시간에 1회만 가능합니다.
IF EXISTS (SELECT nReporter FROM tReport(NOLOCK) 
		WHERE 	nReporter = @nReporter AND 
			nWorldNo=@nWorldNo AND
			dReportDate > dateadd(hour, -1, getdate()) )
BEGIN
	SET @nReportRegNo = -1
	RETURN @ret
END
--동일 캐릭터에 대해서는 1일 1회만 가능합니다.
IF EXISTS (SELECT nReporter FROM tReport(NOLOCK) 
		WHERE 	nReporter = @nReporter AND 
			nWorldNo=@nWorldNo AND
			dReportDate > dateadd(day, -1, getdate()) AND
			sHarmer = @sHarmer)
BEGIN
	SET @nReportRegNo = -2
	RETURN @ret
END
INSERT tReport (sReportType,nReporter,sReporterCharID,sHarmer,nWorldNo,sReason,sChatLog)
VALUES (@sReportType,@nReporter, @sReporterCharID, @sHarmer, @nWorldNo,@sReason,@sChatLog )
SET @ret = @@error
IF @ret <> 0 
	SET @nReportRegNo = 0
ELSE
	SET @nReportRegNo = @@IDENTITY
RETURN @ret

0촁珁Ē
GO
