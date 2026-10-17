CREATE PROCEDURE [dbo].[p_Char_GetUserNo]
/*
Get User No by Char ID
Input: nCharNo
Output:
    @nUserNo, 0 = Not found user
              ? = Found user No
*/
@nCharNo int,
@nUserNo int OUTPUT
AS
SET NOCOUNT ON
SET @nUserNo = 0
SELECT @nUserNo = nUserNo FROM tCharacter WHERE nCharNo = @nCharNo
-- END

0ᄁť
GO
