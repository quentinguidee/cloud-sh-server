package com.quentinguidee.services.storage

import com.quentinguidee.dao.nodesDAO
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

class NodesServices {
    suspend fun getChildren(parentUUID: String) =
        getChildren(UUID.fromString(parentUUID))

    private suspend fun getChildren(parentUUID: UUID) = transaction {
        nodesDAO.getChildren(parentUUID)
    }
}

val nodesServices = NodesServices()
