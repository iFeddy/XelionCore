CREATE PROCEDURE [dbo].[usp_SpamBlock_setTime]
	@nCharNo int
,	@nBlockTime int
,	@nStandbyTime int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tSpamBlockTime WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tSpamBlockTime
		   SET sys_update_date = GETDATE()
		     , nBlockTime = @nBlockTime
		     , nStandbyTime = @nStandbyTime
		 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
		BEGIN
			SET @nRet = -2893
		END
	END
	ELSE BEGIN
		INSERT INTO tSpamBlockTime (nCharNo, nBlockTime, nStandbyTime)
		VALUES (@nCharNo, @nBlockTime, @nStandbyTime)
		IF @@ERROR <> 0
		BEGIN
			SET @nRet = -2892
		END
	END

	SET NOCOUNT ON
END

0谁䈘)
GO
