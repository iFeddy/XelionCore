CREATE PROCEDURE [dbo].[p_HolyPromise_GetMember]
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

	SELECT @nNumOfMember=COUNT( P.nCharNo) FROM tHolyPromise AS P
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nUpCharNo = @nCharNo

	IF @nNumOfMember > 0 
	BEGIN
		SELECT P.nCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
		INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
		WHERE P.nUpCharNo = @nCharNo
	END
END

0
GO
