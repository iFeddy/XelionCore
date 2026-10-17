create procedure [dbo].[pFindCCU]
/*written by amadeus
*/
@start datetime,
@end	datetime
as
select dDate,w0,w1,w0+w1 total from (
select convert(char(16),dDate,120) dDate ,max(case nWorld when 0 then nCount else 0 end) w0, max(case nWorld when 1 then nCount else 0 end) w1
from tSameTimeChar where dDate between @start and @end
group by convert(char(16),dDate,120)
) as t1
GO
