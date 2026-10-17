CREATE PROCEDURE [dbo].[p_GameLog_Add_Type3]
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
@nInt1 int
AS
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  sMap,  nMapX,  nMapY,  nMapZ,  nTargetCharNo,  nTargetID,  nInt1 )
VALUES          ( @nType, @nCharNo, @sMap, @nMapX, @nMapY, @nMapZ, @nTargetCharNo, @nTargetID, @nInt1 )
SET @nLogNo = @@IDENTITY
-- end

0趎Ĺ
GO
