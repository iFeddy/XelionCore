CREATE PROCEDURE [dbo].[p_Quest_Add]
/*
Quest Set
2004.10 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nCharNo int,
@nQuestNo int,
@nStatus tinyint,
@sData varbinary(100),
@nRet int OUTPUT

AS

SET NOCOUNT ON

INSERT tQuest (  nCharNo,  nQuestNo,  nStatus,  sData )
VALUES        ( @nCharNo, @nQuestNo, @nStatus, @sData )

SET @nRet = @@ERROR

-- end

0码瑧ĉ
GO
