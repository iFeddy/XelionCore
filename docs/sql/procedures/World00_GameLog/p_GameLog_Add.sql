CREATE PROCEDURE [dbo].[p_GameLog_Add]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int,
@nBigint1 bigint
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3,  nBigint1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3, @nBigint1 )
SET @nLogNo = @@IDENTITY
-- end

0䌁넢Ķ
GO
