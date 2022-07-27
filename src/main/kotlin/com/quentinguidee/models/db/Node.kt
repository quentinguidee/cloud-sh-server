package com.quentinguidee.models.db

import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.UUIDEntity
import org.jetbrains.exposed.dao.UUIDEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.UUIDTable
import org.jetbrains.exposed.sql.transactions.transaction
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

    var parent by Node optionalReferencedOn Nodes.parent
    var bucket by Bucket referencedOn Nodes.bucket
    var name by Nodes.name
    var type by Nodes.type
    var mime by Nodes.mime
    var size by Nodes.size

    fun toJSON() = transaction {
        return@transaction buildJsonObject {
            put("parent", parent?.id.toString())
            put("bucket", bucket.toJSON())
            put("name", name)
            put("type", type)
            put("mime", mime)
            put("size", size)
        }
    }
}
