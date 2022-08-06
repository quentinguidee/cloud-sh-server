package com.quentinguidee.routes.storage

import com.quentinguidee.services.storage.bucketsServices
import com.quentinguidee.utils.json
import com.quentinguidee.utils.putObject
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
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
                putObject("rootNode", rootNode)
            })
        }
    }
}
