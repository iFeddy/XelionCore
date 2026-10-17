CREATE PROCEDURE [usp_Log_setLogout]
	@WORLD_NO TINYINT		-- 서버
,	@USER_NO INT			-- 유저 번호
,	@CHAR_NO INT			-- 캐릭터 번호
,	@CLASS TINYINT			-- 캐릭터 클래스
,	@LEVEL SMALLINT			-- 캐릭터 현재레벨
,	@EXPERIENCE BIGINT		-- 획득 경험치
,	@PLAYTIME SMALLINT		-- 플레이시간
,	@RET INT = 0 OUTPUT		-- 리턴값(0:실패/1:성공)
AS

BEGIN

	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000

	INSERT INTO tLogout (worldNo, userNo, charNo, class, level, experience, playTime)
	VALUES (@WORLD_NO, @USER_NO, @CHAR_NO, @CLASS, @LEVEL, @EXPERIENCE, @PLAYTIME)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		SET @RET = 0
	ELSE
		SET @RET = 1

	SET NOCOUNT OFF
	SET LOCK_TIMEOUT -1

END

0∁靕ď
GO
