CREATE PROCEDURE [dbo].[usp_Attendance_getRewardList]
	@nCharNo    int
,	@dDate      datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	SELECT dEventDate, nRewardID  
	  FROM tAttendanceReward  
	 WHERE nCharNo = @nCharNo
	   AND dEventDate > @dDate
END

0氁生Ģ
GO
