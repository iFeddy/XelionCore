CREATE PROCEDURE [dbo].[p_Char_SavePKCount] 
/*
Character Set nPKCount
2007.2 By CJC
 input:
   All args
 output: 
   @@error
*/
@nCharNo int,
@nPKCount int

AS
SET NOCOUNT ON

UPDATE tCharacter SET nPKCount=@nPKCount WHERE nCharNo=@nCharNo

return @@error

0䨁옌Ų
GO
