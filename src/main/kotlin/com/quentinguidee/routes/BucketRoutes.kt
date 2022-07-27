package com.quentinguidee.routes

import com.quentinguidee.models.UserSession
import com.quentinguidee.services.bucketService
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.sessions.*

fun Route.bucketRoutes() {
    get {
        val session = call.sessions.get<UserSession>() ?: return@get call.respondText(
            "failed to retrieve user session",
            status = HttpStatusCode.InternalServerError
        )

        var bucket = bucketService.bucket(session.userID)
        if (bucket == null) {
            bucket = bucketService.createBucket(session.userID)
        }

        call.respond(bucket.toJSON())
    }
}
