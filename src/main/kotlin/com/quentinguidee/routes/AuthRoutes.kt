package com.quentinguidee.routes

import com.quentinguidee.services.authServices
import com.quentinguidee.services.sessionsServices
import com.quentinguidee.services.usersServices
import com.quentinguidee.utils.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.util.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.encodeToJsonElement
import kotlinx.serialization.json.put

@Serializable
data class CallbackParams(
    val code: String,
    val state: String,
)

fun Route.authRoutes() {
    route("/auth") {
        get {
            val methods = authServices.methods()

            call.respond(buildJsonObject {
                put("methods", Json.encodeToJsonElement(methods))
            })
        }

        authenticated {
            post("/logout") {
                sessionsServices.revokeSession(call.session.token)
                call.ok()
            }
        }

        route("/{oauth_method_name}") {
            get("/login") {
                val oAuthMethodName = call.parameters.getOrFail("oauth_method_name")
                val method = authServices.method(oAuthMethodName)

                call.respond(buildJsonObject {
                    put("url", method.getLoginURL())
                })
            }

            post("/callback") {
                val oAuthMethodName = call.parameters.getOrFail("oauth_method_name")
                val method = authServices.methodPrivate(oAuthMethodName)

                val params = call.receive<CallbackParams>()

                val code = params.code

                // TODO: Check that states are equals
                // val state = params.state

                // TODO: Handle exchange fail
                val token = method.exchange(code).accessToken
                val githubUserBody = authServices.fetchGitHubUser(token)

                val session = try {
                    val oAuthUser = authServices.oAuthUser(githubUserBody.login)
                    sessionsServices.createSession(oAuthUser.userID)
                } catch (e: NoSuchElementException) {
                    authServices.createAccount(githubUserBody, method)
                }

                val user = usersServices.get(session.userID)

                call.respond(json(session) {
                    putObject("user", user)
                })
            }
        }
    }
}
