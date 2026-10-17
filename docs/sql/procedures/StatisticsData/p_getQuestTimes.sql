CREATE procedure [dbo].[p_getQuestTimes]
@nQuestNo int,
@top smallint
as
	set nocount on;
	select * from (
		select top (@top) 0 ServerNo, c.sID,nTimes, convert(char(17),dLastComplete,113) dLastComplete from w00_Character.dbo.tQuestTimes q, w00_Character.dbo.tCharacter c
		where nQuestNo=@nQuestNo and q.nCharNo=c.nCharNo and c.bDeleted=0 
		order by nTimes desc, q.dLastComplete  
	) as t0
	union all
	select * from (	
		select top (@top) 1 ServerNo, c.sID,nTimes, convert(char(17),dLastComplete,113) dLastComplete from w01_Character.dbo.tQuestTimes q, w01_Character.dbo.tCharacter c
		where nQuestNo=@nQuestNo and q.nCharNo=c.nCharNo and c.bDeleted=0 
		order by nTimes desc, q.dLastComplete  
	) as t1
0
GO
