CREATE  PROCEDURE [dbo].[p_Guild_Make2]
/*
Make Guild
2006.12.15 By CJC
Input:
	See SQL
	
Output:
	nNo = 0,   Error, Remain char money
	nNo = -1,  Error( Same guild name exist )
	nNo = -2,  Error( Not exists need item )
	nNo = -3,  Error( Other guild member )
	nNo = -4,  Error( Not enough maker money )
	nNo = -5,  Error( Make Guild error )
	nNo = -6,  Error( Maker to Join Guild master )
	nNo = 1~?, OK( Maked guild No )
*/
@sName nvarchar(18),
@sPassword nvarchar(9),
@nMaxMembers smallint,
@nType tinyint,
@nMakerCharNo int,
@nMakerPay bigint,
@nMakerNeedItemID int,
@nNo int output,
@nRemainCharMoney bigint output
AS
SET NOCOUNT ON
-- 에러코드 일단 세팅
SET @nNo = 0
-- 해당 길드가 있는지 확인
IF EXISTS (SELECT sName FROM tGuild(NOLOCK) WHERE sName = @sName)
BEGIN
	SET @nNo = -1
	RETURN
END
-- 필요아이템을 사용하는경우, 해당 아이템이 있는지 확인
IF @nMakerNeedItemID <> 65535
BEGIN
	IF NOT EXISTS (SELECT nItemID FROM tItem(NOLOCK) WHERE nOwner = @nMakerCharNo AND nItemID = @nMakerNeedItemID)
	BEGIN
		SET @nNo = -2
		RETURN
	END
END
-- 이미 다른 길드 가입했는지 확인
IF EXISTS (SELECT nCharNo FROM tGuildMember(NOLOCK) WHERE nCharNo = @nMakerCharNo)
BEGIN
	SET @nNo = -3
	RETURN
END
-- 길드 창설비용 적용
DECLARE @nRet tinyint
BEGIN TRAN
EXEC p_Char_MoneySub @nMakerCharNo, @nMakerPay, @nRet OUTPUT
IF @nRet = 0 -- 에러 확인
BEGIN
	ROLLBACK TRAN
	SET @nNo = -4
	RETURN
END
-- 길드 창설비용을 쓰고 남은돈 얻기
EXEC p_Char_MoneyGet @nMakerCharNo, @nRemainCharMoney OUTPUT, @nRet OUTPUT
IF @nRet = 0 -- 에러 확인
BEGIN
	ROLLBACK TRAN
	SET @nNo = -4
	RETURN
END
-- 길드 만들기
INSERT tGuild (  sName,  sPassword,  nMaxMembers,  nType )
VALUES        ( @sName, @sPassword, @nMaxMembers, @nType )
SET @nNo = @@IDENTITY
-- 길드 만들기 오류 확인
IF @@ERROR <> 0 
BEGIN
	ROLLBACK TRAN
	SET @nNo = -5
	RETURN
END
-- 만든 캐릭을 길마으로 추가
EXEC p_Guild_Member_Join @nNo, @nMakerCharNo, 0, @nRet OUTPUT
IF @nRet = 0
BEGIN 
	ROLLBACK TRAN
	SET @nNo = -6
	RETURN
END
-- 길드아카데미 생성
EXEC usp_GuildAcademy_setAcademy @nNo, @nMakerCharNo, 50, @nRet OUTPUT
IF @nRet <> 0 -- 길드아카데미 생성에서는 정상처리시 0을 return value로 가짐
BEGIN 
	ROLLBACK TRAN
	SET @nNo = @nRet
	RETURN
END
-- 만들기 성공
COMMIT TRAN
RETURN
-- end

0褁㧔ħ
GO
