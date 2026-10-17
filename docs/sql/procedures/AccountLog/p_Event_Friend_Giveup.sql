CREATE PROCEDURE [p_Event_Friend_Giveup]
@nEventUserNo 	int,
@nFriendCharNo 	int OUTPUT,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

SET @nFriendCharNo = 0

SELECT @nFriendCharNo = nFriendCharNo FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
IF @@ROWCOUNT = 0
	BEGIN
		SET @err = 1
	END

ELSE

	BEGIN
		DELETE FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
		SET @err = 0
	END

RETURN @@ERROR
GO
