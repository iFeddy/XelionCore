CREATE PROCEDURE [dbo].[p_Prison_Get] 
@nCharID	int


AS
SET NOCOUNT ON

SELECT TOP 1 sReason, sRemark FROM tPrison
WHERE nCharID = @nCharID 
ORDER BY dPrisonDate DESC

RETURN @@error

0؁谟ć
GO
