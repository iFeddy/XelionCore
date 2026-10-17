CREATE PROCEDURE [p_Report_Del]
@nRegNumber int,
@nRet smallint OUTPUT
AS
SET NOCOUNT ON
DECLARE	
@dReportDate	datetime,
@nResult	tinyint
--NoError(ErrorCode.h와 연동됨)
SET @nRet = 0
SELECT @dReportDate=dReportDate, @nResult=nResult FROM tReport WHERE nRegNumber=@nRegNumber 
--잘못된 등록번호
IF @@ROWCOUNT = 0
BEGIN
	SET @nRet = 1
	RETURN @@error
END
--1시간 지난 신고는 취소할수 없음.
IF @dReportDate < DATEADD(hour, -1, GETDATE())
BEGIN
	SET @nRet = 2
	RETURN @@error
END
--처리된 결과는 취소할 수 없음
IF @nResult <> 0
BEGIN
	SET @nRet = 3
	RETURN @@error
END
DELETE FROM tReport WHERE nRegNumber=@nRegNumber
RETURN @@error

0㼁尊Ĕ
GO
