CREATE PROCEDURE [dbo].[p_Item_GetStorageEmptyList]
	@nOwner         int
,	@nStorageType   int
,	@nGetStorageCnt int = 1
,	@nRet           int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nOwner = -1 BEGIN
		SET @nRet = -1000
		RETURN
	END
	
	-- 동적 쿼리 작성 시작
	DECLARE @SQL AS NVARCHAR(3000)
	DECLARE @PARAMS AS NVARCHAR(1000)
	SET @SQL = N''
	SET @PARAMS = N''


	SET @SQL = '
		SELECT TOP ' + CAST(@nGetStorageCnt AS NVARCHAR) + ' A.nStorage
		FROM tItemStorageFindEmpty A LEFT OUTER JOIN (
			SELECT nStorage
			FROM tItem
			WHERE nOwner = ' + CAST(@nOwner AS NVARCHAR) + '
			AND nStorageType = ' + CAST(@nStorageType AS NVARCHAR) + ') B ON A.nStorage = B.nStorage
		WHERE B.nStorage IS NULL
		ORDER BY A.nStorage'
	SET @PARAMS = '@nGetStorageCnt INT, @nOwner INT, @nStorageType INT'

	-- 만들어진 동적 쿼리 실행
	EXEC sp_executesql @SQL, @PARAMS, @nGetStorageCnt, @nOwner, @nStorageType

	SET NOCOUNT OFF
END

0㴁鸩_
GO
