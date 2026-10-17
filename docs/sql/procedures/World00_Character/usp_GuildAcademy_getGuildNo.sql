CREATE PROCEDURE [dbo].[usp_GuildAcademy_getGuildNo]
	@nCharNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nCharNo <= 0 BEGIN
		SET @nRet = -1512
		RETURN
	END

	IF EXISTS(SELECT nNo FROM tGuildAcademyMember WHERE nCharNo = @nCharNo) BEGIN
		SELECT nNo FROM tGuildAcademyMember WHERE nCharNo = @nCharNo
	END
	ELSE BEGIN
		-- 2008.8.28 일 업데이트시 서버프로그램의 수정 실수로 인해
		-- 정상적인 처리가 되지 않는 현상이 발생함.
		-- 2008.9.4 일 업데이트시 정상적으로 진행할 예정임.
		-- SET @nRet = -1521
		SET @nRet = 0
	END

	SET NOCOUNT OFF
END

0樁錅y
GO
