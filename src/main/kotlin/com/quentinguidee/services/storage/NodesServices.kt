package com.quentinguidee.services.storage

import com.quentinguidee.dao.nodesDAO
import com.quentinguidee.models.Node
import com.quentinguidee.utils.MaxRecursionLevelException
import org.jetbrains.exposed.sql.transactions.transaction
import java.nio.file.Path
import java.util.*
import kotlin.io.path.Path
import kotlin.io.path.createDirectory
import kotlin.io.path.createFile

class NodesServices {
    suspend fun getChildren(parentUUID: String) =
        getChildren(UUID.fromString(parentUUID))

    private suspend fun getChildren(parentUUID: UUID) = transaction {
        nodesDAO.getChildren(parentUUID)
    }

    suspend fun getBin(bucketUUID: UUID) = transaction {
        nodesDAO.getDeleted(bucketUUID)
    }

    private fun getNodePath(of: Node): Path {
        var node = of
        var path = Path(node.name)
        var parentUUID = node.parentUUID

        var i = 100
        while (parentUUID != null) {
            node = nodesDAO.get(UUID.fromString(parentUUID))
            path = Path(node.name, path.toString())
            parentUUID = node.parentUUID
            if (parentUUID == null)
                break
            if (i-- <= 0)
                throw MaxRecursionLevelException()
        }

        return Path("data", "buckets", node.bucketUUID, path.toString())
    }

    suspend fun create(bucketUUID: UUID, parentUUID: UUID, name: String, type: String) = transaction {
        val node = nodesDAO.create(bucketUUID, parentUUID, name, type)
        val path = getNodePath(node)
        if (type == "directory") {
            path.createDirectory()
        } else {
            path.createFile()
        }
        return@transaction node
    }

    suspend fun softDelete(nodeUUID: UUID) = transaction {
        nodesDAO.softDelete(nodeUUID)
    }

    suspend fun forceDeleteRecursively(nodeUUID: UUID) {
        val node = transaction {
            nodesDAO.get(nodeUUID)
        }
        forceDeleteRecursively(node)
    }

    suspend fun forceDeleteRecursively(node: Node) {
        val nodes = transaction {
            nodesDAO.getChildren(UUID.fromString(node.uuid))
        }
        for (n in nodes) {
            forceDeleteRecursively(n)
        }
        forceDelete(node)
    }

    private suspend fun forceDelete(node: Node) = transaction {
        val path = getNodePath(node)
        nodesDAO.delete(UUID.fromString(node.uuid))
        path.toFile().deleteRecursively()
    }

    suspend fun emptyBin(bucketUUID: UUID) {
        val nodes = transaction {
            nodesDAO.getDeleted(bucketUUID)
        }
        for (node in nodes) {
            forceDeleteRecursively(node)
        }
    }
}

val nodesServices = NodesServices()
