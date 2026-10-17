CREATE PROCEDURE [dbo].[usp_Minihouse_getVisitorCount]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		SELECT nVisitor FROM tMinihouse WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END

		SELECT nVisitor FROM tMinihouse WHERE nCharNo = @nCharNo
	END

	SET NOCOUNT OFF
END

0ﴁ锚"
GO
