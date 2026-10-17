CREATE PROCEDURE [dbo].[usp_Dice_getDiceRank]
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

	SELECT TOP 20 C.sID, COUNT(B.nGameNo) AS nWinCount
	  FROM tDice AS A, tDiceLog AS B, tCharacter AS C
	 WHERE A.nGameNo = B.nGameNo
	   AND B.nCustomCharNo = C.nCharNo
	   AND A.nCharNo = @nCharNo
	   AND A.nItemKey = @nItemKey
	   AND B.bStatus = 1
	   AND C.bDeleted = 0
	 GROUP BY C.sID
	 ORDER BY nWinCount DESC

	SET NOCOUNT OFF
END

0夁춘[
GO
