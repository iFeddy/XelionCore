CREATE PROCEDURE [dbo].[p_Char_AddExp]
@nCharNo int,
@nExp bigint

AS
SET NOCOUNT ON

UPDATE tCharacter 
	SET 
		nExp=nExp+@nExp
	WHERE nCharNo=@nCharNo

return @@error

0趎Ĺ
GO
