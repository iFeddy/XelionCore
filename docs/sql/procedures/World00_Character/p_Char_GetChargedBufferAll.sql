CREATE PROCEDURE [dbo].[p_Char_GetChargedBufferAll]
/*
Char Get charaged buffer all
2006.6 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int

AS

SET NOCOUNT ON

SELECT nKey, nID, dUseDate, dEndDate
FROM tCharacterChargedBuffer WHERE nCharNo = @nCharNo

-- end

0鰁丆Ŋ
GO
