package com.quentinguidee.dao

import com.quentinguidee.models.Node
import com.quentinguidee.models.Nodes
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import java.util.*

class NodesDAO {
    private fun toNode(row: ResultRow) = Node(
        uuid = row[Nodes.id].value.toString(),
        parentUUID = row[Nodes.parent]?.value.toString(),
        bucketUUID = row[Nodes.bucket].value.toString(),
        name = row[Nodes.name],
        type = row[Nodes.type],
        mime = row[Nodes.mime],
        size = row[Nodes.size],
        createdAt = row[Nodes.createdAt],
        updatedAt = row[Nodes.updatedAt],
        deletedAt = row[Nodes.deletedAt],
    )

    fun getChildren(parentUUID: UUID) = Nodes
        .select { Nodes.parent eq parentUUID }
        .map(::toNode)

    fun create(bucketUUID: UUID, name: String, type: String) = Nodes
        .insert {
            it[Nodes.bucket] = bucketUUID
            it[Nodes.name] = name
            it[Nodes.type] = type
        }.resultedValues!!.map(::toNode).first()

    fun getRoot(bucketUUID: UUID) = Nodes
        .select {
            Nodes.parent eq null and
                    (Nodes.bucket eq bucketUUID)
        }
        .map(::toNode)
        .first()
}

val nodesDAO = NodesDAO()
