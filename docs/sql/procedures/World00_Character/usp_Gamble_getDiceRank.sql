CREATE PROCEDURE [dbo].[usp_Gamble_getDiceRank]
	@nRet     int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	SELECT TOP 20 C.sID, COUNT(B.nGameNo) AS nWinCount
	  FROM tGambleDice AS A, tGambleDiceLog AS B, tCharacter AS C
	 WHERE A.nGameNo = B.nGameNo
	   AND B.nCharNo = C.nCharNo
	   AND B.bStatus = 1
	   AND C.bDeleted = 0
	 GROUP BY C.sID
	 ORDER BY nWinCount DESC

	SET NOCOUNT OFF
END

0
GO
