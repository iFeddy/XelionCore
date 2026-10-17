CREATE PROCEDURE [dbo].[usp_Collection_getFavoritesList]
	@nCharNo      int
,	@nRet         int = 0 output
AS
BEGIN
	SET NOCOUNT ON

	SET @nRet = 0
	
	IF ISNULL(@nCharNo, 0) = 0 BEGIN
		SET @nRet = -1111
		RETURN
	END

	SELECT T1.nSlotNo, T1.nCardID, T2.nSerialNo, T2.nStarNumber 
	  FROM tCollectionFavorites T1, tCollectionCard T2 
	 WHERE T1.nCharNo = T2.nCharNo
	   AND T1.nCardID = T2.nCardID 
	   AND T1.nCharNo = @nCharNo 
END

0堁㡒S
GO
