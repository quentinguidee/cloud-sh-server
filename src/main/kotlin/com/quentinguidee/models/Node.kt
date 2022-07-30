package com.quentinguidee.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.UUIDTable

object Nodes : UUIDTable() {
    val parent = reference("parent_uuid", Nodes).nullable()
    val bucket = reference("bucket_uuid", Buckets)
    val name = varchar("name", 255)
    val type = varchar("type", 255)
    val mime = varchar("mime", 255).nullable()
    val size = integer("size").default(0)
}

@Serializable
data class Node(
    val uuid: String,
    @SerialName("parent_uuid")
    val parentUUID: String?,
    @SerialName("bucket_uuid")
    val bucketUUID: String,
    val name: String,
    val type: String,
    val mime: String?,
    val size: Int,
)
