package com.quentinguidee.services

import com.quentinguidee.models.db.Node
import com.quentinguidee.models.db.Nodes
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

class NodeService {
    suspend fun nodes(parentUUID: String) =
        nodes(UUID.fromString(parentUUID))

    private suspend fun nodes(parentUUID: UUID) = transaction {
        Node
            .find { Nodes.parent eq parentUUID }
            .toList()
    }
}

val nodeService = NodeService()
