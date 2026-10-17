CREATE PROCEDURE [dbo].[usp_Statistics_setCharacterLevel]
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tStatCharLevel
	SELECT CONVERT(NVARCHAR(10), GETDATE(), 120) AS dDate, 
		   A.nLevel, 
		   ISNULL(A.nCnt, 0) AS 'Total(Before 2009-04-22)', 
		   ISNULL(B.nCnt, 0) AS 'Exp>0(Before 2009-04-22)', 
		   ISNULL(C.nCnt, 0) AS 'Total(After 2009-04-22)', 
		   ISNULL(D.nCnt, 0) AS 'Exp>0(After 2009-04-22)'
	  FROM (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND dCreateDate < '2009-04-22'
			 GROUP BY nLevel) A LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND nExp > 0
			   AND dCreateDate < '2009-04-22'
			 GROUP BY nLevel) B ON A.nLevel = B.nLevel LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND dCreateDate >= '2009-04-22'
			 GROUP BY nLevel) C ON A.nLevel = C.nLevel LEFT OUTER JOIN
		   (
			SELECT nLevel, COUNT(nCharNo) AS nCnt
			  FROM tCharacter WITH(NOLOCK)
			 WHERE bDeleted = 0 
			   AND nAdminLevel = 0 
			   AND nExp > 0
			   AND dCreateDate >= '2009-04-22'
			 GROUP BY nLevel) D ON A.nLevel = D.nLevel
	 ORDER BY A.nLevel

	SET NOCOUNT OFF
END

0
GO
