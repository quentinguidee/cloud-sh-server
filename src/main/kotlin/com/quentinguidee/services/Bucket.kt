package com.quentinguidee.services

import com.quentinguidee.models.db.*
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

class BucketService {
    suspend fun bucket(userID: Int) = transaction {
        val query = Buckets
            .innerJoin(UserBuckets)
            .innerJoin(Users)
            .select {
                Users.id eq userID and
                        (Buckets.type eq BucketType.USER_BUCKET)
            }
            .first()

        val bucket = Bucket.wrapRow(query)

        bucket.rootNode = Node
            .find { Nodes.bucket eq bucket.id and Nodes.parent.isNull() }
            .first()

        return@transaction bucket
    }

    suspend fun createBucket(userID: Int) =
        createBucket(userID, "User bucket", BucketType.USER_BUCKET)

    private suspend fun createBucket(userID: Int, name: String, type: BucketType) = transaction {
        val bucket = Bucket.new {
            this.name = name
            this.type = type
        }

        UserBuckets.insert {
            it[accessType] = AccessType.ADMIN
            it[user] = userID
            it[UserBuckets.bucket] = bucket.id
        }

        val rootNode = Node.new {
            this.bucket = bucket
            this.name = "root"
            this.type = "directory"
        }

        bucket.rootNode = rootNode

        return@transaction bucket
    }

    private suspend fun accessType(bucketUUID: UUID, userID: Int) = transaction {
        val query = UserBuckets
            .select {
                UserBuckets.user eq userID and
                        (UserBuckets.bucket eq bucketUUID)
            }
            .firstOrNull() ?: return@transaction null

        return@transaction query.let {
            it[UserBuckets.accessType]
        }
    }

    suspend fun authorize(desiredAccessType: AccessType, bucketUUID: UUID, userID: Int): Boolean {
        val accessType = accessType(bucketUUID, userID) ?: return false
        println("${accessType.ordinal} ${desiredAccessType.ordinal}")
        return accessType >= desiredAccessType
    }
}

val bucketService = BucketService()
