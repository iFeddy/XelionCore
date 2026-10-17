CREATE PROCEDURE [dbo].[p_Item_SetTitle]
/*
Item Set Title
2004.8 By CJC
 input:
	see SQL
 output:
	@nRet	0 = OK, 
	        ? = Error code
*/
@nItemKey bigint,
@nTitleType smallint,
@nTitleCharNo int,
@sTitle nvarchar(24),
-- Output var
@nRet int OUTPUT
AS
SET NOCOUNT ON
UPDATE tItemTitle SET 
nTitleType = @nTitleType,
nTitleCharNo = @nTitleCharNo,
sTitle = @sTitle
WHERE nItemKey = @nItemKey
IF @@ERROR <> 0 OR @@ROWCOUNT = 0
BEGIN
	INSERT tItemTitle (  nItemKey,  nTitleType,  nTitleCharNo,  sTitle )
	VALUES			  ( @nItemKey, @nTitleType, @nTitleCharNo, @sTitle )
END
SET @nRet = @@ERROR
-- end

0䬁羋ů
GO
