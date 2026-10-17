CREATE PROCEDURE [dbo].[p_Char_ChangeClass]
@nCharNo int,
@nClass tinyint

AS
SET NOCOUNT ON

UPDATE tCharacterShape 
	SET 
		nClass=@nClass
	WHERE nCharNo=@nCharNo

return @@error

0쐁墄Ĺ
GO
