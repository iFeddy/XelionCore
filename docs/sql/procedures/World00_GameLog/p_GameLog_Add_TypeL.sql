CREATE PROCEDURE [dbo].[p_GameLog_Add_TypeL]
@nLogNo bigint OUTPUT,
@nType int,
@nCharNo int,
@nTargetCharNo int,
@nTargetID int,
@nItemKey bigint,
@nInt1 int,
@nInt2 int,
@nInt3 int,
@nBigint1 bigint
AS
BEGIN
SET NOCOUNT ON
SET @nLogNo = 0
INSERT tGameLog (  nType,  nCharNo,  nTargetCharNo,  nTargetID,  nItemKey,  nInt1,  nInt2,  nInt3, nBigint1 )
VALUES          ( @nType, @nCharNo, @nTargetCharNo, @nTargetID, @nItemKey, @nInt1, @nInt2, @nInt3, @nBigint1 )
SET @nLogNo = @@IDENTITY

END

0렚Ŋ
GO
