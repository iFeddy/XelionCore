CREATE PROCEDURE [dbo].[usp_Pet_setCreate]
@nOwerNo int,
@sPetName nvarchar(16),
@nPetMind tinyint,
@nPetStress tinyint,
@nRet int OUTPUT
AS

INSERT INTO tPetNew (nOwnerNo, sPetName, nPetMind, nPetStress)
VALUES (@nOwerNo, @sPetName, @nPetMind, @nPetStress)

SET @nRet = @@ROWCOUNT

IF @@ERROR <> 0 OR @@ROWCOUNT <> 1
BEGIN
	SET @nRet = -1
	RETURN
END

SELECT @@IDENTITY
GO
