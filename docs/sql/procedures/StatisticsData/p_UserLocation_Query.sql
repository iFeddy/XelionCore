CREATE PROCEDURE [dbo].[p_UserLocation_Query] 
/*
  UserLocation Query - 주어진 조건에 맞는 데이타를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET  @SQLCMD = N'SELECT nLogNo, nYear, nMonth, nDay, nWeek, nWeekDay, sLocation, nCount FROM tLocation WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0嬁譹
GO
