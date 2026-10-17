CREATE PROCEDURE [dbo].[p_DayConnectionTime] 
AS
SET NOCOUNT ON
declare @nMon int
IF DAY(GETDATE()) = 1
set @nMon =MONTH(GETDATE())-1 
ELSE 
set @nMon = MONTH(GETDATE()) 
select min(nPlayMin) nMin, max(nPlayMin) nMax, avg(nPlayMin) nAvg, nWorldNo, DAY(dDate) nDay 
from [AccountLog].[dbo].[tAccountLog]			
where MONTH(dDate) = @nMon and nType = 2 
group by nWorldNo, DAY(dDate) 
order by nWorldNo, DAY(dDate) 
return 

0
GO
