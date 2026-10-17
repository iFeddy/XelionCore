CREATE PROCEDURE [dbo].[p_Event_Friend_Reward]
@nFriendCharNo	int,
@nItemKey	bigint,
@err 		smallint OUTPUT

AS
SET NOCOUNT ON

DECLARE
@nFriendUserNo		int,
@nLevel 			smallint,
@nRewardItemID		int,
@nLevelReward_1_23	int,
@nLevelReward_24_47	int,
@nLevelReward_48_100	int,
@nStorage_Min		smallint,
@nStorage_Max		smallint,
@nStorageType		tinyint,
@nItemID		int,
@nOptionType		smallint,
@nOptionData		bigint,
@nOldItemKey		bigint



-- Hard Fixed...

SET @nStorage_Max 		= 143 	--개정창고 4 맨끝 위치.

SET @nStorageType 		= 6 	--개정창고

SET @nLevelReward_1_23 		= 8004 	-- 축복의 릭스(강도1)
SET @nLevelReward_24_47 	= 8013 	-- 축복의 릭스(강도2)
SET @nLevelReward_48_100 	= 8022 	-- 축복의 릭스(강도3)

SET @nOptionType 		= 1	-- Max = 255


SELECT @nLevel = nLevel, @nFriendUserNo=nUserNo  FROM tCharacter WITH(READUNCOMMITTED) 
	WHERE nCharNo = @nFriendCharNo
IF @@ROWCOUNT = 0
BEGIN
	SET @err = 1
	PRINT('Not Exitist Friend')
	RETURN @@ERROR
END
PRINT('@nLeve=')
PRINT(@nLevel)

IF @nLevel > 0 AND @nLevel < 24
	SET @nRewardItemID = @nLevelReward_1_23
ELSE IF @nLevel > 23 AND @nLevel < 48
	SET @nRewardItemID = @nLevelReward_24_47
ELSE
	SET @nRewardItemID = @nLevelReward_48_100

PRINT('@nRewardItemID=')
PRINT(@nRewardItemID)

SELECT @nStorage_Min = ISNULL(MIN( nStorage ), 0) FROM tItem WHERE nOwner=@nFriendUserNo AND nStorageType =@nStorageType 
	AND ( nItemID= @nLevelReward_1_23 OR nItemID= @nLevelReward_24_47 OR nItemID= @nLevelReward_48_100 )

PRINT('@nStorage_Min=')
PRINT(@nStorage_Min)

-- 개정창고 2 3 4 번 다 썼다.(@nStorage_Min=36)
IF @nStorage_Min > 0 AND @nStorage_Min < 37 
	BEGIN
		PRINT('IF @nStorage_Min > 0 AND @nStorage_Min < 36 ')
		GOTO CRITICAL
	END

IF @nStorage_Min = 0
	BEGIN
		SET	@nStorage_Min = @nStorage_Max
		SET	@nOptionData = 5


	END
ELSE
	BEGIN
		SELECT @nItemID = I.nItemID, @nOptionData=O.nOptionData, @nOldItemKey=I.nItemKey FROM tItem I
			INNER JOIN tItemOptions O ON I.nItemKey = O.nItemKey AND O.nOptionType = @nOptionType
			WHERE I.nOwner=@nFriendUserNo AND I.nStorageType =@nStorageType 
			AND I.nStorage = @nStorage_Min 

		PRINT('@@ROWCOUNT=')
		PRINT(@@ROWCOUNT)

		PRINT('@nItemID=')
		PRINT(@nItemID)

		PRINT('@nOptionData=')
		PRINT(@nOptionData)

		IF @nRewardItemID = @nItemID
			BEGIN
				-- 기존값에 추가한다.
				IF @nOptionData > 0 AND @nOptionData < 6 
				BEGIN
					SET	@nOptionData = @nOptionData+5
					--SET	@nStorage_Min = @nStorage_Min
				END
				ELSE
				BEGIN
					SET 	@nOptionData = 5
					SET	@nStorage_Min = @nStorage_Min-1
					PRINT('SET	@nStorage_Min = @nStorage_Min-1')
					PRINT(@nStorage_Min)

				END

			END
		ELSE
			BEGIN
				SET	@nStorage_Min = @nStorage_Min-1
				SET	@nOptionData = 5
			END
		
	END

-- 넣으려는 자리에 다른 아이템이 있는지 체크.(이벤트 전에 들어가 있던 아이템)
IF @nOptionData = 5
BEGIN
	IF EXISTS (SELECT nItemID  FROM tItem WHERE nOwner=@nFriendUserNo AND nStorageType =@nStorageType 
		AND nStorage=@nStorage_Min )
	BEGIN
		GOTO CRITICAL
	END

END

-- 아템 넣기.
BEGIN TRANSACTION

-- 새롭게 아템을 넣는다.
IF @nOptionData = 5
	BEGIN
		PRINT('INSERT tItem nItemKey')
		PRINT(@nItemKey)

		INSERT INTO tItem(nItemKey, nStorageType, nOwner, nStorage, nItemID, nFlags)
			VALUES(@nItemKey, @nStorageType, @nFriendUserNo, @nStorage_Min, @nRewardItemID, 0 )
		IF @@ERROR <> 0
		BEGIN
			GOTO ERR
		END
		PRINT('INSERT tItemOptions nItemKey')
		PRINT(@nItemKey)
		INSERT INTO tItemOptions(nItemKey, nOptionType, nOptionData)
			VALUES(@nItemKey, @nOptionType, @nOptionData )
		IF @@ERROR <> 0
		BEGIN
			GOTO ERR
		END
	END
-- 기존에 개수를 늘린
GO
