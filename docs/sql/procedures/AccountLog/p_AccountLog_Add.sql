CREATE PROCEDURE [dbo].[p_AccountLog_Add] 
/*
 Account log add
 input:
  로그에 기록될 정보들(See SQL)
 output: 
  nLogNo = 추가된 로그의 번호(에러시 LogNo = 0 )
*/
@nType		tinyint,
@nUserNo	int,
@nWorldNo	tinyint,
@nPlayMin	smallint,
@nErrorCode smallint,
@sIP		varchar(16),
@sUserID	nvarchar(16),
@sUserPassword	nvarchar(16),
@nLogNo bigint OUTPUT
AS
SET NOCOUNT ON
SET @nLogNo = 0

/*-------------------------------------------------------------------------------------
Project: FiestaHeroes - NA2016
Author: xkl

Update the IP associated with the account upon login.
------------------------------------------------------------------------------------*/
IF @nType = 1
BEGIN
    UPDATE Account..tUser SET sUserIP = @sIP Where nUserNo = @nUserNo
END
/*-------------------------------------------------------------------------------------
 Continue as normal ~ 
------------------------------------------------------------------------------------*/

INSERT tAccountLog  WITH (ROWLOCK) (  nType,  nUserNo,  nWorldNo,  nPlayMin,  nErrorCode,  sIP,  sUserID,  sUserPassword )
VALUES             ( @nType, @nUserNo, @nWorldNo, @nPlayMin, @nErrorCode, @sIP, @sUserID, @sUserPassword )
SET @nLogNo = @@IDENTITY
-- end
GO
