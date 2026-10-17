CREATE  PROCEDURE [dbo].[p_SameTimeWorld]
AS
BEGIN
	SET NOCOUNT ON
	SET LOCK_TIMEOUT 5000
declare 
@CurrentTime datetime,
@TimeText nvarchar(24)
SELECT TOP 1 @TimeText= convert(nvarchar(24),dDate,0), @CurrentTime=dDate  FROM tSameTimeChar WITH(READUNCOMMITTED)
WHERE dDate >= DATEADD(minute, -5, GETDATE() )
order by dDate desc
IF @@ROWCOUNT = 0 
	BEGIN
		SELECT GETDATE() AS CurTime, 0 AS CurCount
	END
ELSE
	BEGIN
		SELECT @CurrentTime AS CurTime , SUM(nCount) AS CurCount  FROM tSameTimeChar WITH(READUNCOMMITTED)
		WHERE convert(nvarchar(24),dDate,0) = @TimeText
	END
END

0䀁
GO
