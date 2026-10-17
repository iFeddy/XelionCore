CREATE PROCEDURE [dbo].[p_HolyPromise_SetDate]
/*
Holy promise set date
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
*/
@nCharNo int,
@dRejoinableDate datetime,
@dMemberAcceptableDate datetime,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	-- Update data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		UPDATE tHolyPromise SET dRejoinableDate=@dRejoinableDate, dMemberAcceptableDate=@dMemberAcceptableDate WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
	END
	ELSE	
	BEGIN
		INSERT tHolyPromise (  nCharNo,  nUpCharNo,  nFlags, dRejoinableDate, dMemberAcceptableDate )
		VALUES ( @nCharNo, 0, 0, @dRejoinableDate, @dMemberAcceptableDate )
		SET @nRet = @@ERROR
	END
	
	RETURN
END

0
GO
