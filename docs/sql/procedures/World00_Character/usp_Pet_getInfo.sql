CREATE PROCEDURE [dbo].[usp_Pet_getInfo]
@nPetNo int,
@nRet int OUTPUT
AS

SELECT
nOwnerNo,
sPetName,
nPetMind,
nPetStress
FROM tPetNew
WHERE nPetNo = @nPetNo

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END
GO
