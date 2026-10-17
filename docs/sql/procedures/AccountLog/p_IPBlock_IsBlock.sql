CREATE  PROCEDURE p_IPBlock_IsBlock
/*
 Is IP block ?
 input:
  IP
 output: 
  bIsBlock, 0=No block, 1=Block
*/
	@nUserNo    int
,	@sIP 		nvarchar(20)
,	@bIsBlock	tinyint OUTPUT
AS
BEGIN
	SET NOCOUNT ON

	SET @bIsBlock = 0

	IF EXISTS(SELECT sIPStart FROM tIPBlock WHERE @sIP >= sIPStart AND @sIP <= sIPEnd) BEGIN
		IF NOT EXISTS(SELECT TOP 1 nLogNo FROM tAccountLog WHERE nUserNo = @nUserNo AND nType = 1) BEGIN
			SET @bIsBlock = 1
		END
	END

	SET NOCOUNT OFF
END
-- end
GO
