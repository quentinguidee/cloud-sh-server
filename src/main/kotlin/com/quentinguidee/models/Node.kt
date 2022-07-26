package com.quentinguidee.models

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.UUIDTable
import java.util.*

object Nodes : UUIDTable() {
    val parent = reference("parent", Nodes).nullable()
    val bucket = reference("bucket", Buckets)
    val name = varchar("name", 255)
    val type = varchar("type", 255)
    val mime = varchar("mime", 255)
    val size = integer("size")
}

class Node(id: EntityID<UUID>) : UUIDEntity(id) {
    companion object : UUIDEntityClass<Node>(Nodes)

    val parent by Nodes.parent
    val bucket by Nodes.bucket
    val name by Nodes.name
    val type by Nodes.type
    val mime by Nodes.mime
    val size by Nodes.size

    fun toJSON(): JsonObject {
        return buildJsonObject {
            put("parent", parent?.value.toString())
            put("bucket", bucket.value.toString())
            put("name", name)
            put("type", type)
            put("mime", mime)
            put("size", size)
        }
    }
}
