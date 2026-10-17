CREATE PROCEDURE [dbo].[usp_Minihouse_getGateList]
	@nCharNo  int
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1611
		RETURN
	END

	IF NOT EXISTS(SELECT nCharNo FROM tMinihouse WHERE nCharNo = @nCharNo) BEGIN
		INSERT INTO tMinihouse (nCharNo) VALUES (@nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1692
			RETURN
		END
	END

	SELECT nCharNo, nItemID, bType, ISNULL(dEndDate, 0) AS dEndDate
	  FROM tMinihouseGate
	 WHERE nCharNo = @nCharNo AND (dEndDate > GETDATE() OR dEndDate IS NULL)

	SET NOCOUNT OFF
END

0쐁ꃶ!
GO
