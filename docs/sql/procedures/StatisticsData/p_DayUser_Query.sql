CREATE PROCEDURE [dbo].[p_DayUser_Query] 
/*
  DayUser Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nIndex, dDate, nUserCount, nConnectCount FROM tDayUser WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0෪Ć
GO
