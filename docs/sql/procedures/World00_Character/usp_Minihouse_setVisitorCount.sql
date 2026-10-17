CREATE PROCEDURE [dbo].[usp_Minihouse_setVisitorCount]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		UPDATE tMinihouse SET nVisitor = nVisitor + 1 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1693
		END
	END
	ELSE BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END

		UPDATE tMinihouse SET nVisitor = nVisitor + 1 WHERE nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1693
		END
	END

	SET NOCOUNT OFF
END

0漁絣$
GO
