CREATE PROCEDURE [dbo].[usp_Char_SetRestExpExecTime]
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

	-- 존재하지 않는 케릭터가 요청으로 들어오면 에러
	IF NOT EXISTS(SELECT nCharNo FROM tCharacter WHERE nCharNo = @nCharNo)
	BEGIN
		SET @nRet = -5654;
		RETURN
	END

	IF EXISTS(SELECT nCharNo, dRestExpLastExcuteTime FROM tRestExp WHERE nCharNo = @nCharNo)
	BEGIN
		UPDATE tRestExp
		SET dRestExpLastExcuteTime = getdate()
		WHERE nCharNo = @nCharNo

		SET @nRet = 0
		RETURN
	END
	ELSE
	BEGIN
		INSERT INTO tRestExp (nCharNo, dRestExpLastExcuteTime)
		VALUES (@nCharNo, getdate())

		SET @nRet = 0
		RETURN
	END
END

0꜁레Ĩ
GO
