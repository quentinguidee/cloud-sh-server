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

        val bucket = bucketService.bucket(session.userID) ?: return@get call.respondText(
            "failed to find the user bucket",
            status = HttpStatusCode.NotFound
        )

        call.respond(bucket.toJSON())
    }
}
