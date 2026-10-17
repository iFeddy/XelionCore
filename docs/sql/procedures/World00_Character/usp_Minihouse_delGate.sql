CREATE PROCEDURE [dbo].[usp_Minihouse_delGate]
	@nCharNo  int
,	@nItemID  int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF EXISTS(SELECT nCharNo FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID) BEGIN
		DELETE FROM tMinihouseGate WHERE nCharNo = @nCharNo AND nItemID = @nItemID
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1694
		END
	END
	ELSE BEGIN
		SET @nRet = -1622
	END

	SET NOCOUNT OFF
END

0謁곒
GO
