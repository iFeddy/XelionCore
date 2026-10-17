CREATE PROCEDURE [p_Event_WorldCup_Enter]
@nCharNo	int,
@nUserNo	int,
@nLevel		tinyint,
@nWorldNo	tinyint,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

IF GETDATE() > '2006-06-07'
BEGIN
	SET @err = 3
	RETURN @@ERROR
END

IF EXISTS ( SELECT nEventCharNo FROM tEventWorldCup WHERE nEventCharNo=@nCharNo AND nLevel=@nLevel AND nWorldNo=@nWorldNo )
	BEGIN
		SET @err = 1
	END 
ELSE
	BEGIN
		INSERT INTO tEventWorldCup(nEventCharNo, nEventUserNo, nLevel, nWorldNo)
			VALUES(@nCharNo, @nUserNo, @nLevel, @nWorldNo )
		IF @@ERROR <> 0 --Error
			SET @err = 2
		ELSE
			SET @err = 0
	END

RETURN @@ERROR



0嬁譹Đ
GO
