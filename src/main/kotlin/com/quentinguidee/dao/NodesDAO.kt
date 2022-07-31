package com.quentinguidee.dao

import com.quentinguidee.models.Node
import com.quentinguidee.models.Nodes
import com.quentinguidee.models.UsersNodes
import org.jetbrains.exposed.sql.*
import java.time.LocalDateTime
import java.util.*

class NodesDAO {
    private fun toNode(row: ResultRow) = Node(
        uuid = row[Nodes.id].value.toString(),
        parentUUID = row[Nodes.parent]?.value?.toString(),
        bucketUUID = row[Nodes.bucket].value.toString(),
        name = row[Nodes.name],
        type = row[Nodes.type],
        mime = row[Nodes.mime],
        size = row[Nodes.size],
        createdAt = row[Nodes.createdAt],
        updatedAt = row[Nodes.updatedAt],
        deletedAt = row[Nodes.deletedAt],
    )

    fun get(uuid: UUID) = Nodes
        .select { Nodes.id eq uuid }
        .map(::toNode)
        .first()

    fun softDelete(uuid: UUID) = Nodes
        .update({ Nodes.id eq uuid }) {
            it[deletedAt] = LocalDateTime.now()
        }

    fun delete(uuid: UUID) = Nodes
        .deleteWhere { Nodes.id eq uuid }

    fun getChildren(parentUUID: UUID) = Nodes
        .select { Nodes.parent eq parentUUID and (Nodes.deletedAt eq null) }
        .map(::toNode)

    fun getDeleted(bucketUUID: UUID) = Nodes
        .select { Nodes.bucket eq bucketUUID and (Nodes.deletedAt neq null) }
        .map(::toNode)

    fun getRecent(bucketUUID: UUID, userID: Int) = Nodes
        .innerJoin(UsersNodes)
        .select {
            UsersNodes.user eq userID and
                    (UsersNodes.seenAt neq null) and
                    (Nodes.deletedAt eq null) and
                    (Nodes.bucket eq bucketUUID)
        }
        .orderBy(UsersNodes.seenAt, SortOrder.DESC)
        .map(::toNode)

    fun create(bucketUUID: UUID, parentUUID: UUID? = null, name: String, type: String, size: Int = 0) = Nodes
        .insert {
            it[Nodes.parent] = parentUUID
            it[Nodes.bucket] = bucketUUID
            it[Nodes.name] = name
            it[Nodes.type] = type
            it[Nodes.size] = size
        }.resultedValues!!.map(::toNode).first()

    fun getRoot(bucketUUID: UUID) = Nodes
        .select {
            Nodes.parent eq null and
                    (Nodes.bucket eq bucketUUID)
        }
        .map(::toNode)
        .first()

    fun rename(uuid: UUID, name: String, type: String) = Nodes
        .update({ Nodes.id eq uuid }) {
            it[Nodes.name] = name
            it[Nodes.type] = type
        }
}

val nodesDAO = NodesDAO()
