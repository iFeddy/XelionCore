CREATE PROCEDURE [dbo].[p_OPLog_Add]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogType	tinyint,
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@nUserNo	int,
@sServer	nvarchar(50),
@nCharNo	int,
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
BEGIN
SET NOCOUNT ON
SET @nLogNo = 0
DECLARE @STMT	AS NVARCHAR(200)
SET @STMT = 'INSERT tOpLog ' + 
	CASE @nLogType
		WHEN 0 THEN		'( nOperNo,		sIP,	nType,	nUserNo,	sServer,	nCharNo,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@nCharNo,	@sDesc )'
		WHEN 1 THEN		'( nOperNo,		sIP,	nType,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @sDesc )'
		WHEN 2 THEN		'( nOperNo,		sIP,	nType,	nUserNo,	sServer,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@sDesc )'
		WHEN 3 THEN		'( nOperNo,		sIP,	nType,	sServer,	sDesc )
				VALUES	 ( @nOperNo,	@sIP,	@nType, @sServer,	@sDesc )'
	END
PRINT @STMT
EXEC sp_executesql @STMT
SET @nLogNo = @@IDENTITY
END
-- end

0긁蒀ż
GO
