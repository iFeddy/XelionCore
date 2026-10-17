CREATE PROCEDURE [dbo].[p_Char_GetChargedBufferCount]
/*
Char Get charaged buffer all count
2006.6 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nCount int OUTPUT

AS

SET NOCOUNT ON

SELECT @nCount = COUNT(*) FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo

-- end

0！眴ŋ
GO
