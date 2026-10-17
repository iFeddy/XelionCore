CREATE PROCEDURE [dbo].[p_KingdomQuest_Set]
/*
Kingdom Quest Set
2005.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nQuestNo smallint,
@sData varchar(256),
@nRet int OUTPUT

AS

SET NOCOUNT ON

UPDATE tKingdomQuest SET sData = @sData WHERE nQuestNo = @nQuestNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tKingdomQuest ( nQuestNo, sData ) VALUES ( @nQuestNo, @sData )
END

SET @nRet = @@ERROR

-- end

0䬁孒Ż
GO
