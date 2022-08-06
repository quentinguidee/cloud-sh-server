package com.quentinguidee.models

import kotlinx.serialization.Serializable
import org.jetbrains.exposed.sql.Table

enum class AccessType {
    NONE,
    READ,
    WRITE,
    MODERATOR,
    ADMIN;
}

object UsersBuckets : Table("user_buckets") {
    val bucket = reference("bucket_uuid", Buckets)
    val user = reference("user_id", Users)
    val accessType = enumerationByName("access_type", 63, AccessType::class)

    override val primaryKey = PrimaryKey(bucket, user)
}

@Serializable
data class UserBucket(
    val bucketUUID: String,
    val userID: Int,
    val accessType: AccessType,
)
