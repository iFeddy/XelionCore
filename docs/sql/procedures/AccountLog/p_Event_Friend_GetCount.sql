CREATE PROCEDURE [p_Event_Friend_GetCount]
@nUserNo int,
@nCount int OUTPUT 

AS
SET NOCOUNT ON


SELECT @nCount = COUNT(*) FROM tEvent_Friend WITH(READUNCOMMITTED)
WHERE nFriendUserNo = @nUserNo

RETURN @@error



0㸁웄ċ
GO
