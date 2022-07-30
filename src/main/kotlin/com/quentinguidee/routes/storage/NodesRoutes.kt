package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.ok
import com.quentinguidee.utils.user
import io.ktor.http.content.*
import io.ktor.server.application.*
import io.ktor.server.plugins.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import java.util.*

fun Route.nodesRoutes() {
    route("/nodes") {
        get {
            val parentUUID = call.parameters.getOrFail("parent_uuid")
            val bucketUUID = nodesServices.getNode(UUID.fromString(parentUUID)).bucketUUID

            if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            val nodes = nodesServices.getChildren(parentUUID)

            call.respond(nodes)
        }

        put {
            val parentUUID = call.parameters.getOrFail("parent_uuid")
            val bucketUUID = nodesServices.getNode(UUID.fromString(parentUUID)).bucketUUID
            val type = call.parameters.getOrFail("type")
            val name = call.parameters.getOrFail("name")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            nodesServices.create(
                UUID.fromString(bucketUUID),
                UUID.fromString(parentUUID),
                name,
                type,
            )

            call.ok()
        }

        patch {
            val nodeUUID = call.parameters.getOrFail("node_uuid")
            val bucketUUID = nodesServices.getNode(UUID.fromString(nodeUUID)).bucketUUID
            val name = call.parameters.getOrFail("new_name")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            nodesServices.rename(UUID.fromString(nodeUUID), name)

            call.ok()
        }

        delete {
            val nodeUUID = call.parameters.getOrFail("node_uuid")
            val bucketUUID = nodesServices.getNode(UUID.fromString(nodeUUID)).bucketUUID
            val softDelete = call.parameters.get("soft_delete")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            if (softDelete == "false")
                nodesServices.forceDeleteRecursively(UUID.fromString(nodeUUID))
            else
                nodesServices.softDelete(UUID.fromString(nodeUUID))

            call.ok()
        }

        post("/upload") {
            val parentUUID = call.parameters.getOrFail("parent_uuid")
            val parent = nodesServices.getNode(UUID.fromString(parentUUID))
            val bucketUUID = parent.bucketUUID

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id))
                throw UnauthorizedException(call.user)

            var name = "unnamed"
            var bytes: ByteArray? = null
            call.receiveMultipart().forEachPart { part ->
                if (part !is PartData.FileItem)
                    throw BadRequestException("The upload is not a file")

                name = part.originalFileName ?: "unnamed"
                bytes = part.streamProvider().readBytes()
                part.dispose()
            }

            nodesServices.create(
                bucketUUID = UUID.fromString(bucketUUID),
                parentUUID = UUID.fromString(parentUUID),
                name = name,
                type = "file",
                size = bytes?.size ?: 0,
                bytes = bytes,
            )

            call.ok()
        }
    }
}
