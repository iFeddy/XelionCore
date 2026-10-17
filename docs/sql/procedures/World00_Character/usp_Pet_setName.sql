CREATE PROCEDURE [dbo].[usp_Pet_setName]
@nPetNo int,
@sPetName nvarchar(16),
@nRet int OUTPUT
AS

UPDATE tPetNew
SET sPetName = @sPetName
WHERE nPetNo = @nPetNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
GO
