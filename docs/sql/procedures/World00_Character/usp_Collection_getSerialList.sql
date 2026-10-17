CREATE PROCEDURE [dbo].[usp_Collection_getSerialList]
	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	SELECT nCardID, nCount
	  FROM tCollectionSerialCount 
END

0鄁ⱶT
GO
