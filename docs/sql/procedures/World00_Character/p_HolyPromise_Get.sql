CREATE PROCEDURE [dbo].[p_HolyPromise_Get]
/*
Get holy primise
2008.3 By CJC
 input:
	see SQL
 output:
	nRet > 0 fetch to take data
	= 0 No data
*/

@nCharNo int,
@nRet int OUTPUT

AS
BEGIN
	SET NOCOUNT ON
	SELECT @nRet=COUNT(nUpCharNo) FROM tHolyPromise WHERE nCharNo = @nCharNo
	IF @nRet > 0
	BEGIN
		IF EXISTS(
			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
			INNER JOIN tCharacter AS C ON C.nCharNo =  P.nUpCharNo AND C.bDeleted = 0
			WHERE P.nCharNo = @nCharNo)

			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, C.sID, C.nLevel,C.dLastLoginDate FROM tHolyPromise AS P
			INNER JOIN tCharacter AS C ON C.nCharNo =  P.nUpCharNo AND C.bDeleted = 0
			WHERE P.nCharNo = @nCharNo

		ELSE
			SELECT P.nUpCharNo, P.nFlags, P.dRegDate, P.dRejoinableDate, P.dMemberAcceptableDate, '', 1, GetDate() FROM tHolyPromise AS P WHERE P.nCharNo = @nCharNo

	END
END

0
GO
