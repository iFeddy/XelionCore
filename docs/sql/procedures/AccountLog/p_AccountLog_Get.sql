CREATE PROCEDURE [p_AccountLog_Get] 
/*
  AccountLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nchar(512)
SET  @SQLCMD = N'SELECT nLogNo, dDate, nType, nUserNo, nWorldNo, nPlayMin, nErrorCode, sIP, sUserID, sUserPassword FROM tAccountLog WITH(NOLOCK) WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0㴁ㅾă
GO
