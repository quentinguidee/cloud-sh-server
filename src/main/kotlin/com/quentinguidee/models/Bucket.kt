package com.quentinguidee.models

import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable

enum class BucketType {
    USER_BUCKET
}

object Buckets : UUIDTable("buckets") {
    val name = varchar("name", 255)
    val type = enumerationByName("type", 63, BucketType::class)
    val size = integer("size").default(0)
    val maxSize = integer("max_size").nullable()
}

@Serializable
data class Bucket(
    val uuid: String,
    val name: String,
    val type: BucketType,
    val size: Int,
    val maxSize: Int?,
)
