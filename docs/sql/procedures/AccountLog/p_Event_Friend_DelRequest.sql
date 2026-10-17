CREATE PROCEDURE [p_Event_Friend_DelRequest]
@nEventUserNo 	int,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

DELETE FROM tEvent_Friend WHERE nEventUserNo= @nEventUserNo
IF @@ROWCOUNT = 0
	SET @err = 1
ELSE
	SET @err = 0	

RETURN @@ERROR




0錁Ĉ
GO
