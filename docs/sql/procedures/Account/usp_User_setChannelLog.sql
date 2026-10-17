CREATE PROCEDURE [dbo].[usp_User_setChannelLog]
	@userNo    int
,	@channel   nvarchar(20)
,	@issueDate nvarchar(20)
As
BEGIN
	SET NOCOUNT ON

	INSERT INTO tChannelLog (userNo, channel, issueDate)
	VALUES (@userNo, @channel, @issueDate)
	IF @@ERROR <> 0 OR @@IDENTITY = 0
		RETURN (0)
	ELSE
		RETURN (-1)

	SET NOCOUNT OFF
END

0
GO
