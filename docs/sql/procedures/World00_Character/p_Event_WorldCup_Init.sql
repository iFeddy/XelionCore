CREATE PROCEDURE [dbo].[p_Event_WorldCup_Init]
@nCharNo int,
@bNewJoiner tinyint OUTPUT

AS
SET NOCOUNT ON

DECLARE
@dRegisterDate	datetime,
@dEventStartDate	datetime,
@dEventEndDate	datetime


SET @dEventStartDate = '2006-05-30' --테스트
SET @dEventEndDate =  '2006-06-07' --6월 06일까지.

IF GETDATE() > @dEventEndDate
BEGIN
	SET	@bNewJoiner = 0
	RETURN @@error
END

SELECT @dRegisterDate=dCreateDate   
	FROM tCharacter
	WHERE nCharNo = @nCharNo --AND 

IF @@ROWCOUNT = 0
BEGIN
	SET	@bNewJoiner = 0
	RETURN @@error
END

IF @dRegisterDate > @dEventStartDate AND @dRegisterDate < @dEventEndDate
	BEGIN
		SET	@bNewJoiner = 1
	END
ELSE
	BEGIN
		SET	@bNewJoiner = 0
	END

RETURN @@error

0␁Ĕ
GO
