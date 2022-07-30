package com.quentinguidee.routes.storage

import com.quentinguidee.models.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.json
import com.quentinguidee.utils.putObject
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import java.util.*

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

    route("{bucket_uuid}") {
        get {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")
            val parentUUID = call.parameters.getOrFail("parent_uuid")

            if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.user.id)) {
                throw UnauthorizedException(call.user)
            }

            val nodes = nodesServices.getChildren(parentUUID)

            call.respond(nodes)
        }
    }
}
