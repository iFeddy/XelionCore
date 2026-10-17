CREATE PROCEDURE [dbo].[usp_Char_GetRestExpExecTime]
/*
get cen reward
2009.4 By Luunar
 input:
	see SQL
 output:
	= 0 No Error
	< 0 Error
*/

@nCharNo		int,
@nRet			int OUTPUT

AS
BEGIN

	SET NOCOUNT ON

	IF EXISTS(SELECT nCharNo, dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo)
	BEGIN
		SELECT dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo
		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
--		SELECT nCharNo, getdate() FROM tRestExp WHERE nCharNo = @nCharNo
		SELECT getdate()
		SET @nRet = 0
		RETURN
	END
END

0
GO
