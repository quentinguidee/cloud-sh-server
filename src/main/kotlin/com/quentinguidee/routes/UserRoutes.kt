package com.quentinguidee.routes

import com.quentinguidee.services.usersServices
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*

fun Route.userRoutes() {
    route("/user") {
        get {
            call.respond(call.user)
        }

        patch {
            val user = usersServices.updateInfo(
                userID = call.user.id,
                name = call.parameters["name"],
                email = call.parameters["email"],
                profilePicture = call.parameters["profile_picture"],
            )

            call.respond(user)
        }

        get("/{username}") {
            val username = call.parameters.getOrFail("username")

            val user = usersServices.get(username)

            call.respond(user)
        }
    }
}
