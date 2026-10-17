CREATE PROCEDURE [dbo].[p_Char_SaveLevelExpFrame] 
@nCharNo int,
@nLevel smallint,
@nExp bigint,
@nFame int

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nLevel=@nLevel,
		nExp=@nExp,
		nFame=@nFame
	WHERE nCharNo=@nCharNo

return @@error

0긁ꢹŰ
GO
