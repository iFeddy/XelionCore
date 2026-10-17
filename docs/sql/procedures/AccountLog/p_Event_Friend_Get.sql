CREATE PROCEDURE [p_Event_Friend_Get]
	@nFriendUserNo	int

AS
BEGIN

	SET NOCOUNT ON
	 
	SELECT 
		t.nEventUserNo,
		t.nEventCharNo,
		t.nFriendCharNo,
		t.nWorldNo,
		t.bEventComplete,
		t.dDate
	FROM tEvent_Friend t
	WHERE t.nFriendUserNo = @nFriendUserNo

	SET NOCOUNT OFF

	RETURN @@ERROR
END
0
GO
