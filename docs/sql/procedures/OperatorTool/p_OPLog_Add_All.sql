CREATE PROCEDURE [dbo].[p_OPLog_Add_All]
/*
 Operator log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nOperNo	smallint,
@sIP		nvarchar(20),
@nType		smallint,
@nUserNo	int,
@sServer	nvarchar(50),
@nCharNo	int,
@sDesc		nvarchar(200),
@nLogNo		int OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tOpLog ( nOperNo,	sIP,	nType,	nUserNo,	sServer,	nCharNo,	sDesc )
VALUES        ( @nOperNo,	@sIP,	@nType, @nUserNo,	@sServer,	@nCharNo,	@sDesc )
SET @nLogNo = @@IDENTITY
-- end

0㸁웄ċ
GO
