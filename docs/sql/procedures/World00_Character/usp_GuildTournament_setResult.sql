CREATE PROCEDURE [dbo].[usp_GuildTournament_setResult]
	@nGTNo    int
,	@nType    smallint
,	@nData    int
,	@nGuildNo int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nGTNo FROM tGuildTournament WHERE  nGTNo = @nGTNo AND nType = @nType) BEGIN
		UPDATE tGuildTournament 
		   SET nData = @nData 
		     , nGuildNo = @nGuildNo 
		 WHERE nGTNo = @nGTNo 
		   AND nType = @nType
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1493
		END		   
	END
	ELSE BEGIN
		INSERT tGuildTournament (nGTNo, nType, nData, nGuildNo)
		VALUES (@nGTNo, @nType, @nData, @nGuildNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1492
		END			
	END
END

0ď
GO
