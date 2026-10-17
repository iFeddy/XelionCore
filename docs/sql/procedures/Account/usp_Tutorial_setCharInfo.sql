CREATE PROCEDURE [dbo].[usp_Tutorial_setCharInfo]
	@nUserNo int
,	@nCharNo int
,	@nState  tinyint
,	@nStep   tinyint
,	@nRet    int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	SET @nRet = 0
	
	IF EXISTS(SELECT nCharNo FROM tTutorial WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo) BEGIN
		UPDATE tTutorial SET nState = @nState, nStep  = @nStep 
		 WHERE nUserNo = @nUserNo AND nCharNo = @nCharNo 
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -1
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tTutorial (nUserNo, nCharNo, nState, nStep)
		VALUES (@nUserNo, @nCharNo, @nState, @nStep)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -1
			RETURN
		END		
	END
END

0괁፳Ũ
GO
