CREATE PROCEDURE [p_Report_Get] 
@nReporter int,
@nWorldNo tinyint
AS
SET NOCOUNT ON
SELECT 	nRegNumber,nResult,sRemark  FROM tReport
WHERE 	nReporter = @nReporter AND nWorldNo=@nWorldNo AND
	dReportDate > dateadd(hour, -1, getdate()) 
RETURN @@error

0㡶
GO
