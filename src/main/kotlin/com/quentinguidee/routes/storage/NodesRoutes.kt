package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.ok
import com.quentinguidee.utils.user
import io.ktor.server.application.*
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
    }
}
