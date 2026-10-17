CREATE PROCEDURE [dbo].[usp_Attendance_setCall]
	@nCharNo        int
,	@dAttendaceDate datetime
,	@nRet       int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tAttendance WHERE nCharNo = @nCharNo AND dAttendanceDate = @dAttendaceDate) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tAttendance (nCharNo, dAttendanceDate)
		VALUES (@nCharNo, @dAttendaceDate)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -9999
		END
	END
END

0儁ＯĢ
GO
