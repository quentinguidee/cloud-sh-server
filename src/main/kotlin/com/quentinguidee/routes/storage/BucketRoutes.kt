package com.quentinguidee.routes.storage

import com.quentinguidee.models.db.AccessType
import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.services.storage.nodesServices
import com.quentinguidee.utils.userID
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import java.util.*

fun Route.bucketRoutes() {
    route("/bucket") {
        get {
            val userID = call.userID

            val bucket = try {
                bucketsServices.bucket(userID)
            } catch (e: NoSuchElementException) {
                bucketsServices.createBucket(userID)
            }

            call.respond(bucket)
        }
    }

    route("{bucket_uuid}") {
        get {
            val bucketUUID = call.parameters.getOrFail("bucket_uuid")

            if (!bucketsServices.authorize(AccessType.READ, UUID.fromString(bucketUUID), call.userID)) {
                call.respondText("unauthorized access", status = HttpStatusCode.Unauthorized)
            }

            val parentUUID = call.parameters.getOrFail("parent_uuid")

            val nodes = nodesServices.getChildren(parentUUID)

            call.respond(nodes)
        }
    }
}
