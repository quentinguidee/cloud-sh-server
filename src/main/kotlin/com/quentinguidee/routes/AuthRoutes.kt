package com.quentinguidee.routes

import com.quentinguidee.services.authServices
import com.quentinguidee.services.sessionsServices
import com.quentinguidee.utils.OAuth
import com.quentinguidee.utils.OAuthConfig
import com.quentinguidee.utils.ok
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put

@Serializable
data class CallbackParams(
    val code: String,
    val state: String,
)

@Serializable
data class LogoutParams(
    val token: String,
)

fun Route.authRoutes() {
    val environment = environment

    route("/auth") {
        post("/logout") {
            val params = call.receive<LogoutParams>()
            sessionsServices.revokeSession(params.token)
            call.ok()
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

                call.respond(session)
            }
        }
    }
}
