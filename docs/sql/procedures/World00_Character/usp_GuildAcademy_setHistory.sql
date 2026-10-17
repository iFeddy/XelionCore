CREATE PROCEDURE [dbo].[usp_GuildAcademy_setHistory]
/*
Set Guild Token
2008.12 By Lunar
Input:
    See SQL
    
Output:
    nRet = -1,  Error(Not found guild)
    nRet = -2,  Error(Not found guild)
    nRet = 1,  OK
*/
@nAcademyID int,
@nCharNo int,
@sID nvarchar(50),
@nAction int,
@nRet int output

AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	INSERT INTO tGuildAcademyHistory (nAcademyID, nCharNo, sID, dDate, nAction)
	VALUES (@nAcademyID, @nCharNo, @sID, CURRENT_TIMESTAMP, @nAction)

	IF @@ERROR <> 0
	BEGIN
		SET @nRet = -2
		RETURN
	END

	SET NOCOUNT OFF
END
GO
