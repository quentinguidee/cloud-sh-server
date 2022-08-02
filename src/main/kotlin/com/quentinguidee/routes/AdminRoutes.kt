package com.quentinguidee.routes

import com.quentinguidee.services.adminsServices
import com.quentinguidee.services.authServices
import com.quentinguidee.utils.UnauthorizedException
import com.quentinguidee.utils.ok
import com.quentinguidee.utils.user
import io.ktor.server.application.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.encodeToJsonElement

fun Route.adminRoutes() {
    route("/admin") {
        get("/auth") {
            val methods = authServices.methodsPrivate()

            methods.forEach { method ->
                method.clientSecret = "***"
            }

            call.respond(buildJsonObject {
                put("methods", Json.encodeToJsonElement(methods))
            })
        }

        post("/reset") {
            if (call.user.role != "admin") {
                throw UnauthorizedException(call.user)
            }

            adminsServices.reset()

            call.ok()
        }
    }
}
