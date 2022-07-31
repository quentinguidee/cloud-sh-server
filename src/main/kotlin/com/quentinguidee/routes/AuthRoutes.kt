package com.quentinguidee.routes

import com.quentinguidee.services.authServices
import com.quentinguidee.services.sessionsServices
import com.quentinguidee.services.usersServices
import com.quentinguidee.utils.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
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
    val environment = environment

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

        route("/github") {
            val oAuthConfig = OAuthConfig(
                clientID = environment!!.config.property("auth.github.client_id").getString(),
                clientSecret = environment.config.property("auth.github.client_secret").getString(),
                authorizeURL = "https://github.com/login/oauth/authorize",
                accessTokenURL = "https://github.com/login/oauth/access_token",
                redirectURL = "http://localhost:3000/login",
            )

            val oAuth = OAuth(oAuthConfig)

            get("/login") {
                call.respond(buildJsonObject {
                    put("url", oAuth.getLoginURL())
                })
            }

            post("/callback") {
                val params = call.receive<CallbackParams>()

                val code = params.code

                // TODO: Check that states are equals
                // val state = params.state

                // TODO: Handle exchange fail
                val token = oAuth.exchange(oAuthConfig, code).accessToken
                val githubUserBody = authServices.fetchGitHubUser(token)

                val session = try {
                    val githubUser = authServices.githubUser(githubUserBody.login)
                    sessionsServices.createSession(githubUser.userID)
                } catch (e: NoSuchElementException) {
                    authServices.createAccount(githubUserBody)
                }

                val user = usersServices.get(session.userID)

                call.respond(json(session) {
                    putObject("user", user)
                })
            }
        }
    }
}
