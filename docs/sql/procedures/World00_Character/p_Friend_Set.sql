CREATE PROCEDURE [dbo].[p_Friend_Set]
/*
Char Set friend
2004.12 By CJC

 input:
	see SQL
	
 output:
	@nRet	0 = OK, 
			1 = friend and me are same char
			2 = i'm over max friend count
			3 = friend is over max friend count
			547 = DB error, nCharNo or nFriendCharNo is not existed Character(Foreign Key checked)
			2627 = DB error, already friend marked
	        ? = Other DB error code
	        
*/
@nMAXFRIEND smallint,
@nCharNo int,
@nFriendCharNo int,
@nFlags tinyint,
@nRet int OUTPUT

AS

SET NOCOUNT ON

-- check firend and me are same char.
IF @nCharNo = @nFriendCharNo
BEGIN
	SET @nRet = 1
	RETURN
END

-- check over max firend
DECLARE @nFriendCount int
EXEC p_Friend_GetNumOfFriend @nCharNo, @nFriendCount OUTPUT
IF @nFriendCount >= @nMAXFRIEND
BEGIN
	SET @nRet = 2
	RETURN
END

EXEC p_Friend_GetNumOfFriend @nFriendCharNo, @nFriendCount OUTPUT
IF @nFriendCount >= @nMAXFRIEND
BEGIN
	SET @nRet = 3
	RETURN
END

-- add friend
BEGIN TRAN

INSERT tFriend (  nCharNo,  nFriendCharNo,  nFlags )
VALUES         ( @nCharNo, @nFriendCharNo, @nFlags )
SET @nRet = @@ERROR
if( @nRet <> 0 )
BEGIN
	ROLLBACK TRAN
	RETURN
END

-- add friend cross
INSERT tFriend (  nCharNo,  nFriendCharNo,  nFlags )
VALUES         ( @nFriendCharNo, @nCharNo, @nFlags )
SET @nRet = @@ERROR
if( @nRet <> 0 )
BEGIN
	ROLLBACK TRAN
	RETURN
END

COMMIT TRAN
RETURN

-- end

0댁錗ě
GO
