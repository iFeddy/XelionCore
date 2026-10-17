CREATE PROCEDURE [dbo].[p_HolyPromise_GetMemberCount]
/*
Get primise member count
2008.3 By CJC
 input:
	see SQL
 output:
	fetch to take data
*/

@nCharNo int,
@nNumOfMember int OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SET @nNumOfMember = 0
	SELECT @nNumOfMember = COUNT(P.nUpCharNo) FROM tHolyPromise AS P 
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nUpCharNo = @nCharNo
END

<
GO
