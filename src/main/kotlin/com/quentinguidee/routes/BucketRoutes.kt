package com.quentinguidee.routes

import com.quentinguidee.services.bucketService
import com.quentinguidee.utils.session
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*

fun Route.bucketRoutes() {
    get {
        val userID = call.session.user.id.value

        val bucket = bucketService.bucket(userID) ?: bucketService.createBucket(userID)

        call.respond(bucket.toJSON())
    }
}
