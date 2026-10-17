CREATE PROCEDURE [dbo].[usp_GuildAcademy_setMemberLevelup]
	@nNo       int
,	@nPoint    int
,	@nBuffTime int
,	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	DECLARE @BUFF_END_TIME AS DATETIME

	IF EXISTS(SELECT nNo FROM tGuildAcademy WHERE nNo = @nNo) BEGIN
		BEGIN TRAN

		SELECT @BUFF_END_TIME = ISNULL(dBuffEndTime, GETDATE()) FROM tGuildAcademy WHERE nNo = @nNo
		IF @BUFF_END_TIME > GETDATE() BEGIN
			SET @BUFF_END_TIME = DATEADD ( minute, @nBuffTime, @BUFF_END_TIME ) 
		END
		ELSE BEGIN
			SET @BUFF_END_TIME = DATEADD ( minute, @nBuffTime, GETDATE() ) 
		END

		UPDATE tGuildAcademy SET nPoint = nPoint + @nPoint, dBuffEndTime = @BUFF_END_TIME
		 WHERE nNo = @nNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1593
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN			
			COMMIT TRAN
		END
	END
	ELSE BEGIN
		SET @nRet = -1521
	END	

	SET NOCOUNT OFF
END

0ᘁҸ
GO
