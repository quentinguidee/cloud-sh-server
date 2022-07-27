package com.quentinguidee.services

import com.quentinguidee.models.db.*
import org.jetbrains.exposed.sql.and
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
                        (Buckets.type eq "user_bucket")
            }
            .firstOrNull() ?: return@transaction null

        return@transaction Bucket.wrapRow(query)
    }
}

val bucketService = BucketService()
