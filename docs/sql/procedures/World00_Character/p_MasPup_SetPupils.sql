CREATE PROCEDURE [dbo].[p_MasPup_SetPupils]
/*
Set My Pupils
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK,또는 이미 기록 되어 있음 
			1 = 자신과 사제를 맺을수 없다.
			2 = 제자가 4명이상 기록이 안된다.
			? = Error code(또는 스승이 2명이상 기록안된다)
*/
@nCharNo int, -- As Master
@nPupils int,
@nFlags  int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 자신과 사제를 맺는지 확인한다.
IF @nPupils = @nCharNo
BEGIN
	SET @nRet = 1
	RETURN
END
-- 이미 기록이 된경우 확인, Flag만 바꾼다.
IF EXISTS( SELECT * FROM tMasterPupils WHERE (nMaster = @nCharNo) AND ( nPupils = @nPupils ) )
BEGIN
	UPDATE tMasterPupils SET nFlags = @nFlags WHERE (nMaster = @nCharNo) AND (nPupils = @nPupils)
	SET @nRet = @@ERROR
	RETURN
END
-- 제자가 4명이상 기록이 안되게 확인
DECLARE @PupilsCount int
SELECT @PupilsCount = COUNT(*) FROM tMasterPupils WHERE nMaster = @nCharNo
IF @PupilsCount >= 3
BEGIN
	SET @nRet = 2
	RETURN
END
-- 새로 추가된 제자 기록(스승이 2명이상 기록안되게 체크는 nPupils가 Alternative key로 확인된다)
INSERT tMasterPupils (  nMaster,  nPupils,  nFlags )
VALUES				 ( @nCharNo, @nPupils, @nFlags )
SET @nRet = @@ERROR
-- end

0ࠁ＝
GO
