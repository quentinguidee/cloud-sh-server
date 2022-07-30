package com.quentinguidee.routes

import com.quentinguidee.services.adminsServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.ok
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.routing.*

fun Route.adminRoutes() {
    route("/admin") {
        post("/reset") {
            if (call.user.role != "admin") {
                throw UnauthorizedException(call.user)
            }

            adminsServices.reset()

            call.ok()
        }
    }
}
