CREATE PROCEDURE [dbo].[p_OPLog_Get]
/*
  OPLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT nNo, dDate, nOperNo, sIP, nType, nUserNo, sServer, nCharNo, sDesc FROM tOpLog ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0碤Ž
GO
