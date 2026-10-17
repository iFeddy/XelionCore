CREATE PROCEDURE [dbo].[p_GameLog_Add_Type7]
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
@nInt1 int,
@nInt2 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nInt1,  nInt2 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nInt1, @nInt2 )
SET @nLogNo = @@IDENTITY
-- end

0툁帟Ľ
GO
