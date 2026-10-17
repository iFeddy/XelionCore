CREATE PROCEDURE [dbo].[p_MasPup_SetMaster]
/*
Set My Master
2004.12 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
			1 = 자신과는 스승관계를 맺을수 없다.
			2 = 제자가 4명이상 기록이 안된다.
			? = Error code(또는 스승이 2명이상 기록안된다)
*/
@nCharNo int,	-- As Pupils
@nMaster int,
@nFlags  int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 자신과 사제를 맺는지 확인한다.
IF @nMaster = @nCharNo
BEGIN
	SET @nRet = 1
	RETURN
END
-- 제자가 4명이상 기록이 안되게 확인
DECLARE @PupilsCount int
SELECT @PupilsCount = COUNT(*) FROM tMasterPupils WHERE nMaster = @nMaster
IF @PupilsCount >= 3
BEGIN
	SET @nRet = 2
	RETURN
END
-- 기록된 사항에 수정인 경우 수정처리 한다.
UPDATE tMasterPupils SET 
	nMaster = @nMaster,
	nFlags  = @nFlags
WHERE nPupils = @nCharNo
-- 처음 기록의 경우 처리한다.
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tMasterPupils (  nMaster,  nPupils,  nFlags )
	VALUES				 ( @nMaster, @nCharNo, @nFlags )
END
SET @nRet = @@ERROR
-- end
GO
