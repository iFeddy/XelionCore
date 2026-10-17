CREATE PROCEDURE [dbo].[usp_User_setDrawComment]
	@userNo  int
,	@causeNo tinyint = 0
,	@comment nvarchar(1000) = ''
,	@userIP  nvarchar(30)
AS
BEGIN
	SET NOCOUNT ON

	INSERT INTO tDrawComment (nUserNo, nCauseNo, sComment, sUserIP)
	VALUES (@userNo, @causeNo, @comment, @userIP)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		RETURN (-9)
	ELSE
		RETURN (0)

	SET NOCOUNT OFF
END

!bNôÙ¾ dF(ä¿u PF+
GO
