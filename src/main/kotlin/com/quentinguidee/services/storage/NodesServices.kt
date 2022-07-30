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
        bucketUUID: UUID,
        parentUUID: UUID,
        name: String,
        type: String,
        size: Int = 0,
        bytes: ByteArray? = null
    ) = transaction {
        val node = nodesDAO.create(bucketUUID, parentUUID, name, type, size)
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
        nodesDAO.rename(nodeUUID, name)
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
