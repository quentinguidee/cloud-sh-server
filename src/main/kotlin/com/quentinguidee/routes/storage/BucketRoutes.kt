package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import kotlinx.serialization.Serializable
import java.util.*

@Serializable
data class CreateNodeParams(
    val name: String,
    val type: String,
)

fun Route.bucketRoutes() {
    route("/bucket") {
        get {
            val userID = call.user.id

            val bucket = try {
                bucketsServices.bucket(userID)
            } catch (e: NoSuchElementException) {
                bucketsServices.createBucket(userID)
            }

            val rootNode = bucketsServices.getRoot(UUID.fromString(bucket.uuid))

            call.respond(json(bucket) {
                putObject("root_node", rootNode)
            })
        }
    }

    route("/{bucket_uuid}") {
        get {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")
            val parentUUID = call.parameters.getOrFail("parent_uuid")

            if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.user.id)) {
                throw UnauthorizedException(call.user)
            }

            val nodes = nodesServices.getChildren(parentUUID)

            call.respond(nodes)
        }

        put {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")
            val parentUUID = call.parameters.getOrFail("parent_uuid")

            val params = call.receive<CreateNodeParams>()

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id)) {
                throw UnauthorizedException(call.user)
            }

            nodesServices.create(
                UUID.fromString(bucketUUID),
                UUID.fromString(parentUUID),
                params.name,
                params.type,
            )

            call.ok()
        }

        patch {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")
            val nodeUUID = call.parameters.getOrFail("node_uuid")
            val name = call.parameters.getOrFail("new_name")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id)) {
                throw UnauthorizedException(call.user)
            }

            nodesServices.rename(UUID.fromString(nodeUUID), name)

            call.ok()
        }

        delete {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")
            val nodeUUID = call.parameters.getOrFail("node_uuid")
            val softDelete = call.parameters.get("soft_delete")

            if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id)) {
                throw UnauthorizedException(call.user)
            }

            if (softDelete == "false") {
                nodesServices.forceDeleteRecursively(UUID.fromString(nodeUUID))
            } else {
                nodesServices.softDelete(UUID.fromString(nodeUUID))
            }

            call.ok()
        }

        route("/bin") {
            get {
                val bucketUUID = call.parameters.getOrFail("bucket_uuid")

                if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.user.id)) {
                    throw UnauthorizedException(call.user)
                }

                val nodes = nodesServices.getBin(UUID.fromString(bucketUUID))

                call.respond(nodes)
            }

            delete {
                val bucketUUID = call.parameters.getOrFail("bucket_uuid")

                if (!bucketsServices.authorize(AccessType.WRITE, UUID.fromString(bucketUUID), call.user.id)) {
                    throw UnauthorizedException(call.user)
                }

                nodesServices.emptyBin(UUID.fromString(bucketUUID))

                call.ok()
            }
        }
    }
}
