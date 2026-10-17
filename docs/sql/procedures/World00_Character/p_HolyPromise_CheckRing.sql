CREATE PROCEDURE [dbo].[p_HolyPromise_CheckRing]
/*
Holy promise check ring
2008.4 By CJC

 input:
	see SQL
	
 output:
	@nRet
	0 = No ring 
	1 = Yes ring
	-1 = DB error
*/
@nCharNo int,		-- 현재 캐릭
@nUpCharNo int,	-- 상위로 등록될 캐릭
@nRet int OUTPUT	-- 아미 상위로 등록되어 링구조가 되는지 결과확인값 리턴

AS

BEGIN
	SET NOCOUNT ON

	-- check are same char.
	IF @nCharNo = @nUpCharNo
	BEGIN
		SET @nRet = 1
		RETURN
	END

	DECLARE @nTempCharNo int
	DECLARE @nTempUpCharNo int

	SET @nTempCharNo = @nCharNo

	LOOP_UP:
	SET @nTempUpCharNo = 0
	SELECT @nTempUpCharNo = P.nUpCharNo FROM tHolyPromise AS P 
	INNER JOIN tCharacter AS C ON C.nCharNo =  P.nCharNo AND C.bDeleted = 0
	WHERE P.nCharNo = @nTempCharNo

	IF @nTempUpCharNo > 0
	BEGIN
		-- 링구조가 형성되는지 검사
		IF @nUpCharNo = @nTempUpCharNo
		BEGIN
			SET @nRet = 1
			RETURN
		END

		-- 다음 상위로 이동
		SET @nTempCharNo = @nTempUpCharNo
		GOTO LOOP_UP
	END

	-- 링구조 없음.
	SET @nRet = 0
	RETURN

END

0똁⺰Ł
GO
