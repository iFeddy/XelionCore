CREATE PROCEDURE [dbo].[p_Char_Delete]
/*
Character Delete
2004.6 By CJC
  Input:
    nCharNo = Character number to be delete.
  Output:
    0 = Delete OK
    1 = 해당 캐릭이 이미 지워진 캐릭임
    2 = 해당 캐릭을 찾을수 없음
    3 = 해당 캐릭을 지울수 없음
  Desc:
    - 지워진 캐릭의 ID는 '아뒤~날짜#시도횟수' 형태로 바뀐다.
    - 버그 방지로 아뒤에는 '~'는 사용되면 안된다.
*/
@nCharNo int,
@nRet int OUTPUT
AS
SET NOCOUNT ON
-- 이미 지워진 캐릭인지 조사
IF EXISTS(SELECT nCharNo FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo AND bDeleted > 0)
BEGIN
	SET @nRet = 1
	RETURN
END
-- 지워질 캐릭의 이름 얻기
DECLARE @sID nvarchar(40)
SET @sID = ''
SELECT @sID = sID FROM tCharacter(NOLOCK) WHERE nCharNo = @nCharNo
IF LEN(@sID) = 0
BEGIN
	SET @nRet = 2 -- 해당 캐릭이 없음
	RETURN
END
-- 아뒤내에 최초 삭제마크 이후의 문자는 모두 지운다.(지운거 또 지울수 있음으로 오류방지용)
DECLARE @DeleteMark nvarchar(2)
SET @DeleteMark = N'~'                       -- 삭제 마크 세팅
DECLARE @StartDeleteMark int
SET @StartDeleteMark = CHARINDEX( @DeleteMark, @sID )
IF @StartDeleteMark > 0
	SET @sID = LEFT( @sID, @StartDeleteMark-1 )
-- 지워질 캐릭의 ID를 변경
DECLARE @CurrentDate datetime
SET @CurrentDate = GETDATE()
DECLARE @nDeleteCount int
SET @nDeleteCount = 0
DECLARE @sChangedID nvarchar(40)
NAME_CHANGE: -- 이름 만들기 '아뒤' + 삭제마크 + '날짜' + '#시도횟수'
SET @sChangedID = 
	RTRIM(@sID) + 
	@DeleteMark + DATENAME(year, @CurrentDate) + DATENAME(month, @CurrentDate) + DATENAME(day, @CurrentDate) +
	N'#' + CAST( @nDeleteCount AS nvarchar(2) )
-- 캐릭 지우기
UPDATE tCharacter SET 
bDeleted = 1,
dDeletedDate = GETDATE(),
sID = @sChangedID
WHERE nCharNo = @nCharNo
-- 실패시 50번 이름을 바꾸면서 시도한다.
IF @@ERROR = 2601 -- 인덱스키 중복인 경우 
BEGIN
	SET @nDeleteCount = @nDeleteCount + 1
	IF @nDeleteCount < 50
		GOTO NAME_CHANGE
		
	-- 실패	
	SET @nRet = 3 
	RETURN
END	
-- 성공
SET @nRet = 0 
-- end

0戁쒜ŀ
GO
