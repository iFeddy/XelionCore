CREATE PROCEDURE [dbo].[p_GameLog_GetCount]
/*
  GameLog Get Count - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - @Count = 해당 조건에 만족하는 레코드수
*/
@WhereExpr nvarchar(256),
@Count int OUTPUT
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT @RowCount = COUNT(*) FROM tGameLog WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD, N'@RowCount AS int OUTPUT', @RowCount = @Count OUTPUT
-- end

0戁ꁣŌ
GO
