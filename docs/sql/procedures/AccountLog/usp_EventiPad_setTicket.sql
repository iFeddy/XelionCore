CREATE PROCEDURE dbo.usp_EventiPad_setTicket
	@nUserNo   int,
	@nPlayTime int,
	@nRet      int = 0 output
AS
BEGIN
	SET NOCOUNT ON
	
	DECLARE @nTicket AS int SET @nTicket = 0
	SET @nRet = 0
	
	IF ISNULL(@nPlayTime, 0) < 60 BEGIN
		SET @nRet = -1
		RETURN
	END
	ELSE BEGIN
		SET @nTicket = @nPlayTime / 60
	END
	
	IF EXISTS(SELECT nUserNo FROM tEventiPad WHERE nUserNo = @nUserNo) BEGIN
		UPDATE tEventiPad SET
			nTicket = nTicket + @nTicket,
			sys_update_date = GETDATE()
		 WHERE nUserNo = @nUserNo
		IF @@ERROR <> 0 OR @@ROWCOUNT <> 1 BEGIN
			SET @nRet = -8
			RETURN
		END
	END
	ELSE BEGIN
		INSERT INTO tEventiPad (nUserNo, nTicket)
		VALUES (@nUserNo, @nTicket)
		IF @@ERROR <> 0 BEGIN
			SET @nRet = -7
			RETURN
		END
	END
	
	SET NOCOUNT OFF
END
0䌁넢Ķ
GO
