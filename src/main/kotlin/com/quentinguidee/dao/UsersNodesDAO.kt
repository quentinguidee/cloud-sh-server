package com.quentinguidee.dao

import com.quentinguidee.models.UsersNodes
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.deleteWhere
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.update
import java.time.LocalDateTime
import java.util.*

class UsersNodesDAO {
    fun create(nodeUUID: UUID, userID: Int) = UsersNodes.insert {
        it[node] = nodeUUID
        it[user] = userID
        it[seenAt] = LocalDateTime.now()
        it[editedAt] = LocalDateTime.now()
    }

    fun delete(nodeUUID: UUID) = UsersNodes
        .deleteWhere { UsersNodes.node eq nodeUUID }

    fun updateSeenAt(uuid: UUID, userID: Int) = UsersNodes
        .update({ UsersNodes.user eq userID and (UsersNodes.node eq uuid) }) {
            it[seenAt] = LocalDateTime.now()
        }
}

val usersNodesDAO = UsersNodesDAO()
