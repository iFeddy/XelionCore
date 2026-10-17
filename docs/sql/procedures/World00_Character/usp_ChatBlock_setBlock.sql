CREATE PROCEDURE [dbo].[usp_ChatBlock_setBlock]
	@nUserNo int
,	@sID     nvarchar(40) = ''
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	DECLARE @cnt AS int
	DECLARE @maxCnt AS int
	DECLARE @nCharNo AS int

	SET @nRet = 0
	SET @nCharNo = 0
	SET @maxCnt = 100 -- 최대등록 차단수(차후 변경시 DB 저장프로시저만 수정하면 반영됨)

	IF ISNULL(@nUserNo, 0) <= 0 BEGIN
		SET @nRet = -2311
		RETURN
	END

	IF LEN(ISNULL(@sID, '')) = 0 BEGIN
		SET @nRet = -2312
		RETURN
	END

	SELECT @nCharNo = nCharNo FROM tCharacter WHERE sID = @sID AND bDeleted = 0
	IF @nCharNo <= 0 BEGIN
		SET @nRet = -2323
		RETURN
	END

	IF EXISTS(SELECT nUserNo FROM tChatBlock WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		SET @nRet = -2321
		RETURN
	END

	SELECT @cnt = ISNULL(COUNT(*), 0) FROM tChatBlock WHERE nUserNo = @nUserNo
	IF @cnt > @maxCnt BEGIN
		SET @nRet = -2322
		RETURN
	END		
	ELSE BEGIN
		INSERT INTO tChatBlock (nUserNo, nCharNo)
		VALUES (@nUserNo, @nCharNo)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -2392
			RETURN
		END

		SELECT T1.nCharNo, T2.sID, T2.nAdminLevel, T2.bDeleted
		  FROM tChatBlock T1, tCharacter T2
		 WHERE T1.nUserNo = @nUserNo 
		   AND T1.nCharNo = @nCharNo
		   AND T1.nCharNo = T2.nCharNo
	END
END

0ᄁ㺓N
GO
