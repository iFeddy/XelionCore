CREATE  PROCEDURE [dbo].[usp_CharacterLoading]
/*
char loading save
2009.6 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@dWorldTime		datetime,
@sLoadingMapName	nvarchar(16),
@nFailType		smallint,
@nLoadingFailCnt	smallint,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	DECLARE @nSeqNo int
	SET @nSeqNo       = 0

	BEGIN TRAN

	SELECT TOP 1 @nSeqNo = nSeqNo
	  FROM tCharacterLoading
	 WHERE nCharNo = @nCharNo 
	 ORDER BY nSeqNo DESC

	SET @nSeqNo       = @nSeqNo + 1

	INSERT INTO tCharacterLoading( nCharNo, nSeqNo, dWorldTime, sLoadingMapName, nFailType, nLoadingFailCnt )
	VALUES (@nCharNo, @nSeqNo, @dWorldTime, @sLoadingMapName, @nFailType, @nLoadingFailCnt)
	IF @@ERROR <> 0 BEGIN
		-- 파일 에러처리
		SET @nRet = -1001
		ROLLBACK TRAN
		RETURN
	END

	SET @nRet = 0
	COMMIT TRAN
END

0ᴁᦢA
GO
