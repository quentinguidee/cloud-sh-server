package com.quentinguidee.routes

import com.quentinguidee.services.bucketService
import com.quentinguidee.utils.userID
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.bucketRoutes() {
    route("/bucket") {
        get {
            val userID = call.userID

            val bucket = bucketService.bucket(userID) ?: bucketService.createBucket(userID)

            call.respond(bucket.toJSON())
        }
    }
}
