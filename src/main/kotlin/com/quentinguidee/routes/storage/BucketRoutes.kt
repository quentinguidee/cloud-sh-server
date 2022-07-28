package com.quentinguidee.routes.storage

import com.quentinguidee.models.db.AccessType
import com.quentinguidee.services.bucketService
import com.quentinguidee.services.nodeService
import com.quentinguidee.utils.userID
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.serialization.json.buildJsonArray
import java.util.*

fun Route.bucketRoutes() {
    route("/bucket") {
        get {
            val userID = call.userID

            val bucket = bucketService.bucket(userID) ?: bucketService.createBucket(userID)

            call.respond(bucket.toJSON())
        }
    }

    route("{bucket_uuid}") {
        get {
            val bucketUUID = call.parameters["bucket_uuid"] ?: return@get call.respondText(
                "missing bucket_uuid",
                status = HttpStatusCode.BadRequest
            )

            if (!bucketService.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.userID)) {
                call.respondText(
                    "unauthorized access",
                    status = HttpStatusCode.Unauthorized
                )
            }

            val parentUUID = call.parameters["parent_uuid"] ?: return@get call.respondText(
                "missing parent_uuid",
                status = HttpStatusCode.BadRequest
            )

            val nodes = nodeService.nodes(parentUUID)

            call.respond(buildJsonArray {
                nodes.forEach { add(it.toJSON()) }
            })
        }
    }
}
