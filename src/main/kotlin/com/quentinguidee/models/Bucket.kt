package com.quentinguidee.models

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.UUIDTable
import java.util.*

object Buckets : UUIDTable() {
    val name = varchar("name", 255)
    val type = varchar("type", 255)
    val size = integer("size").default(0)
    val maxSize = integer("max_size").nullable()
    val rootNode = reference("root_node", Nodes)
}

class Bucket(id: EntityID<UUID>) : UUIDEntity(id) {
    companion object : UUIDEntityClass<Bucket>(Buckets)

    val name by Buckets.name
    val type by Buckets.type
    val size by Buckets.size
    val maxSize by Buckets.maxSize
    val rootNode by Buckets.rootNode


    fun toJSON(): JsonObject {
        return buildJsonObject {
            put("name", name)
            put("type", type)
            put("size", size)
            put("max_size", maxSize)
            put("root_node", rootNode.value.toString())
        }
    }
}
