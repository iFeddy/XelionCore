CREATE PROCEDURE [dbo].[usp_InGameGuide_setGuideNo]
	@nCharNo  int
,	@nGuideNo int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tInGameGuide WITH(NOLOCK) WHERE nCharNo = @nCharNo AND nGuideNo = @nGuideNo) BEGIN
		SET @nRet = -1111
	END
	ELSE BEGIN
		INSERT INTO tInGameGuide (nCharNo, nGuideNo)
		VALUES (@nCharNo, @nGuideNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2222
		END
	END
END

0洁勗
GO
