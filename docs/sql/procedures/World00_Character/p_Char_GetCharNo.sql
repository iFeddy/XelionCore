CREATE PROCEDURE [dbo].[p_Char_GetCharNo]
/*
Get Char No by Char ID
Input: sID
Output:
    @nCharNo, 0 = Not found user
              ? = Found user No
*/
@sID nvarchar(40),
@nCharNo int OUTPUT
AS
SET NOCOUNT ON
SET @nCharNo = 0
SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID
-- END

0鬁钇ō
GO
