CREATE PROCEDURE [dbo].[usp_Character_setChangeID]
	@nRowNo  int
,	@nCharNo int
,	@sID     nvarchar(16)
,	@bInit   bit
,	@nUserNo int
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	DECLARE @sOldID AS nvarchar(50)
	DECLARE @nDBUserNo AS int
	DECLARE @nGrade AS tinyint
	DECLARE @nMyGuildNo AS int
	
	SELECT @sOldID = T1.sID, @nDBUserNo = T1.nUserNo, @nGrade = nGrade
	  FROM tCharacter T1 LEFT OUTER JOIN tGuildMember T2 ON T1.nCharNo = T2.nCharNo
	 WHERE T1.nCharNo = @nCharNo AND T1.bDeleted = 0
	   AND bDeleted = 0
	IF @sOldID = @sID BEGIN
		SET @nRet = -3039
		RETURN
	END	
	IF @nDBUserNo IS NULL BEGIN
		SET @nRet = -3023
		RETURN
	END
	ELSE BEGIN
		-- mismatch userNo
		IF @nDBUserNo <> @nUserNo BEGIN
			SET @nRet = -3024
			RETURN
		END
		
		-- impossible guildmaster
		IF @nGrade IS NOT NULL AND @nGrade = 0 AND @bInit = 1 BEGIN
			SET @nRet = -3025
			RETURN
		END

		IF NOT EXISTS(SELECT nCharNo FROM tCharacterChangeID WHERE bFlag = 0 AND nCharNo = @nCharNo) BEGIN
			SET @nRet = -3024
			RETURN
		END

		BEGIN TRAN
		
		-- change ID
		IF EXISTS(SELECT nCharNo FROM tCharacter WHERE sID = @sID) BEGIN		
			SET @nRet = -3039
			ROLLBACK TRAN
			RETURN
		END
		ELSE BEGIN
			UPDATE tCharacter
			   SET sID = @sID
			 WHERE nCharNo = @nCharNo
			IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
				SET @nRet = -3026
				ROLLBACK TRAN
				RETURN
			END		
		END
		
		-- init Community
		IF @bInit = 1 BEGIN
			-- delete Friend
			IF EXISTS(SELECT nCharNo FROM tFriend WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tFriend WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3027
					ROLLBACK TRAN
					RETURN
				END
			END
			IF EXISTS(SELECT nCharNo FROM tFriend WHERE nFriendCharNo = @nCharNo) BEGIN
				DELETE FROM tFriend WHERE nFriendCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3027
					ROLLBACK TRAN
					RETURN
				END				
			END
			
			-- delete Marriage
			IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tMarriage WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3028
					ROLLBACK TRAN
					RETURN
				END
			END
			IF EXISTS(SELECT nCharNo FROM tMarriage WHERE nPartnerCharNo = @nCharNo) BEGIN				
				DELETE FROM tMarriage WHERE nPartnerCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3028
					ROLLBACK TRAN
					RETURN
				END				
			END			

			-- delete HolyPromise me
			IF EXISTS(SELECT nCharNo FROM tHolyPromise WHERE nCharNo = @nCharNo) BEGIN
				DELETE FROM tHolyPromise WHERE nCharNo = @nCharNo
				IF @@ERROR <> 0 OR @@ROWCOUNT < 1 BEGIN
					SET @nRet = -3029
					ROLLBACK TRAN
					RETURN
				END
			END
GO
