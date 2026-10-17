CREATE PROCEDURE [dbo].[usp_Attendance_getCallDateList]
	@nCharNo    int
,	@dStartDate datetime
,	@dEndDate   datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT dAttendanceDate 
	  FROM tAttendance 
	 WHERE nCharNo = @nCharNo
	   AND dAttendanceDate >= @dStartDate 
	   AND dAttendanceDate <= @dEndDate 	
END

0ꘁ⋂Ġ
GO
