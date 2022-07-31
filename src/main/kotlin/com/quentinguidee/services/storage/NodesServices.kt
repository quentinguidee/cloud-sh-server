package com.quentinguidee.services.storage

import com.quentinguidee.dao.nodesDAO
import com.quentinguidee.dao.usersNodesDAO
import com.quentinguidee.models.Node
import com.quentinguidee.models.deduceNodeTypeByName
import com.quentinguidee.utils.MaxRecursionLevelException
import org.jetbrains.exposed.sql.transactions.transaction
import java.nio.file.Path
import java.util.*
import kotlin.io.path.Path
import kotlin.io.path.createDirectory
import kotlin.io.path.createFile
import kotlin.io.path.writeBytes

class NodesServices {
    fun getNode(uuid: UUID) = transaction {
        nodesDAO.get(uuid)
    }

    fun getChildren(parentUUID: String) =
        getChildren(UUID.fromString(parentUUID))

    private fun getChildren(parentUUID: UUID) = transaction {
        nodesDAO.getChildren(parentUUID)
    }

    fun getBin(bucketUUID: UUID) = transaction {
        nodesDAO.getDeleted(bucketUUID)
    }

    fun getRecent(bucketUUID: UUID, userID: Int) = transaction {
        nodesDAO.getRecent(bucketUUID, userID)
    }

    fun getFile(node: Node, userID: Int) = transaction {
        val file = getNodePath(node).toFile()
        usersNodesDAO.updateSeenAt(UUID.fromString(node.uuid), userID)
        return@transaction file
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

    fun create(
        userID: Int,
        bucketUUID: UUID,
        parentUUID: UUID,
        name: String,
        type: String,
        mime: String? = null,
        size: Int = 0,
        bytes: ByteArray? = null
    ) = transaction {
        var deducedType: String = type
        if (deducedType == "file")
            deducedType = deduceNodeTypeByName(name)

        val node = nodesDAO.create(bucketUUID, parentUUID, name, deducedType, mime, size)
        usersNodesDAO.create(UUID.fromString(node.uuid), userID)

        val path = getNodePath(node)
        if (bytes != null) {
            path.writeBytes(bytes)
        } else if (type == "directory") {
            path.createDirectory()
        } else {
            path.createFile()
        }
        return@transaction node
    }

    fun rename(nodeUUID: UUID, name: String) = transaction {
        val path = getNodePath(nodesDAO.get(nodeUUID))
        val node = nodesDAO.get(nodeUUID)
        var type = node.type
        if (type != "directory")
            type = deduceNodeTypeByName(name)

        nodesDAO.rename(nodeUUID, name, type)
        val file = path.toFile()
        file.renameTo(path.parent.resolve(name).toFile())
    }

    fun softDelete(nodeUUID: UUID) = transaction {
        nodesDAO.softDelete(nodeUUID)
    }

    fun forceDeleteRecursively(nodeUUID: UUID) {
        val node = transaction {
            nodesDAO.get(nodeUUID)
        }
        forceDeleteRecursively(node)
    }

    fun forceDeleteRecursively(node: Node) {
        val nodes = transaction {
            nodesDAO.getChildren(UUID.fromString(node.uuid))
        }
        for (n in nodes) {
            forceDeleteRecursively(n)
        }
        forceDelete(node)
    }

    private fun forceDelete(node: Node) = transaction {
        val path = getNodePath(node)
        usersNodesDAO.delete(UUID.fromString(node.uuid))
        nodesDAO.delete(UUID.fromString(node.uuid))
        path.toFile().deleteRecursively()
    }

    fun emptyBin(bucketUUID: UUID) {
        val nodes = transaction {
            nodesDAO.getDeleted(bucketUUID)
        }
        for (node in nodes) {
            forceDeleteRecursively(node)
        }
    }
}

val nodesServices = NodesServices()
