CREATE PROCEDURE [dbo].[usp_Attendance_setReward]
	@nCharNo     int
,	@dEventDate  datetime	
,	@nRewardID   int
,	@nRet        int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tAttendanceReward WHERE nCharNo = @nCharNo AND dEventDate = CAST(@dEventDate as datetime) AND nRewardID = @nRewardID) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tAttendanceReward (nCharNo, dEventDate, dRewardDate, nRewardID)
		VALUES (@nCharNo, @dEventDate, GETDATE(), @nRewardID)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -9999
		END
	END
END

0嵧Ĥ
GO
