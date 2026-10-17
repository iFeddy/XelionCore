CREATE PROCEDURE [dbo].[p_Guild_Member_SetGrade]
	@nNo     int
,	@nCharNo int
,	@nGrade  tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cnt int, @guildMaster int
	DECLARE @lastChangeDate datetime

	SET @nRet = 0
	SET @guildMaster = 0

	IF ISNULL(@nNo, 0) = 0 BEGIN
		SET @nRet = -1311
		RETURN
	END
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1313
		RETURN
	END
	
	-- Ä³ž¯ÅÍží º¯°æ Œ­ºñœº(ÀÎ°£°ü°è ÃÊ±âÈ­) ŸÆÀÌÅÛ Ãß°¡·Î ÀÎÇÑ ŒöÁ€ by kodong
	SELECT @cnt = ISNULL(COUNT(nCharNo), 0) 
	  FROM tCharacterChangeID 
	 WHERE nCharNo = @nCharNo 
	   AND bFlag = 0 -- ¹Ý¿µ ŸÈµÈ ±â·ÏÀÌ ÀÖŽÂ °æ¿ì
	   AND bInit = 1 -- ÀÎ°£°ü°è ÃÊ±âÈ­žŠ Œ³Á€ÇÑ °æ¿ì
	IF @nGrade = 0 AND @cnt = 1 BEGIN
		SET @nRet = -1323
		RETURN
	END
	ELSE BEGIN
		BEGIN TRAN
		
		-- ±æµåž¶œºÅÍ º¯°æ ¿äÃ»œÃ ±âÁž ±æµåž¶œºÅÍ Ä³ž¯ÅÍ¹øÈ£žŠ °¡Á®¿È
		IF @nGrade = 0 BEGIN
			SELECT @guildMaster = ISNULL(MAX(nCharNo), 0) FROM tGuildMember WHERE nNo = @nNo AND nGrade = 0
			
			-- ±æµåž¶œºÅÍ ÀÌŸç ÄðÅžÀÓ ÃŒÅ© by kodong 2011.9.30
			SELECT @lastChangeDate = MAX(dDate) FROM tLog_GuildMasterChange WITH(HOLDLOCK) WHERE nGuildNo = @nNo
			IF @lastChangeDate IS NOT NULL AND DATEADD(Day, 7, @lastChangeDate) > GETDATE() BEGIN
				ROLLBACK TRAN
				SET @nRet = -1324
				RETURN
			END
		END	
		
		-- ±æµåž¶œºÅÍ º¯°æ ¿äÃ»°ú ±âÁž ±æµåž¶œºÅÍ Ä³ž¯ÅÍ¹øÈ£°¡ ÁžÀçÇÏžé ±âÁž ±æµåž¶œºÅÍ žðµÎžŠ ÀÏ¹Ý ±æµå¿øÀž·Î º¯°æ		
		IF @nGrade = 0 AND @guildMaster <> 0 BEGIN
			UPDATE tGuildMember SET nGrade = 6 WHERE nNo = @nNo AND nGrade = 0
			IF @@ERROR <> 0 BEGIN
				ROLLBACK TRAN
				SET @nRet = -1325
				RETURN
			END
		END		
		
		-- ±æµå¿ø ±ÇÇÑ º¯°æ	
		UPDATE tGuildMember SET nGrade = @nGrade WHERE nNo = @nNo AND nCharNo = @nCharNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			ROLLBACK TRAN
			SET @nRet = -1325
		END
		ELSE BEGIN
			-- ±æµåž¶œºÅÍ ÀÌŸç ÄðÅžÀÓ ÃŒÅ© by kodong 2011.9.30				
			IF @nGrade = 0 BEGIN
				INSERT INTO tLog_GuildMasterChange (nGuildNo, nOldGuildMaster, nNewGuildMaster)
				VALUES (@nNo, @guildMaster, @nCharNo)
			END
			
			COMMIT TRAN
			SET @nRet = 0
		END
	END
END

0ﴁ烡Į
GO
