package com.quentinguidee.services

import com.quentinguidee.models.db.*
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import org.jetbrains.exposed.sql.transactions.transaction

class BucketService {
    suspend fun bucket(userID: Int) = transaction {
        val query = Buckets
            .innerJoin(UserBuckets)
            .innerJoin(Users)
            .select {
                Users.id eq userID and
                        (UserBuckets.accessType eq AccessType.ADMIN) and
                        (Buckets.type eq BucketType.USER_BUCKET)
            }
            .firstOrNull() ?: return@transaction null

        val bucket = Bucket.wrapRow(query)

        bucket.rootNode = Node
            .find { Nodes.bucket eq bucket.id and Nodes.parent.isNull() }
            .firstOrNull()

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
}

val bucketService = BucketService()
