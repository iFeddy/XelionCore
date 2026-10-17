CREATE PROCEDURE [dbo].[p_Char_Find]
/*
Character Find
2004.6 By CJC
  input:
    sID = 찾을 ID
  output:
    0 = 없음
    1~? = 찾은 캐릭번호
*/
@sID nvarchar(40),
@nCharNo int OUTPUT
AS
SET NOCOUNT ON
SET @nCharNo = 0
SELECT @nCharNo = nCharNo FROM tCharacter(NOLOCK) WHERE sID = @sID
-- end

0⨁斾ň
GO
