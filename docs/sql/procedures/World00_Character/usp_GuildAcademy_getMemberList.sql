CREATE PROCEDURE [dbo].[usp_GuildAcademy_getMemberList]
	@nNo  int
,	@nRet int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0

	IF @nNo <= 0 BEGIN
		SET @nRet = -1511
		RETURN
	END

	-- add G.nChatBanTime
	SELECT G.nCharNo, G.nJoinLevel, G.dJoinDate, G.nGrade, G.sIntro, G.nChatBanTime
		 , C.sID, C.nLevel, C.sLoginZone, C.dLastLoginDate
		 , S.nClass
	  FROM tGuildAcademyMember AS G INNER JOIN 
		   tCharacter          AS C ON C.nCharNo =  G.nCharNo INNER JOIN 
		   tCharacterShape     AS S ON S.nCharNo =  G.nCharNo
	 WHERE G.nNo = @nNo AND C.bDeleted = 0

	SET NOCOUNT OFF
END

0
GO
