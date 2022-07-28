package com.quentinguidee.models.db

import org.jetbrains.exposed.sql.Table

enum class AccessType {
    NONE,
    READ,
    WRITE,
    MODERATOR,
    ADMIN;
}

object UserBuckets : Table("user_buckets") {
    val bucket = reference("bucket_uuid", Buckets)
    val user = reference("user_id", Users)
    val accessType = enumerationByName("access_type", 63, AccessType::class)

    override val primaryKey = PrimaryKey(bucket, user, name = "PK_user_buckets")
}
