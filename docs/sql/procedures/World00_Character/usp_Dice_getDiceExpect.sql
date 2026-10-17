CREATE PROCEDURE [dbo].[usp_Dice_getDiceExpect]
	@nCharNo  int
,	@nItemKey bigint
,	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF ISNULL(@nCharNo, 0) <= 0 BEGIN
		SET @nRet = -2411
		RETURN
	END

	IF ISNULL(@nItemKey, 0) <= 0 BEGIN
		SET @nRet = -2412
		RETURN
	END

	SELECT ISNULL(COUNT(DISTINCT(A.nGameNo)), 0) AS nCount, ISNULL(SUM(nTotalDividend), 0) AS nTotalDividend
	  FROM tDice AS A, tDiceResult AS B
	 WHERE A.nGameNo = B.nGameNo
	   AND A.nCharNo = @nCharNo
	   AND A.nItemKey = @nItemKey

	SET NOCOUNT OFF
END

0Y
GO
