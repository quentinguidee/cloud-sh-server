package com.quentinguidee.models.db

import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.sql.Table
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

enum class AccessType {
    ADMIN
}

object UserBuckets : Table("user_buckets") {
    val bucket = reference("bucket", Buckets)
    val user = reference("user", Users)
    val accessType = enumeration<AccessType>("access_type")

    override val primaryKey = PrimaryKey(bucket, user)
}

class UserBucket(id: EntityID<UUID>) : UUIDEntity(id) {
    companion object : UUIDEntityClass<Bucket>(Buckets)

    var bucket by Bucket referencedOn UserBuckets.bucket
    var user by User referencedOn UserBuckets.user
    var accessType by UserBuckets.accessType

    fun toJSON() = transaction {
        return@transaction buildJsonObject {
            put("bucket", bucket.toJSON())
            put("user", user.toJSON())
            put("access_type", accessType.name)
        }
    }
}
