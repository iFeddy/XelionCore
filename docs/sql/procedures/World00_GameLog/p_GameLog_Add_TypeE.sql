CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeE]
/*
 Game log add
 input:
  All Args = 로그에 기록될 정보
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nLogNo bigint OUTPUT,
@nType int,
@sMap varchar(16),
@nMapX int,
@nMapY int,
@nMapZ int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetID,  nItemKey,  nInt1,  nInt2 )
VALUES          ( @nType, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetID, @nItemKey, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end
GO
