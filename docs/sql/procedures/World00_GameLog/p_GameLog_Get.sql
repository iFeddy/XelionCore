CREATE PROCEDURE [dbo].[p_GameLog_Get]
/*
  GameLog Get - 주어진 조건에 맞는 로그를 모두검색해서 얻는다.
  input:
   - 주어진 조건절
  output:
   - Select된 데이타
*/
@WhereExpr nvarchar(256)
AS
SET NOCOUNT ON
DECLARE @SQLCMD nvarchar(512)
SET @SQLCMD = N'SELECT nLogNo, dDate, nCharNo, nType, sMap, nMapX, nMapY, nMapZ, nTargetCharNo, nTargetID, nItemKey, nInt1, nInt2, nInt3, nBigint1 FROM tGameLog WHERE ' + @WhereExpr
EXEC sp_executesql @SQLCMD
-- end

0⤁갿ŋ
GO
