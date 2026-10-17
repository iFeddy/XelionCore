CREATE PROCEDURE [dbo].[p_HolyPromise_Del]
/*
Delete holy promise
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = OK, 
	? = Other DB error code
	        
*/
@nCharNo int,
@nRet int OUTPUT

AS

BEGIN
	SET NOCOUNT ON

	-- Update/Inset data
	IF EXISTS( SELECT nCharNo FROM  tHolyPromise WHERE nCharNo=@nCharNo)
	BEGIN
		UPDATE tHolyPromise SET nUpCharNo=0, nFlags=0, dRegDate=GetDate() WHERE nCharNo=@nCharNo
		SET @nRet = @@ERROR
	END
	ELSE
		SET @nRet = 0

	RETURN
END

0
GO
