CREATE PROCEDURE [p_Event_Friend_SetHideMsg]
@nUserNo 	int,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

IF EXISTS (SELECT nUserNo FROM tEvent_Friend_HideMsg WITH(READUNCOMMITTED)
		WHERE nUserNo= @nUserNo )
	BEGIN
		SET @err = 1
	END
ELSE
	BEGIN
		INSERT INTO tEvent_Friend_HideMsg(nUserNo) VALUES(@nUserNo)
		IF @@ERROR <> 0
			SET @err = @@ERROR
		ELSE
			SET @err = 0
	END


RETURN




0뀁꼌č
GO
