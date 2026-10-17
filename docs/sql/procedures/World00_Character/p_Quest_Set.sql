CREATE PROCEDURE [dbo].[p_Quest_Set]
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

--------------------------------------------------------- add start
-- added by CJC for for wedding quest 2008-06-03
DECLARE @nOldStatus tinyint
IF EXISTS(SELECT nStatus FROM tQuest WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo)
	SELECT @nOldStatus = nStatus FROM tQuest WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo
ELSE
	SET @nOldStatus = 0
/*-------------------------------------------------------------------------------------
	Project: FiestaHeroes - NA2016
	Author: xkl
	
	Quest item fix. Simply add the Quest ID and Item ID to the table, and it will
	prevent players from repeatedly turning in the same item for the quest.
------------------------------------------------------------------------------------*/
	DECLARE @QuestID INT, @ItemID INT;
	SELECT @QuestID = QuestID, @ItemID = ItemID FROM QuestItemFix WHERE QuestID = @nQuestNo;
	
	IF @nQuestNo = @QuestID AND @nStatus = 6
	BEGIN
		IF EXISTS (SELECT nItemID FROM tItem WHERE nOwner = @nCharNo AND nItemID = @ItemID)
		RETURN;
	END
/*-------------------------------------------------------------------------------------
	Continue as normal ~ 
------------------------------------------------------------------------------------*/

UPDATE tQuest SET 
    nStatus = @nStatus,
	sData = @sData
WHERE nCharNo = @nCharNo AND nQuestNo = @nQuestNo

IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tQuest (  nCharNo,  nQuestNo,  nStatus,  sData )
	VALUES        ( @nCharNo, @nQuestNo, @nStatus, @sData )
END
-- end of event on 2008-06-24 KST
--------------------------------------------------------- add start
--added by Amadeus for wedding quest 2008-05-29, modify by CJC 2008-06-03
-- nStatus 4 = PQS_REPEAT, Repeat quest done, now re-acceptable.
-- nStatus 6 = PQS_IN_PROGRESS 

/*IF @nStatus = 4 and @nQuestNo=20016 and @nOldStatus = 6
BEGIN
	IF EXISTS(SELECT nCharNo FROM tQuestTimes WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo)
		BEGIN
			UPDATE tQuestTimes SET nTimes=nTimes+1, dLastComplete=getdate()
			WHERE nCharNo=@nCharNo and nQuestNo=@nQuestNo
		END
	ELSE
		BEGIN
			INSERT tQuestTimes 
			VALUES(@nCharNo, @nQuestNo, 1, getdate())
		END
END
--------------------------------------------------------- end of edit by Amadeus
*/
SET @nRet = @@ERROR
--end

<
GO
