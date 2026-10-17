CREATE PROCEDURE [dbo].[p_HolyPromise_DelChar]
/*
Delete char holy promise
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
	        
*/
@nDelCharNo int,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	DELETE FROM tHolyPromise WHERE nUpCharNo = @nDelCharNo
	SET @nRet = @@ERROR

	RETURN
END

0
GO
