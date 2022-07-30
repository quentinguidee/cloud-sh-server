package com.quentinguidee.services.storage

import com.quentinguidee.dao.bucketsDAO
import com.quentinguidee.dao.nodesDAO
import com.quentinguidee.dao.usersBucketsDAO
import com.quentinguidee.models.AccessType
import com.quentinguidee.models.BucketType
import org.jetbrains.exposed.sql.transactions.transaction
import java.nio.file.Files
import java.util.*
import kotlin.io.path.Path

class BucketsServices {
    suspend fun bucket(userID: Int) = transaction {
        return@transaction bucketsDAO.getUserBucket(userID)
    }

    suspend fun createBucket(userID: Int) =
        createBucket(userID, "User bucket", BucketType.USER_BUCKET)

    private suspend fun createBucket(userID: Int, name: String, type: BucketType) = transaction {
        val bucket = bucketsDAO.create(name, type)
        val bucketUUID = UUID.fromString(bucket.uuid)
        usersBucketsDAO.create(bucketUUID, userID, AccessType.ADMIN)
        nodesDAO.create(bucketUUID, "root", "directory")

        val path = Path("data", "buckets", bucket.uuid, "root")
        try {
            Files.createDirectories(path)
        } catch (e: Exception) {
            rollback()
        }

        return@transaction bucket
    }

    suspend fun authorize(desiredAccessType: AccessType, bucketUUID: UUID, userID: Int) = transaction {
        val accessType = usersBucketsDAO.get(bucketUUID, userID).accessType
        println("${accessType.ordinal} ${desiredAccessType.ordinal}")
        return@transaction accessType >= desiredAccessType
    }

    suspend fun getRoot(bucketUUID: UUID) = transaction {
        nodesDAO.getRoot(bucketUUID)
    }
}

val bucketsServices = BucketsServices()
