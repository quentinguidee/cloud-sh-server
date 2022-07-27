package com.quentinguidee.models.db

import org.jetbrains.exposed.sql.Table

enum class AccessType {
    ADMIN
}

object UserBuckets : Table("user_buckets") {
    val bucket = reference("bucket_uuid", Buckets)
    val user = reference("user_id", Users)
    val accessType = enumeration<AccessType>("access_type")

    override val primaryKey = PrimaryKey(bucket, user, name = "PK_user_buckets")
}
