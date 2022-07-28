package com.quentinguidee.routes.storage

import com.quentinguidee.models.db.AccessType
import com.quentinguidee.services.bucketService
import com.quentinguidee.services.nodeService
import com.quentinguidee.utils.userID
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import kotlinx.serialization.json.buildJsonArray
import java.util.*

fun Route.bucketRoutes() {
    route("/bucket") {
        get {
            val userID = call.userID

            val bucket = try {
                bucketService.bucket(userID)
            } catch (e: NoSuchElementException) {
                bucketService.createBucket(userID)
            }

            call.respond(bucket.toJSON())
        }
    }

    route("{bucket_uuid}") {
        get {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")

            if (!bucketService.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.userID)) {
                call.respondText("unauthorized access", status = HttpStatusCode.Unauthorized)
            }

            val parentUUID = call.parameters.getOrFail("parent_uuid")

            val nodes = nodeService.nodes(parentUUID)

            call.respond(buildJsonArray {
                nodes.forEach { add(it.toJSON()) }
            })
        }
    }
}
