CREATE PROCEDURE [p_Event_Friend_EventComplete]
@nEventUserNo 	int,
@nFriendCharNo	int OUTPUT,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

SET @nFriendCharNo = 0

SELECT @nFriendCharNo = nFriendCharNo FROM tEvent_Friend WITH(READUNCOMMITTED) WHERE nEventUserNo = @nEventUserNo
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	RETURN @@ERROR
END

UPDATE tEvent_Friend SET  bEventComplete = 1 WHERE nEventUserNo = @nEventUserNo
IF @@ERROR <> 0
	SET @err = @@ERROR
ELSE
	SET @err = 0

RETURN @@ERROR




0찁ĉ
GO
