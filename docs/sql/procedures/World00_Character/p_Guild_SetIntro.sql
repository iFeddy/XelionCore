CREATE PROCEDURE [dbo].[p_Guild_SetIntro]

/*
Set Guild Intro
2006.2 By CJC

Input:
	See SQL
	
Output:
	nRet = 0,  Error(Not found guild)
	nRet = 1,  OK
*/

@nNo int,
@sIntro nvarchar(128),

@nRet tinyint output

AS

SET NOCOUNT ON
SET @nRet = 0

UPDATE tGuild SET sIntro = @sIntro WHERE nNo = @nNo
SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = 0
	RETURN
END

-- end

<䜍Ķ
GO
