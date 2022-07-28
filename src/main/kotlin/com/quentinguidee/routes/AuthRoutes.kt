package com.quentinguidee.routes

import com.quentinguidee.models.UserSession
import com.quentinguidee.services.authService
import com.quentinguidee.utils.OAuth
import com.quentinguidee.utils.OAuthConfig
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.request.*
import io.ktor.server.response.*
import io.ktor.server.routing.*
import io.ktor.server.sessions.*
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.sql.transactions.transaction

@Serializable
data class CallbackParams(
    val code: String?,
    val state: String?,
)

fun Route.authRoutes() {
    val environment = environment

    get("/session") {
        val userSession = call.sessions.get<UserSession>() ?: return@get call.respondText(
            "not authenticated",
            status = HttpStatusCode.Unauthorized
        )

        val session = authService.session(userSession.username)

        call.respond(session.toJSON())
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

            val code = params.code ?: return@post call.respondText(
                "missing code in request parameters",
                status = HttpStatusCode.BadRequest
            )

            // TODO: Check that states are equals
            val state = params.state ?: return@post call.respondText(
                "missing state in request parameters",
                status = HttpStatusCode.BadRequest
            )

            // TODO: Handle exchange fail
            val token = oAuth.exchange(oAuthConfig, code).accessToken
            val githubUserBody = authService.fetchGitHubUser(token)
            val githubUser = authService.githubUser(githubUserBody.login)

            val session = if (githubUser == null) {
                authService.createAccount(githubUserBody)
            } else {
                authService.session(githubUser.username)
            }

            transaction {
                val user = session.user
                call.sessions.set(UserSession(user.id.value, user.username))
            }

            call.respond(session.toJSON())
        }
    }
}
