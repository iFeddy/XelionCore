CREATE PROCEDURE [dbo].[usp_Collection_setReward]
	@nCharNo      int
,	@nReward      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	IF EXISTS(SELECT nCharNo FROM tCollectionReward WHERE nCharNo = @nCharNo AND nReward = @nReward) BEGIN
		SET @nRet = -2222
		RETURN
	END
	ELSE BEGIN
		INSERT INTO tCollectionReward (nCharNo, nReward)
		VALUES (@nCharNo, @nReward)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -3333
			RETURN
		END
	END
END

0㰁ࣣW
GO
